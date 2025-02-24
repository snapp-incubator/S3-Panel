package server

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/domain/objectstorage"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/infra/config"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/platform/health"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/platform/repository"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	Config     config.Config
	cancelCtx  context.Context
	cancelFunc context.CancelFunc
	db         objectstorage.ObjectStorage
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

	err := s.registerCephRepository()
	if err != nil {
		return nil, err
	}

	s.registerRouter()
	s.registerRoutes()
	s.registerPruner()

	return s, nil
}

func (s *Server) registerCephRepository() error {
	s.db = repository.NewCephObjectStorage()
	return nil
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
	apiRoutesBuckets.Use(s.TimeOutMiddleware(30))
	{
		apiRoutesBuckets.GET("/list", s.HandleBucketList())
		apiRoutesBuckets.GET("/quota", s.HandleBucketQuota())
		apiRoutesBuckets.POST("/create", s.HandleBucketCreate)
		apiRoutesBuckets.DELETE("/delete", s.HandleBucketDelete())
	}

	apiRoutesObjects := apiRoutes.Group("/object")
	apiRoutesObjects.Use(s.TimeOutMiddleware(300))
	{
		apiRoutesObjects.GET("/list", s.HandleObjectList())
		apiRoutesObjects.POST("/upload", s.HandleObjectUpload())
		apiRoutesObjects.GET("/download", s.HandleObjectDownload())
		apiRoutesObjects.GET("/head", s.HandleObjectHead())
		apiRoutesObjects.DELETE("/delete", s.HandleObjectsDelete())
	}

	apiRoutesUsers := apiRoutes.Group("/user")
	apiRoutesUsers.Use(s.TimeOutMiddleware(30))
	{
		apiRoutesUsers.GET("/quota", s.HandleUserQuota())
		apiRoutesUsers.GET("/id", s.HandleUserIdentification())
	}
}

func (s *Server) registerPruner() {
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		shouldBreak := false
		for {
			select {
			case <-s.cancelCtx.Done():
				s.logger.Warn("Shutting Down Ticker due to canceling context")
				shouldBreak = true
			case <-ticker.C:
				s.logger.Warn(fmt.Sprintf("Triggered Pruner: %s", time.Now()))
				errPrune := PruneObjectPathDir(s.Config.ServerConfigs.DownloadPath)
				if errPrune != nil {
					s.logger.Error(errPrune.Error())
				}
			}
			if shouldBreak {
				break
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
		logger.Warn("Done Context. Shutting down services")
	case sig := <-sigChan:
		logger.Warn(fmt.Sprintf("Received %s signal, gracefully shutting down services", sig.String()))
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
