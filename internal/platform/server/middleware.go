package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"strings"
	"time"
)

// CORSMiddleware to handle CORS config
func (s *Server) CORSMiddleware() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     s.config.ServerCorsConfigs.AllowedOrigins, // Specify your allowed origin(s)
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	})
}

func (s *Server) AuthMiddleware() echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Skipper: func(c echo.Context) bool {
			if isServerAuthEnabled(s.config.ServerConfigs.AuthEnabled) {
				return false
			} else {
				return true
			}
		},
		KeyLookup: s.config.ServerConfigs.AuthKeyLookup,
		Validator: func(auth string, c echo.Context) (bool, error) {
			if len(strings.TrimSpace(s.config.ServerConfigs.AuthToken)) == 0 {
				s.logger.Warn("Authentication is enabled but token is not set")
			}
			if strings.Compare(auth, s.config.ServerConfigs.AuthToken) == 0 {
				return true, nil
			}
			s.logger.Info("Unauthenticated Request", zap.String("address", c.Request().Host))
			return false, nil
		},
	})
}

func (s *Server) TimeOutMiddleware() echo.MiddlewareFunc {
	// TODO: Make it a config
	return middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 30 * time.Second,
	})
}
