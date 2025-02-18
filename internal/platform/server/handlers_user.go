package server

import (
	"github.com/labstack/echo/v4"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/domain/objectstorage"
	"net/http"
)

// HandleUserQuota function to handle the /user/quota endpoint
// caution: this function only works with AccessKey and does not use SecretKey
//
//	@Summary		Fetch User Quota from AccessKey
//	@Description	Fetches Quota Information
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			access_key	header		string							true	"User given AccessKey"
//	@Success		200			{object}	objectstorage.UserQuotaResponse	"Successful response with user quota"
//	@Failure		400			{object}	string							"Bad Request"
//	@Failure		401			{object}	string							"Unauthorized"
//	@Failure		500			{object}	string							"Internal server error"
//	@Router			/api/user/quota [get]
func HandleUserQuota(s *Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req objectstorage.UserRequestMeta
		err := (&echo.DefaultBinder{}).BindHeaders(c, &req)
		if err != nil {
			s.logger.Error(err.Error())
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err = c.Validate(req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		userQuota, errUserQuota := s.db.UserQuota(s.Config.ObjectStorageConfigs, req)
		if errUserQuota.Message != nil {
			return c.JSON(errUserQuota.Code, errUserQuota.Message.Error())
		}
		return c.JSON(http.StatusOK, userQuota)
	}
}

// HandleUserIdentification function to handle the /user/id endpoint
// caution: this function only works with AccessKey and does not use SecretKey
//
//	@Summary		Fetch User Identification information using AccessKey
//	@Description	Fetch User Identification information using AccessKey
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			access_key	header		string										true	"User given AccessKey"
//	@Success		200			{object}	objectstorage.UserIdentificationResponse	"Successful response with user identification"
//	@Failure		400			{object}	map[string]string							"Bad Request"
//	@Failure		401			{object}	string										"Unauthorized"
//	@Failure		500			{object}	map[string]string							"Internal server error"
//	@Router			/api/user/id [get]
func HandleUserIdentification(s *Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req objectstorage.UserRequestMeta
		err := (&echo.DefaultBinder{}).BindHeaders(c, &req)
		if err != nil {
			s.logger.Error(err.Error())
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err = c.Validate(req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		userData, errUserID := s.db.UserIdentification(s.Config.ObjectStorageConfigs, req)
		if errUserID.Message != nil {
			return c.JSON(errUserID.Code, errUserID.Message.Error())
		}
		return c.JSON(http.StatusOK, userData)
	}
}
