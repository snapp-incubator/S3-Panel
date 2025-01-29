package server

import (
	"github.com/labstack/echo/v4"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/domain/objectstorage"
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
//	@Failure		400			{object}	map[string]string					"Bad Request"
//	@Failure		500			{object}	map[string]string					"Internal server error"
//	@Router			/api/bucket/list [get]
func HandleBucketList(s *Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req objectstorage.BucketInfoRequestMeta
		err := (&echo.DefaultBinder{}).BindHeaders(c, &req)
		if err != nil {
			s.logger.Error(err.Error())
			return c.JSON(http.StatusBadRequest, err)
		}
		err = c.Validate(req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		buckets, errBucketList := s.db.BucketList(s.config.ObjectStorageConfigs, req)
		if errBucketList != nil {
			s.logger.Error(errBucketList.Error())
			return c.JSON(http.StatusInternalServerError, errBucketList)
		}
		return c.JSON(http.StatusOK, buckets)
	}
}

// HandleBucketQuota
//
//	@Summary		Quota of buckets of a user
//	@Description	Fetches Quota of buckets owned by a user that is specified by AccessKey and SecretKey
//	@Tags			Bucket
//	@Accept			json
//	@Produce		json
//	@Param			access_key	header		string								true	"User given AccessKey"
//	@Param			secret_key	header		string								true	"User given SecretKey"
//	@Success		200			{object}	[]objectstorage.BucketQuotaResponse	"Successful response with buckets quota"
//	@Failure		400			{object}	map[string]string					"Bad Request"
//	@Failure		500			{object}	map[string]string					"Internal server error"
//	@Router			/api/bucket/quota [get]
func HandleBucketQuota(s *Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req objectstorage.BucketInfoRequestMeta
		err := (&echo.DefaultBinder{}).BindHeaders(c, &req)
		if err != nil {
			s.logger.Error(err.Error())
			return c.JSON(http.StatusBadRequest, err)
		}
		err = c.Validate(req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		quotaInfo, errBucketQuota := s.db.BucketQuota(s.config.ObjectStorageConfigs, req)
		if errBucketQuota != nil {
			s.logger.Error(errBucketQuota.Error())
			return c.JSON(http.StatusInternalServerError, errBucketQuota)
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
//	@Failure		400			{object}	map[string]string					"Bad Request"
//	@Failure		500			{object}	map[string]string					"Internal server error"
//	@Router			/api/bucket/create [post]
func HandleBucketCreate(s *Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req objectstorage.BucketActionRequestMeta
		err := c.Bind(&req)
		if err != nil {
			s.logger.Error(err.Error())
			return c.JSON(http.StatusBadRequest, err)
		}

		err = (&echo.DefaultBinder{}).BindHeaders(c, &req)
		if err != nil {
			s.logger.Error(err.Error())
			return c.JSON(http.StatusBadRequest, err)
		}

		err = c.Validate(req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		createdBucket, errBucketCreate := s.db.BucketCreate(s.config.ObjectStorageConfigs, req)
		if errBucketCreate != nil {
			s.logger.Error(errBucketCreate.Error())
			return c.JSON(http.StatusInternalServerError, errBucketCreate)
		}
		if createdBucket.AlreadyExist {
			return c.JSON(http.StatusOK, createdBucket)
		}
		return c.JSON(http.StatusCreated, createdBucket)
	}
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
//	@Failure		400			{object}	map[string]string					"Bad Request"
//	@Failure		500			{object}	map[string]string					"Internal server error"
//	@Router			/api/bucket/delete [delete]
func HandleBucketDelete(s *Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req objectstorage.BucketActionRequestMeta
		err := c.Bind(&req)
		if err != nil {
			s.logger.Error(err.Error())
			return c.JSON(http.StatusBadRequest, err)
		}

		err = (&echo.DefaultBinder{}).BindHeaders(c, &req)
		if err != nil {
			s.logger.Error(err.Error())
			return c.JSON(http.StatusBadRequest, err)
		}

		err = c.Validate(req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		deleteBucket, errBucketDelete := s.db.BucketDelete(s.config.ObjectStorageConfigs, req)
		if errBucketDelete != nil {
			s.logger.Error(errBucketDelete.Error())
			return c.JSON(http.StatusInternalServerError, errBucketDelete)
		}
		return c.JSON(http.StatusOK, deleteBucket)
	}
}
