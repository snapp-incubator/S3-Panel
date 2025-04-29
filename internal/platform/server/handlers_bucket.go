package server

import (
	"github.com/labstack/echo/v4"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/domain/objectstorage"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/platform/repository"
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
//	@Param			access_key		header		string								true	"User given AccessKey"
//	@Param			secret_key		header		string								true	"User given SecretKey"
//	@Param			search_string	query		string								true	"search by given string, could be empty"
//	@Success		200				{object}	objectstorage.BucketListResponse	"Successful response with bucket list"
//	@Failure		400				{object}	objectstorage.OperationErrWithMsg	"Bad Request"
//	@Failure		401				{object}	string								"Unauthorized"
//	@Failure		422				{object}	objectstorage.OperationErrWithMsg	"Action didn't complete"
//	@Failure		500				{object}	objectstorage.OperationErrWithMsg	"Internal server error"
//	@Router			/api/bucket/list [get]
func (s *Server) HandleBucketList() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req objectstorage.BucketListAndQuotaRequestMeta
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

		radosClient, err := repository.NewRadosClient(s.Config.ObjectStorageConfigs.URL, s.Config.ObjectStorageConfigs.AccessKeyAdmin, s.Config.ObjectStorageConfigs.SecretKeyAdmin)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, objectstorage.OperationErrWithMsg{Message: err.Error()})
		}

		userID, errGetUser, _ := FindUserID(s, radosClient, req.AccessKey)
		if errGetUser != nil {
			return c.JSON(http.StatusUnprocessableEntity, objectstorage.OperationErrWithMsg{Message: errGetUser.Error()})
		}
		req.UID = userID

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
//	@Description	Fetches Quota of buckets owned by a user that is specified by AccessKey. if status code is 200 and the BucketQuotaResponse hard_bytes or hard_objects is -1, mean no limitation on quota_storage or quota_objects is applied.
//	@Tags			Bucket
//	@Accept			json
//	@Produce		json
//	@Param			access_key		header		string								true	"User given AccessKey"
//	@Param			secret_key		header		string								true	"User given SecretKey"
//	@Param			search_string	query		string								true	"search by given string, could be empty"
//	@Success		200				{object}	objectstorage.BucketQuotaResponse	"Successful response with buckets quota"
//	@Failure		400				{object}	objectstorage.OperationErrWithMsg	"Bad Request"
//	@Failure		401				{object}	string								"Unauthorized"
//	@Failure		422				{object}	objectstorage.OperationErrWithMsg	"Action didn't complete"
//	@Failure		500				{object}	objectstorage.OperationErrWithMsg	"Internal server error"
//	@Router			/api/bucket/quota [get]
func (s *Server) HandleBucketQuota() echo.HandlerFunc {
	return func(c echo.Context) error {
		var BucketQuotaV1Threshold = 20
		var req objectstorage.BucketListAndQuotaRequestMeta
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

		radosClient, err := repository.NewRadosClient(s.Config.ObjectStorageConfigs.URL, s.Config.ObjectStorageConfigs.AccessKeyAdmin, s.Config.ObjectStorageConfigs.SecretKeyAdmin)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, objectstorage.OperationErrWithMsg{Message: err.Error()})
		}

		userID, errGetUser, _ := FindUserID(s, radosClient, req.AccessKey)
		if errGetUser != nil {
			return c.JSON(http.StatusUnprocessableEntity, objectstorage.OperationErrWithMsg{Message: errGetUser.Error()})
		}
		req.UID = userID

		// first we try to get a BucketList to make sure the provided credentials are correct
		bucketList, errBucketList := s.db.BucketList(s.Config.ObjectStorageConfigs, req)
		if errBucketList.Message != nil {
			s.logger.Error(errBucketList.Message.Error())
			return c.JSON(errBucketList.Code, objectstorage.OperationErrWithMsg{Message: errBucketList.Message.Error()})
		}

		if bucketList.Total <= BucketQuotaV1Threshold {
			quotaInfo, errBucketQuota := s.db.BucketQuota(s.Config.ObjectStorageConfigs, req, bucketList)
			if errBucketQuota.Message != nil {
				s.logger.Error(errBucketQuota.Message.Error())
				return c.JSON(errBucketQuota.Code, objectstorage.OperationErrWithMsg{Message: errBucketQuota.Message.Error()})
			}
			return c.JSON(http.StatusOK, quotaInfo)
		} else {
			quotaInfo, errBucketQuota := s.db.BucketQuotaV2(s.Config.ObjectStorageConfigs, req, bucketList)
			if errBucketQuota.Message != nil {
				s.logger.Error(errBucketQuota.Message.Error())
				return c.JSON(errBucketQuota.Code, objectstorage.OperationErrWithMsg{Message: errBucketQuota.Message.Error()})
			}
			return c.JSON(http.StatusOK, quotaInfo)
		}
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
//	@Failure		403			{object}	objectstorage.OperationErrWithMsg	"Bucket Already Exists / Bucket creation quota exceeded"
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

	bucketListInput := objectstorage.BucketListAndQuotaRequestMeta{
		AccessKey: req.AccessKey,
		SecretKey: req.SecretKey,
		MaxKeys:   1000,
		Page:      1,
	}
	bucketList, errBucketList := s.db.BucketList(s.Config.ObjectStorageConfigs, bucketListInput)
	if errBucketList.Message != nil {
		return c.JSON(errBucketList.Code, objectstorage.OperationErrWithMsg{Message: errBucketList.Message.Error()})
	}
	for _, bucket := range bucketList.Items {
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

		radosClient, err := repository.NewRadosClient(s.Config.ObjectStorageConfigs.URL, s.Config.ObjectStorageConfigs.AccessKeyAdmin, s.Config.ObjectStorageConfigs.SecretKeyAdmin)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, objectstorage.OperationErrWithMsg{Message: err.Error()})
		}

		userID, errGetUser, _ := FindUserID(s, radosClient, req.AccessKey)
		if errGetUser != nil {
			return c.JSON(http.StatusUnprocessableEntity, objectstorage.OperationErrWithMsg{Message: errGetUser.Error()})
		}

		reqBucketQuota := objectstorage.BucketListAndQuotaRequestMeta{
			AccessKey: req.AccessKey,
			SecretKey: req.SecretKey,
			UID:       userID,
			MaxKeys:   1000,
			Page:      1,
		}
		bucketList, bucketListQuotaErr := s.db.BucketList(s.Config.ObjectStorageConfigs, reqBucketQuota)
		if bucketListQuotaErr.Message != nil {
			return c.JSON(bucketListQuotaErr.Code, bucketListQuotaErr.Message.Error())
		}

		bucketFound := false
		for _, bucket := range bucketList.Items {
			if bucket == req.Bucket {
				bucketFound = true
				limitedBucketList := objectstorage.BucketListResponse{Total: 1, Items: []string{bucket}}
				bucketQuota, bucketQuotaErr := s.db.BucketQuota(s.Config.ObjectStorageConfigs, reqBucketQuota, limitedBucketList)
				if bucketQuotaErr.Message != nil {
					return c.JSON(bucketQuotaErr.Code, bucketQuotaErr.Message.Error())
				}
				if bucketQuota.Items[0].UsedObjects != 0 {
					return c.JSON(http.StatusForbidden, "Bucket is not empty for deletion")
				}
				break
			}
		}

		if !bucketFound {
			return c.JSON(http.StatusUnprocessableEntity, objectstorage.OperationErrWithMsg{Message: language.ErrNoSuchBucket})
		}

		deleteBucket, errBucketDelete := s.db.BucketDelete(s.Config.ObjectStorageConfigs, req)
		if errBucketDelete.Message != nil {
			s.logger.Error(errBucketDelete.Message.Error())
			return c.JSON(errBucketDelete.Code, objectstorage.OperationErrWithMsg{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, deleteBucket)
	}
}
