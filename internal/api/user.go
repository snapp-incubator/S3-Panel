package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.snapp.ir/platform/s3-panel/internal/storage"
	"gitlab.snapp.ir/platform/s3-panel/internal/storage/ceph"
)

// HandleUserQuota function to handle the /user/quota endpoint
// caution: this function only works with AccessKey and does not use SecretKey
//
//	@Summary		Fetch User Quota from AccessKey
//	@Description	Fetches Quota Information
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			access_key	header		string								true	"User given AccessKey"
//	@Success		200			{object}	storage.UserQuotaResponse		"Successful response with user quota"
//	@Failure		400			{object}	storage.OperationErrWithMsg	"Bad Request"
//	@Failure		401			{object}	string								"Unauthorized"
//	@Failure		422			{object}	storage.OperationErrWithMsg	"Action didn't complete"
//	@Failure		500			{object}	storage.OperationErrWithMsg	"Internal server error"
//	@Router			/api/user/quota [get]
func (s *Server) HandleUserQuota() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req storage.UserRequestMeta
		err := (&echo.DefaultBinder{}).BindHeaders(c, &req)
		if err != nil {
			s.logger.Error(err.Error())
			return c.JSON(http.StatusBadRequest, storage.OperationErrWithMsg{Message: err.Error()})
		}
		err = c.Validate(req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, storage.OperationErrWithMsg{Message: err.Error()})
		}

		radosClient, err := ceph.NewRadosClient(s.Config.ObjectStorage.URL, s.Config.ObjectStorage.AccessKeyAdmin, s.Config.ObjectStorage.SecretKeyAdmin)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, storage.OperationErrWithMsg{Message: err.Error()})
		}

		userID, _, errGetUser := findUserID(s, radosClient, req.AccessKey)
		if errGetUser != nil {
			return c.JSON(http.StatusUnprocessableEntity, storage.OperationErrWithMsg{Message: errGetUser.Error()})
		}
		req.UID = userID

		userQuota, errUserQuota := s.store.UserQuota(s.Config.ObjectStorage, req)
		if errUserQuota.Message != nil {
			s.logger.Error(errUserQuota.Message.Error())
			return c.JSON(errUserQuota.Code, storage.OperationErrWithMsg{Message: errUserQuota.Message.Error()})
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
//	@Success		200			{object}	storage.UserIdentificationResponse	"Successful response with user identification"
//	@Failure		400			{object}	storage.OperationErrWithMsg			"Bad Request"
//	@Failure		401			{object}	string										"Unauthorized"
//	@Failure		422			{object}	storage.OperationErrWithMsg			"Action didn't complete"
//	@Failure		500			{object}	storage.OperationErrWithMsg			"Internal server error"
//	@Router			/api/user/id [get]
func (s *Server) HandleUserIdentification() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req storage.UserRequestMeta
		err := (&echo.DefaultBinder{}).BindHeaders(c, &req)
		if err != nil {
			s.logger.Error(err.Error())
			return c.JSON(http.StatusBadRequest, storage.OperationErrWithMsg{Message: err.Error()})
		}
		err = c.Validate(req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, storage.OperationErrWithMsg{Message: err.Error()})
		}

		radosClient, err := ceph.NewRadosClient(s.Config.ObjectStorage.URL, s.Config.ObjectStorage.AccessKeyAdmin, s.Config.ObjectStorage.SecretKeyAdmin)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, storage.OperationErrWithMsg{Message: err.Error()})
		}

		userID, _, errGetUser := findUserID(s, radosClient, req.AccessKey)
		if errGetUser != nil {
			return c.JSON(http.StatusUnprocessableEntity, storage.OperationErrWithMsg{Message: errGetUser.Error()})
		}
		req.UID = userID

		userData, errUserID := s.store.UserIdentification(s.Config.ObjectStorage, req)
		if errUserID.Message != nil {
			s.logger.Error(errUserID.Message.Error())
			return c.JSON(errUserID.Code, storage.OperationErrWithMsg{Message: errUserID.Message.Error()})
		}
		return c.JSON(http.StatusOK, userData)
	}
}
