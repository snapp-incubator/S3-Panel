package server

import (
	"github.com/labstack/echo/v4"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/domain/objectstorage"
	"net/http"
)

func HandleBucketList(s *Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req objectstorage.BucketRequestMeta
		err := c.Bind(&req)
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

func HandleBucketQuota(s *Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req objectstorage.BucketRequestMeta
		err := c.Bind(&req)
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

func HandleBucketCreate(s *Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req objectstorage.BucketRequestMeta
		err := c.Bind(&req)
		if err != nil {
			s.logger.Error(err.Error())
			return c.JSON(http.StatusBadRequest, err)
		}
		err = c.Validate(req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		createdBucket, errBucketList := s.db.BucketCreate(s.config.ObjectStorageConfigs, req)
		if errBucketList != nil {
			s.logger.Error(errBucketList.Error())
			return c.JSON(http.StatusInternalServerError, errBucketList)
		}
		if createdBucket.AlreadyExist {
			return c.JSON(http.StatusOK, createdBucket)
		}
		return c.JSON(http.StatusCreated, createdBucket)
	}
}

func HandleBucketDelete(s *Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}
