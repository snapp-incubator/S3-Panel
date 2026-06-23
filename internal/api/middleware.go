package api

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

// CORSMiddleware to handle CORS config
func (s *Server) CORSMiddleware() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     s.Config.Cors.AllowedOrigins, // Specify your allowed origin(s)
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderXCSRFToken, echo.HeaderAccessControlAllowMethods},
		AllowCredentials: true,
	})
}

// isServerAuthEnabled reports whether bearer-token auth is enabled. Auth is
// disabled by default and only turned on when the value is exactly "true".
func isServerAuthEnabled(s string) bool {
	return strings.EqualFold(strings.TrimSpace(s), "true")
}

func (s *Server) AuthMiddleware() echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Skipper: func(c echo.Context) bool {
			return !isServerAuthEnabled(s.Config.Server.AuthEnabled)
		},
		KeyLookup: s.Config.Server.AuthKeyLookup,
		Validator: func(auth string, c echo.Context) (bool, error) {
			if len(strings.TrimSpace(s.Config.Server.AuthToken)) == 0 {
				s.logger.Warn("Authentication is enabled but token is not set")
			}
			if strings.Compare(auth, s.Config.Server.AuthToken) == 0 {
				return true, nil
			}
			s.logger.Info("Unauthenticated Request", zap.String("address", c.Request().Host))
			return false, nil
		},
	})
}
