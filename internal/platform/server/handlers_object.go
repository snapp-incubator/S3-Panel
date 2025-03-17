package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/domain/objectstorage"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/platform/repository"
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
//	@Failure		400			{object}	objectstorage.OperationErrWithMsg		"Bad Request"
//	@Failure		401			{object}	string									"Unauthorized"
//	@Failure		422			{object}	objectstorage.OperationErrWithMsg		"Action didn't complete"
//	@Failure		500			{object}	objectstorage.OperationErrWithMsg		"Internal server error"
//	@Router			/api/object/download [get]
func (s *Server) HandleObjectDownload() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req objectstorage.ObjectRequestMeta
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

		// check if object exists before downloading
		_, existsErr := s.db.ObjectHead(s.Config.ObjectStorageConfigs, req)
		if existsErr.Message != nil {
			s.logger.Error(existsErr.Message.Error())
			return c.JSON(existsErr.Code, objectstorage.OperationErrWithMsg{Message: existsErr.Message.Error()})
		}

		tmpPath, exists, errCreatePath := createObjectPath(s.Config.ServerConfigs.DownloadPath, req.AccessKey, req.Object)
		if errCreatePath != nil {
			s.logger.Error(errCreatePath.Error())
			return c.JSON(http.StatusInternalServerError, objectstorage.OperationErrWithMsg{Message: errCreatePath.Error()})
		}
		req.TemporaryPath = tmpPath

		if !exists {
			_, errObjectDownload := s.db.ObjectDownload(s.Config.ObjectStorageConfigs, req)
			if errObjectDownload.Message != nil {
				s.logger.Error(errObjectDownload.Message.Error())
				return c.JSON(errObjectDownload.Code, objectstorage.OperationErrWithMsg{Message: errObjectDownload.Message.Error()})
			}
		}
		return c.Attachment(req.TemporaryPath, req.Object)
	}
}

// HandleObjectUpload
//
//	@Summary		upload an object to bucket
//	@Description	This functions uploads an object to bucket.
//	@Tags			Object
//	@Accept			mpfd
//	@Produce		json
//	@Param			access_key	header		string								true	"User given AccessKey"
//	@Param			secret_key	header		string								true	"User given SecretKey"
//	@Param			bucket		formData	string								true	"bucket name"
//	@Success		200			{object}	objectstorage.ObjectUploadResponse	"Successful response with bucket upload"
//	@Failure		400			{object}	objectstorage.OperationErrWithMsg	"Bad Request"
//	@Failure		401			{object}	string								"Unauthorized"
//	@Failure		403			{object}	objectstorage.OperationErrWithMsg	"Forbidden"
//	@Failure		409			{object}	objectstorage.OperationErrWithMsg	"Already Exists"
//	@Failure		422			{object}	objectstorage.OperationErrWithMsg	"Action didn't complete"
//	@Failure		500			{object}	objectstorage.OperationErrWithMsg	"Internal server error"
//	@Router			/api/object/upload [post]
func (s *Server) HandleObjectUpload() echo.HandlerFunc {
	formFieldBucket := "bucket"
	var maxUploadSize int64 = 1024 * 1024 * 1024

	return func(c echo.Context) error {
		var req objectstorage.ObjectUploadRequestMeta
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

		req.Bucket = c.FormValue(formFieldBucket)
		if req.Bucket == "" {
			return c.JSON(http.StatusBadRequest, objectstorage.OperationErrWithMsg{Message: fmt.Sprintf("bucket can not be empty")})
		}

		form, errForm := c.MultipartForm()
		if errForm != nil {
			s.logger.Error(errForm.Error())
			return c.JSON(http.StatusBadRequest, objectstorage.OperationErrWithMsg{Message: errForm.Error()})
		}

		files := form.File["files"]

		hasLargeFile := false
		for _, file := range files {
			// do not allow uploading files larger than 1G
			if file.Size > maxUploadSize {
				hasLargeFile = true
				break
			}

			// check if object exists before upload
			headReq := objectstorage.ObjectRequestMeta{
				AccessKey: req.AccessKey,
				SecretKey: req.SecretKey,
				Bucket:    req.Bucket,
				Object:    file.Filename,
			}
			existOut, existsErr := s.db.ObjectHead(s.Config.ObjectStorageConfigs, headReq)
			if existsErr.Message != nil {
				s.logger.Error(existsErr.Message.Error())
				return c.JSON(existsErr.Code, objectstorage.OperationErrWithMsg{Message: existsErr.Message.Error()})
			} else if existOut.Exists {
				return c.JSON(http.StatusConflict, objectstorage.OperationErrWithMsg{Message: fmt.Sprintf("File %s already exists", file.Filename)})
			}
		}

		if hasLargeFile {
			return c.JSON(http.StatusForbidden, objectstorage.OperationErrWithMsg{Message: "files larger than 1G are not allowed to upload"})
		}

		for _, file := range files {
			_, errObjectUpload := s.db.ObjectUpload(s.Config.ObjectStorageConfigs, req, file)
			if errObjectUpload.Message != nil {
				s.logger.Error(errObjectUpload.Message.Error())
				return c.JSON(errObjectUpload.Code, objectstorage.OperationErrWithMsg{Message: errObjectUpload.Message.Error()})
			}
		}
		return c.JSON(http.StatusOK, objectstorage.ObjectUploadResponse{Created: true})
	}
}

