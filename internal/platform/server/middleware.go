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
		AllowOrigins:     s.Config.ServerCorsConfigs.AllowedOrigins, // Specify your allowed origin(s)
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderXCSRFToken, echo.HeaderAccessControlAllowMethods},
		AllowCredentials: true,
	})
}

func (s *Server) AuthMiddleware() echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Skipper: func(c echo.Context) bool {
			if isServerAuthEnabled(s.Config.ServerConfigs.AuthEnabled) {
				return false
			} else {
				return true
			}
		},
		KeyLookup: s.Config.ServerConfigs.AuthKeyLookup,
		Validator: func(auth string, c echo.Context) (bool, error) {
			if len(strings.TrimSpace(s.Config.ServerConfigs.AuthToken)) == 0 {
				s.logger.Warn("Authentication is enabled but token is not set")
			}
			if strings.Compare(auth, s.Config.ServerConfigs.AuthToken) == 0 {
				return true, nil
			}
			s.logger.Info("Unauthenticated Request", zap.String("address", c.Request().Host))
			return false, nil
		},
	})
}

func (s *Server) TimeOutMiddleware() echo.MiddlewareFunc {
	return middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 600 * time.Second,
	})
}
