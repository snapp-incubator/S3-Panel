package server

import (
	"github.com/labstack/echo/v4"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/domain/objectstorage"
	"net/http"
)

// HandleObjectDownload
//
//	@Summary		download an object to bucket
//	@Description	This functions downloads an object from bucket.
//	@Tags			Object
//	@Accept			json
//	@Produce		json
//	@Param			access_key	header		string									true	"User given AccessKey"
//	@Param			secret_key	header		string									true	"User given SecretKey"
//	@Param			bucket		query		string									true	"bucket name"
//	@Param			object		query		string									true	"object name"
//	@Success		200			{object}	objectstorage.ObjectDownloadResponse	"Successful response with bucket download"
//	@Failure		400			{object}	string									"Bad Request"
//	@Failure		401			{object}	string									"Unauthorized"
//	@Failure		500			{object}	string									"Internal server error"
//	@Router			/api/object/download [get]
func HandleObjectDownload(s *Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req objectstorage.ObjectRequestMeta
		err := c.Bind(&req)
		if err != nil {
			s.logger.Error(err.Error())
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		err = (&echo.DefaultBinder{}).BindHeaders(c, &req)
		if err != nil {
			s.logger.Error(err.Error())
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		err = c.Validate(req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		// check if object exists before downloading
		_, existsErr := s.db.ObjectHead(s.config.ObjectStorageConfigs, req)
		if existsErr.Message != nil {
			return c.JSON(existsErr.Code, existsErr.Message.Error())
		}

		objects, errObjectDownload := s.db.ObjectDownload(s.config.ObjectStorageConfigs, req)
		if errObjectDownload.Message != nil {
			return c.JSON(errObjectDownload.Code, errObjectDownload.Message.Error())
		}
		return c.JSON(http.StatusOK, objects)
	}
}

// HandleObjectUpload
//
//	@Summary		upload an object to bucket
//	@Description	This functions uploads an object to bucket.
//	@Tags			Object
//	@Accept			json
//	@Produce		json
//	@Param			access_key	header		string								true	"User given AccessKey"
//	@Param			secret_key	header		string								true	"User given SecretKey"
//	@Param			bucket		query		string								true	"bucket name"
//	@Param			object		query		string								true	"object name"
//	@Param			content		query		string								true	"content"
//	@Success		200			{object}	objectstorage.ObjectUploadResponse	"Successful response with bucket upload"
//	@Failure		400			{object}	string								"Bad Request"
//	@Failure		401			{object}	string								"Unauthorized"
//	@Failure		500			{object}	string								"Internal server error"
//	@Router			/api/object/upload [put]
func HandleObjectUpload(s *Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req objectstorage.ObjectUploadRequestMeta
		err := c.Bind(&req)
		if err != nil {
			s.logger.Error(err.Error())
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		err = (&echo.DefaultBinder{}).BindHeaders(c, &req)
		if err != nil {
			s.logger.Error(err.Error())
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		err = c.Validate(req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		// check if object exists before upload
		headReq := objectstorage.ObjectRequestMeta{
			AccessKey: req.AccessKey,
			SecretKey: req.SecretKey,
			Bucket:    req.Bucket,
			Object:    req.Object,
		}
		_, existsErr := s.db.ObjectHead(s.config.ObjectStorageConfigs, headReq)
		if existsErr.Message != nil {
			return c.JSON(existsErr.Code, existsErr.Message.Error())
		}

		objects, errObjectUpload := s.db.ObjectUpload(s.config.ObjectStorageConfigs, req)
		if err != nil {
			return c.JSON(errObjectUpload.Code, errObjectUpload.Message.Error())
		}
		return c.JSON(http.StatusOK, objects)
	}
}

// HandleObjectList
//
//	@Summary		List of objects of a bucket
//	@Description	Fetches list of buckets owned by a user that is specified by AccessKey and SecretKey
//	@Tags			Object
//	@Accept			json
//	@Produce		json
//	@Param			access_key	header		string								true	"User given AccessKey"
//	@Param			secret_key	header		string								true	"User given SecretKey"
//	@Param			bucket		query		string								true	"bucket name"
//	@Param			max_keys	query		string								true	"max_keys of pagination"
//	@Param			page		query		string								true	"page of pagination"
//	@Success		200			{object}	objectstorage.ObjectListResponse	"Successful response with bucket list"
//	@Failure		400			{object}	string								"Bad Request"
//	@Failure		401			{object}	string								"Unauthorized"
//	@Failure		500			{object}	string								"Internal server error"
//	@Router			/api/object/list [get]
func HandleObjectList(s *Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req objectstorage.ObjectListRequestMeta
		err := c.Bind(&req)
		if err != nil {
			s.logger.Error(err.Error())
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		err = (&echo.DefaultBinder{}).BindHeaders(c, &req)
		if err != nil {
			s.logger.Error(err.Error())
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		err = c.Validate(req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		objects, errObjectList := s.db.ObjectList(s.config.ObjectStorageConfigs, req)
		if err != nil {
			return c.JSON(errObjectList.Code, errObjectList.Message.Error())
		}
		return c.JSON(http.StatusOK, objects)
	}
}

// HandleObjectsDelete
//
//	@Summary		Deletes the list of objects inside a bucket
//	@Description	Deletes a list of objects specified by name
//	@Tags			Object
//	@Accept			json
//	@Produce		json
//	@Param			access_key	header		string								true	"User given AccessKey"
//	@Param			secret_key	header		string								true	"User given SecretKey"
//	@Param			bucket		body		string								true	"bucket name"
//	@Param			objects		body		[]string							true	"objects names"
//	@Success		200			{object}	objectstorage.ObjectDeleteResponse	"Successful response with objects delete"
//	@Failure		400			{object}	string								"Bad Request"
//	@Failure		401			{object}	string								"Unauthorized"
//	@Failure		500			{object}	string								"Internal server error"
//	@Router			/api/object/delete [delete]
func HandleObjectsDelete(s *Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req objectstorage.ObjectDeleteRequestMeta
		err := c.Bind(&req)
		if err != nil {
			s.logger.Error(err.Error())
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		err = (&echo.DefaultBinder{}).BindHeaders(c, &req)
		if err != nil {
			s.logger.Error(err.Error())
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		err = c.Validate(req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		for _, obj := range req.Objects {
			// check if object exists before deleting
			headReq := objectstorage.ObjectRequestMeta{
				AccessKey: req.AccessKey,
				SecretKey: req.SecretKey,
				Bucket:    req.Bucket,
				Object:    obj,
			}
			_, existsErr := s.db.ObjectHead(s.config.ObjectStorageConfigs, headReq)
			if existsErr.Message != nil {
				return c.JSON(existsErr.Code, existsErr.Message.Error())
			}
		}

		objects, errObjectDelete := s.db.ObjectsDelete(s.config.ObjectStorageConfigs, req)
		if err != nil {
			return c.JSON(errObjectDelete.Code, errObjectDelete.Message.Error())
		}
		return c.JSON(http.StatusOK, objects)
	}
}

// HandleObjectHead
//
//	@Summary		check if an object exists
//	@Description	check if an object exists
//	@Tags			Object
//	@Accept			json
//	@Produce		json
//	@Param			access_key	header		string								true	"User given AccessKey"
//	@Param			secret_key	header		string								true	"User given SecretKey"
//	@Param			bucket		body		string								true	"bucket name"
//	@Param			object		body		string								true	"objects name"
//	@Success		200			{object}	objectstorage.ObjectHeadResponse	"Successful response with objects head"
//	@Failure		400			{object}	string								"Bad Request"
//	@Failure		401			{object}	string								"Unauthorized"
//	@Failure		500			{object}	string								"Internal server error"
//	@Router			/api/object/head [get]
func HandleObjectHead(s *Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req objectstorage.ObjectRequestMeta
		err := c.Bind(&req)
		if err != nil {
			s.logger.Error(err.Error())
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		err = (&echo.DefaultBinder{}).BindHeaders(c, &req)
		if err != nil {
			s.logger.Error(err.Error())
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		err = c.Validate(req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		objects, errObjectDelete := s.db.ObjectHead(s.config.ObjectStorageConfigs, req)
		if err != nil {
			return c.JSON(errObjectDelete.Code, errObjectDelete.Message.Error())
		}
		return c.JSON(http.StatusOK, objects)
	}
}
