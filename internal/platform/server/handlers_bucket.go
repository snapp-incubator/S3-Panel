package server

import (
	"github.com/labstack/echo/v4"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/domain/objectstorage"
	language "gitlab.snapp.ir/platform/snapp_object_store/langs/en"
	"net/http"
)

// HandleBucketList
//
//	@Summary		List of buckets of a user
//	@Description	Fetches list of buckets owned by a user that is specified by AccessKey and SecretKey
//	@Tags			Bucket
//	@Accept			json
//	@Produce		json
//	@Param			access_key	header		string								true	"User given AccessKey"
//	@Param			secret_key	header		string								true	"User given SecretKey"
//	@Success		200			{object}	objectstorage.BucketListResponse	"Successful response with bucket list"
//	@Failure		400			{object}	objectstorage.OperationErrWithMsg	"Bad Request"
//	@Failure		401			{object}	string								"Unauthorized"
//	@Failure		422			{object}	objectstorage.OperationErrWithMsg	"Action didn't complete"
//	@Failure		500			{object}	objectstorage.OperationErrWithMsg	"Internal server error"
//	@Router			/api/bucket/list [get]
func (s *Server) HandleBucketList() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req objectstorage.BucketInfoRequestMeta
		err := (&echo.DefaultBinder{}).BindHeaders(c, &req)
		if err != nil {
			s.logger.Error(err.Error())
			return c.JSON(http.StatusBadRequest, objectstorage.OperationErrWithMsg{Message: err.Error()})
		}
		err = c.Validate(req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, objectstorage.OperationErrWithMsg{Message: err.Error()})
		}

		buckets, errBucketList := s.db.BucketList(s.Config.ObjectStorageConfigs, req)
		if errBucketList.Message != nil {
			s.logger.Error(errBucketList.Message.Error())
			return c.JSON(errBucketList.Code, objectstorage.OperationErrWithMsg{Message: errBucketList.Message.Error()})
		}
		return c.JSON(http.StatusOK, buckets)
	}
}

// HandleBucketQuota
//
//	@Summary		Quota of buckets of a user
//	@Description	Fetches Quota of buckets owned by a user that is specified by AccessKey
//	@Tags			Bucket
//	@Accept			json
//	@Produce		json
//	@Param			access_key	header		string								true	"User given AccessKey"
//	@Param			secret_key	header		string								true	"User given SecretKey"
//	@Success		200			{object}	[]objectstorage.BucketQuotaResponse	"Successful response with buckets quota"
//	@Failure		400			{object}	objectstorage.OperationErrWithMsg	"Bad Request"
//	@Failure		401			{object}	string								"Unauthorized"
//	@Failure		422			{object}	objectstorage.OperationErrWithMsg	"Action didn't complete"
//	@Failure		500			{object}	objectstorage.OperationErrWithMsg	"Internal server error"
//	@Router			/api/bucket/quota [get]
func (s *Server) HandleBucketQuota() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req objectstorage.BucketInfoRequestMeta
		err := (&echo.DefaultBinder{}).BindHeaders(c, &req)
		if err != nil {
			s.logger.Error(err.Error())
			return c.JSON(http.StatusBadRequest, objectstorage.OperationErrWithMsg{Message: err.Error()})
		}
		err = c.Validate(req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, objectstorage.OperationErrWithMsg{Message: err.Error()})
		}

		// first we try to get a BucketList to make sure the provided credentials are correct
		_, errBucketList := s.db.BucketList(s.Config.ObjectStorageConfigs, req)
		if errBucketList.Message != nil {
			s.logger.Error(errBucketList.Message.Error())
			return c.JSON(errBucketList.Code, objectstorage.OperationErrWithMsg{Message: errBucketList.Message.Error()})
		}

		quotaInfo, errBucketQuota := s.db.BucketQuota(s.Config.ObjectStorageConfigs, req)
		if errBucketQuota.Message != nil {
			s.logger.Error(errBucketQuota.Message.Error())
			return c.JSON(errBucketQuota.Code, objectstorage.OperationErrWithMsg{Message: errBucketQuota.Message.Error()})
		}
		return c.JSON(http.StatusOK, quotaInfo)
	}
}