// HandleObjectList
//
//	@Summary		List of objects of a bucket
//	@Description	Fetches list of buckets owned by a user that is specified by AccessKey and SecretKey
//	@Tags			Object
//	@Accept			json
//	@Produce		json
//	@Param			access_key		header		string								true	"User given AccessKey"
//	@Param			secret_key		header		string								true	"User given SecretKey"
//	@Param			bucket			query		string								true	"bucket name"
//	@Param			max_keys		query		string								true	"max_keys of pagination"
//	@Param			page			query		string								true	"page of pagination"
//	@Param			search_string	query		string								true	"search by given string, could be empty"
//	@Success		200				{object}	objectstorage.ObjectListResponse	"Successful response with bucket list"
//	@Failure		400				{object}	objectstorage.OperationErrWithMsg	"Bad Request"
//	@Failure		401				{object}	string								"Unauthorized"
//	@Failure		422				{object}	objectstorage.OperationErrWithMsg	"Action didn't complete"
//	@Failure		500				{object}	objectstorage.OperationErrWithMsg	"Internal server error"
//	@Router			/api/object/list [get]
func (s *Server) HandleObjectList() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req objectstorage.ObjectListRequestMeta
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

		reqList := objectstorage.BucketInfoRequestMeta{
			AccessKey: req.AccessKey,
			SecretKey: req.SecretKey,
			UID:       userID,
		}
		outputBucketQuota, errBucketQuota := s.db.BucketQuota(s.Config.ObjectStorageConfigs, reqList)
		if errBucketQuota.Message != nil {
			return c.JSON(errBucketQuota.Code, objectstorage.OperationErrWithMsg{Message: errBucketQuota.Message.Error()})
		}

		objects, errObjectList := s.db.ObjectList(s.Config.ObjectStorageConfigs, req)
		if errObjectList.Message != nil {
			s.logger.Error(errObjectList.Message.Error())
			return c.JSON(errObjectList.Code, objectstorage.OperationErrWithMsg{Message: errObjectList.Message.Error()})
		}

		var bucketFound = false
		for _, bucketData := range outputBucketQuota.Items {
			if bucketData.BucketName == req.Bucket {
				if objects.TotalMatchedItems%int(req.MaxKeys) == 0 {
					objects.TotalPages = objects.TotalMatchedItems / int(req.MaxKeys)
				} else {
					objects.TotalPages = (objects.TotalMatchedItems / int(req.MaxKeys)) + 1
				}
				bucketFound = true
			}
		}

		if !bucketFound {
			return c.JSON(http.StatusUnprocessableEntity, objectstorage.OperationErrWithMsg{Message: "Bucket Not Found"})
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
//	@Param			bucket		query		string								true	"bucket name"
//	@Param			objects		query		[]string							true	"objects names"
//	@Success		200			{object}	objectstorage.ObjectDeleteResponse	"Successful response with objects delete"
//	@Failure		400			{object}	objectstorage.OperationErrWithMsg	"Bad Request"
//	@Failure		401			{object}	string								"Unauthorized"
//	@Failure		403			{object}	objectstorage.OperationErrWithMsg	"Object Does not exist"
//	@Failure		422			{object}	objectstorage.OperationErrWithMsg	"Action didn't complete"
//	@Failure		500			{object}	objectstorage.OperationErrWithMsg	"Internal server error"
//	@Router			/api/object/delete [delete]
func (s *Server) HandleObjectsDelete() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req objectstorage.ObjectDeleteRequestMeta
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

		for _, obj := range req.Objects {
			// check if object exists before deleting
			headReq := objectstorage.ObjectRequestMeta{
				AccessKey: req.AccessKey,
				SecretKey: req.SecretKey,
				Bucket:    req.Bucket,
				Object:    obj,
			}
			existOut, existsErr := s.db.ObjectHead(s.Config.ObjectStorageConfigs, headReq)
			if existsErr.Message != nil {
				s.logger.Error(existsErr.Message.Error())
				return c.JSON(existsErr.Code, objectstorage.OperationErrWithMsg{Message: existsErr.Message.Error()})
			} else if !existOut.Exists {
				return c.JSON(http.StatusForbidden, objectstorage.OperationErrWithMsg{Message: fmt.Sprintf("File %s does not exist", obj)})
			}
		}

		objects, errObjectDelete := s.db.ObjectsDelete(s.Config.ObjectStorageConfigs, req)
		if errObjectDelete.Message != nil {
			s.logger.Error(errObjectDelete.Message.Error())
			return c.JSON(errObjectDelete.Code, objectstorage.OperationErrWithMsg{Message: errObjectDelete.Message.Error()})
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
//	@Failure		400			{object}	objectstorage.OperationErrWithMsg	"Bad Request"
//	@Failure		401			{object}	string								"Unauthorized"
//	@Failure		422			{object}	objectstorage.OperationErrWithMsg	"Action didn't complete"
//	@Failure		500			{object}	objectstorage.OperationErrWithMsg	"Internal server error"
//	@Router			/api/object/head [get]
func (s *Server) HandleObjectHead() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req objectstorage.ObjectRequestMeta
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

		objects, errObjectHead := s.db.ObjectHead(s.Config.ObjectStorageConfigs, req)
		if errObjectHead.Message != nil {
			s.logger.Error(errObjectHead.Message.Error())
			return c.JSON(errObjectHead.Code, objectstorage.OperationErrWithMsg{Message: errObjectHead.Message.Error()})
		}
		return c.JSON(http.StatusOK, objects)
	}
}
