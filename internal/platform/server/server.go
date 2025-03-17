package server

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/domain/cache"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/domain/objectstorage"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/infra/config"
	inMemoryCache "gitlab.snapp.ir/platform/snapp_object_store/internal/platform/cache"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/platform/health"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/platform/repository"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Server struct {
	Config     config.Config
	cancelCtx  context.Context
	cancelFunc context.CancelFunc
	db         objectstorage.ObjectStorage
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
	s.db = repository.NewCephObjectStorage()
	s.logger.Info("### Cache Registered ###")
}

func (s *Server) registerCache() {
	s.logger.Info("### Registering Cache ###")
	cacheStore := sync.Map{}
	s.cache = inMemoryCache.NewInMemoryCache(&cacheStore)
	s.logger.Info("### Cache Registered ###")
}

func (s *Server) initializeCache() error {
	s.logger.Info("### Initializing Cache ###")
	radosClient, err := repository.NewRadosClient(s.Config.ObjectStorageConfigs.URL, s.Config.ObjectStorageConfigs.AccessKeyAdmin, s.Config.ObjectStorageConfigs.SecretKeyAdmin)
	if err != nil {
		return err
	}
	_, err, _ = FindUserID(s, radosClient, "X")
	defer s.logger.Info("### Cache Initialized ###")
	return err
}

func (s *Server) registerRouter() {
	newRouter := echo.New()
	newRouter.Validator = &CustomValidator{validator: NewRawValidator()}
	s.Router = newRouter
}

func (s *Server) registerRoutes() {
	s.Router.GET("/health", health.HandleHealth)
	s.Router.GET("/docs/*", echoSwagger.WrapHandler)

	apiRoutes := s.Router.Group("/api",
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
				break
			case <-ticker.C:
				s.logger.Info("Triggered Pruner")
				errPrune := PruneObjectPathDir(s.Config.ServerConfigs.DownloadPath, prunerInterval)
				if errPrune != nil {
					s.logger.Error(errPrune.Error())
				}
			}
		}
	}()
}

func (s *Server) Start() error {
	return s.Router.Start(fmt.Sprintf("%s:%s", s.Config.ServerConfigs.Address, s.Config.ServerConfigs.Port))
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
