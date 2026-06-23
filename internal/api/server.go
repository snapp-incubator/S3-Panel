package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"

	"github.com/snapp-incubator/S3-Panel/internal/cache"
	"github.com/snapp-incubator/S3-Panel/internal/config"
	"github.com/snapp-incubator/S3-Panel/internal/health"
	"github.com/snapp-incubator/S3-Panel/internal/storage"
	"github.com/snapp-incubator/S3-Panel/internal/storage/ceph"
	"github.com/snapp-incubator/S3-Panel/internal/web"
)

type Server struct {
	Config     config.Config
	cancelCtx  context.Context
	cancelFunc context.CancelFunc
	store      storage.ObjectStorage
	cache      cache.ServerCache
	logger     *zap.Logger
	Router     *echo.Echo
}

func NewServer(ctx context.Context, cancelFunc context.CancelFunc, cfg config.Config, logger *zap.Logger) (*Server, error) {
	s := &Server{
		Config:     cfg,
		cancelCtx:  ctx,
		cancelFunc: cancelFunc,
		logger:     logger,
	}

	s.registerCephRepository()

	s.registerCache()
	err := s.initializeCache()
	if err != nil {
		return nil, err
	}

	s.registerRouter()
	s.registerRoutes()
	s.registerPruner()

	return s, nil
}

func (s *Server) registerCephRepository() {
	s.logger.Info("### Registering Ceph Repository ###")
	s.store = ceph.NewCephObjectStorage()
	s.logger.Info("### Ceph Repository Registered ###")
}

func (s *Server) registerCache() {
	s.logger.Info("### Registering Cache ###")
	cacheStore := sync.Map{}
	s.cache = cache.NewInMemoryCache(&cacheStore)
	s.logger.Info("### Cache Registered ###")
}

func (s *Server) initializeCache() error {
	s.logger.Info("### Initializing Cache ###")
	radosClient, err := ceph.NewRadosClient(s.Config.ObjectStorage.URL, s.Config.ObjectStorage.AccessKeyAdmin, s.Config.ObjectStorage.SecretKeyAdmin)
	if err != nil {
		return err
	}
	_, _, err = findUserID(s, radosClient, "X")
	defer s.logger.Info("### Cache Initialized ###")
	return err
}

func (s *Server) registerRouter() {
	newRouter := echo.New()
	newRouter.Validator = &CustomValidator{validator: NewRawValidator()}
	if s.Config.Server.ServeFrontend {
		newRouter.Use(frontendMiddleware())
		s.logger.Info("### Serving embedded frontend ###")
	} else {
		s.logger.Info("### Frontend serving disabled (API-only) ###")
	}
	s.Router = newRouter
}

// frontendMiddleware serves the embedded frontend SPA for any request that is
// not an API, health or docs route. HTML5 falls back to index.html so
// client-side routes (e.g. /object-storage/...) resolve to the app.
func frontendMiddleware() echo.MiddlewareFunc {
	return middleware.StaticWithConfig(middleware.StaticConfig{
		Filesystem: web.HTTPFS(),
		HTML5:      true,
		Skipper: func(c echo.Context) bool {
			p := c.Request().URL.Path
			return strings.HasPrefix(p, "/api") ||
				strings.HasPrefix(p, "/s3/api") ||
				strings.HasPrefix(p, "/health") ||
				strings.HasPrefix(p, "/docs")
		},
	})
}

func (s *Server) registerRoutes() {
	s.Router.GET("/health", health.HandleHealth)
	s.Router.GET("/docs/*", echoSwagger.WrapHandler)

	// Serve the API under both /api and /s3/api. The bundled frontend calls
	// /s3/api (a convention inherited from the central panel router); direct API
	// clients and the swagger docs use /api.
	s.registerAPIGroup("/api")
	s.registerAPIGroup("/s3/api")
}

func (s *Server) registerAPIGroup(prefix string) {
	apiRoutes := s.Router.Group(prefix,
		s.CORSMiddleware(),
		s.AuthMiddleware(),
	)

	apiRoutes.OPTIONS("/*", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})

	apiRoutesBuckets := apiRoutes.Group("/bucket")
	{
		apiRoutesBuckets.GET("/list", s.HandleBucketList())
		apiRoutesBuckets.GET("/quota", s.HandleBucketQuota())
		apiRoutesBuckets.POST("/create", s.HandleBucketCreate)
		apiRoutesBuckets.DELETE("/delete", s.HandleBucketDelete())
	}

	apiRoutesObjects := apiRoutes.Group("/object")
	{
		apiRoutesObjects.GET("/list", s.HandleObjectList())
		apiRoutesObjects.POST("/upload", s.HandleObjectUpload())
		apiRoutesObjects.GET("/download", s.HandleObjectDownload())
		apiRoutesObjects.GET("/head", s.HandleObjectHead())
		apiRoutesObjects.DELETE("/delete", s.HandleObjectsDelete())
		apiRoutesObjects.GET("/share", s.HandleObjectShare())
	}

	apiRoutesUsers := apiRoutes.Group("/user")
	{
		apiRoutesUsers.GET("/quota", s.HandleUserQuota())
		apiRoutesUsers.GET("/id", s.HandleUserIdentification())
	}
}

func (s *Server) registerPruner() {
	prunerInterval := 1 * time.Hour
	go func() {
		ticker := time.NewTicker(prunerInterval)
		for {
			select {
			case <-s.cancelCtx.Done():
				s.logger.Warn("Shutting Down Ticker due to canceling context")
				ticker.Stop()
				return
			case <-ticker.C:
				s.logger.Info("Triggered Pruner")
				errPrune := pruneDownloadDir(s.Config.Server.DownloadPath, prunerInterval)
				if errPrune != nil {
					s.logger.Error(errPrune.Error())
				}
			}
		}
	}()
}

func (s *Server) Start() error {
	return s.Router.Start(fmt.Sprintf("%s:%s", s.Config.Server.Address, s.Config.Server.Port))
}

func (s *Server) ShutDown() error {
	err := s.Router.Shutdown(s.cancelCtx)
	if err != nil {
		return err
	}
	return nil
}

func StartServer(ctx context.Context, cancelFunc context.CancelFunc, cfg config.Config, logger *zap.Logger) error {
	server, err := NewServer(ctx, cancelFunc, cfg, logger)
	if err != nil {
		return err
	}
	go func() {
		errServerStart := server.Start()
		if errServerStart != nil {
			logger.Error(fmt.Sprintf("Error Starting Router %s", errServerStart))
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		logger.Info("Done Context. Shutting down services")
	case sig := <-sigChan:
		logger.Info(fmt.Sprintf("Received %s signal, gracefully shutting down services", sig.String()))
	}

	// call cancel function of the Server
	server.cancelFunc()
	err = server.ShutDown()
	if err != nil {
		logger.Error("Failed to shutdown server:", zap.Error(err))
	}

	logger.Info("Goodbye...")
	os.Exit(0)

	return nil
}