// HandleBucketCreate
//
//	@Summary		Creates Bucket
//	@Description	Creates bucket for a user
//	@Tags			Bucket
//	@Accept			json
//	@Produce		json
//	@Param			access_key	header		string								true	"User given AccessKey"
//	@Param			secret_key	header		string								true	"User given SecretKey"
//	@Param			bucket		body		string								true	"Bucket Name to Create"
//	@Success		201			{object}	objectstorage.BucketCreateResponse	"Bucket created successfully"
//	@Failure		400			{object}	objectstorage.OperationErrWithMsg	"Bad Request"
//	@Failure		401			{object}	string								"Unauthorized"
//	@Failure		403			{object}	objectstorage.OperationErrWithMsg	"Bucket Already Exists"
//	@Failure		422			{object}	objectstorage.OperationErrWithMsg	"Action didn't complete"
//	@Failure		500			{object}	objectstorage.OperationErrWithMsg	"Internal server error"
//	@Router			/api/bucket/create [post]
func (s *Server) HandleBucketCreate(c echo.Context) error {
	var req objectstorage.BucketActionRequestMeta
	err := c.Bind(&req)
	if err != nil {
		s.logger.Error(err.Error())
		return c.JSON(http.StatusBadRequest, objectstorage.OperationErrWithMsg{Message: err.Error()})
	}

	err = (&echo.DefaultBinder{}).BindHeaders(c, &req)
	if err != nil {
		s.logger.Error(err.Error())
		return c.JSON(http.StatusBadRequest, objectstorage.OperationErrWithMsg{Message: err.Error()})
	}

	err = c.Validate(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, objectstorage.OperationErrWithMsg{Message: err.Error()})
	}

	bucketListInput := objectstorage.BucketInfoRequestMeta{
		AccessKey: req.AccessKey,
		SecretKey: req.SecretKey,
	}
	bucketList, errBucketList := s.db.BucketList(s.Config.ObjectStorageConfigs, bucketListInput)
	if errBucketList.Message != nil {
		return c.JSON(errBucketList.Code, objectstorage.OperationErrWithMsg{Message: errBucketList.Message.Error()})
	}
	for _, bucket := range bucketList.BucketList {
		if bucket == req.Bucket {
			return c.JSON(http.StatusForbidden, objectstorage.OperationErrWithMsg{Message: language.ErrBucketAlreadyExists})
		}
	}

	createdBucket, errBucketCreate := s.db.BucketCreate(s.Config.ObjectStorageConfigs, req)
	if errBucketCreate.Message != nil {
		s.logger.Error(errBucketCreate.Message.Error())
		return c.JSON(errBucketCreate.Code, objectstorage.OperationErrWithMsg{Message: errBucketCreate.Message.Error()})
	}
	return c.JSON(http.StatusCreated, createdBucket)
}

// HandleBucketDelete
//
//	@Summary		Deletes Bucket
//	@Description	Deletes bucket for a user
//	@Tags			Bucket
//	@Accept			json
//	@Produce		json
//	@Param			access_key	header		string								true	"User given AccessKey"
//	@Param			secret_key	header		string								true	"User given SecretKey"
//	@Param			bucket		body		string								true	"Bucket Name to Delete"
//	@Success		200			{object}	objectstorage.BucketDeleteResponse	"Bucket deleted successfully"
//	@Failure		400			{object}	objectstorage.OperationErrWithMsg	"Bad Request"
//	@Failure		401			{object}	string								"Unauthorized"
//	@Failure		422			{object}	objectstorage.OperationErrWithMsg	"Action didn't complete"
//	@Failure		500			{object}	objectstorage.OperationErrWithMsg	"Internal server error"
//	@Router			/api/bucket/delete [delete]
func (s *Server) HandleBucketDelete() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req objectstorage.BucketActionRequestMeta
		err := c.Bind(&req)
		if err != nil {
			s.logger.Error(err.Error())
			return c.JSON(http.StatusBadRequest, objectstorage.OperationErrWithMsg{Message: err.Error()})
		}

		err = (&echo.DefaultBinder{}).BindHeaders(c, &req)
		if err != nil {
			s.logger.Error(err.Error())
			return c.JSON(http.StatusBadRequest, objectstorage.OperationErrWithMsg{Message: err.Error()})
		}

		err = c.Validate(req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, objectstorage.OperationErrWithMsg{Message: err.Error()})
		}

		deleteBucket, errBucketDelete := s.db.BucketDelete(s.Config.ObjectStorageConfigs, req)
		if errBucketDelete.Message != nil {
			s.logger.Error(errBucketDelete.Message.Error())
			return c.JSON(errBucketDelete.Code, objectstorage.OperationErrWithMsg{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, deleteBucket)
	}
}
