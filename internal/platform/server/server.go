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
)

type Server struct {
	Config     config.Config
	cancelCtx  context.Context
	cancelFunc context.CancelFunc
	db         objectstorage.ObjectStorage
	logger     *zap.Logger
	router     *echo.Echo
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

	return s, nil
}

func (s *Server) registerCephRepository() error {
	s.db = repository.NewCephObjectStorage()
	return nil
}

func (s *Server) registerRouter() {
	newRouter := echo.New()
	newRouter.Validator = &CustomValidator{validator: NewRawValidator()}
	s.router = newRouter
}

func (s *Server) registerRoutes() {
	s.router.GET("/health", health.HandleHealth())
	s.router.GET("/docs/*", echoSwagger.WrapHandler)

	apiRoutes := s.router.Group("/api",
		s.CORSMiddleware(),
		s.AuthMiddleware(),
		s.TimeOutMiddleware(),
	)

	apiRoutes.OPTIONS("/*", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})

	apiRoutesBuckets := apiRoutes.Group("/bucket")
	{
		apiRoutesBuckets.GET("/list", HandleBucketList(s))
		apiRoutesBuckets.GET("/quota", HandleBucketQuota(s))
		apiRoutesBuckets.POST("/create", HandleBucketCreate(s))
		apiRoutesBuckets.DELETE("/delete", HandleBucketDelete(s))
	}

	apiRoutesObjects := apiRoutes.Group("/object")
	{
		apiRoutesObjects.GET("/list", HandleObjectList(s))
		apiRoutesObjects.POST("/upload", HandleObjectUpload(s))
		apiRoutesObjects.GET("/download", HandleObjectDownload(s))
		apiRoutesObjects.GET("/head", HandleObjectHead(s))
		apiRoutesObjects.DELETE("/delete", HandleObjectsDelete(s))
	}

	apiRoutesUsers := apiRoutes.Group("/user")
	{
		apiRoutesUsers.GET("/quota", HandleUserQuota(s))
		apiRoutesUsers.GET("/id", HandleUserIdentification(s))
	}
}

func (s *Server) Start() error {
	return s.router.Start(fmt.Sprintf("%s:%s", s.Config.ServerConfigs.Address, s.Config.ServerConfigs.Port))
}

func (s *Server) ShutDown() error {
	err := s.router.Shutdown(s.cancelCtx)
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
