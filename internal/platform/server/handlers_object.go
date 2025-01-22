package server

import (
	"github.com/labstack/echo/v4"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/domain/objectstorage"
	"net/http"
)

func HandleObjectDownload(s *Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}

func HandleObjectUpload(s *Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
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
//	@Failure		400			{object}	map[string]string					"Bad Request"
//	@Failure		500			{object}	map[string]string					"Internal server error"
//	@Router			/api/object/list [get]
func HandleObjectList(s *Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req objectstorage.ObjectRequestMeta
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

		objects, err := s.db.ObjectList(s.config.ObjectStorageConfigs, req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, objects)
	}
}

func HandleObjectDelete(s *Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}
