package server

import (
	"github.com/labstack/echo/v4"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/domain/objectstorage"
	"net/http"
)

// HandleUserQuota function to handle the /user/quota endpoint
// caution: this function only works with AccessKey and does not use SecretKey
func HandleUserQuota(s *Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req objectstorage.UserRequestMeta
		err := c.Bind(&req)
		if err != nil {
			s.logger.Error(err.Error())
			return c.JSON(http.StatusBadRequest, err)
		}
		err = c.Validate(req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		userQuota, err := s.db.UserQuota(s.config.ObjectStorageConfigs, req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, userQuota)
	}
}

// HandleUserIdentification function to handle the /user/id endpoint
// caution: this function only works with AccessKey and does not use SecretKey
func HandleUserIdentification(s *Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req objectstorage.UserRequestMeta
		err := c.Bind(&req)
		if err != nil {
			s.logger.Error(err.Error())
			return c.JSON(http.StatusBadRequest, err)
		}
		err = c.Validate(req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		userData, err := s.db.UserIdentification(s.config.ObjectStorageConfigs, req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		} else if userData.UserNotFound {
			return c.JSON(http.StatusUnauthorized, userData)
		}
		return c.JSON(http.StatusOK, userData)
	}
}
