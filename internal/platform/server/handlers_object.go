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

func HandleObjectList(s *Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req objectstorage.ObjectRequestMeta
		err := c.Bind(&req)
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
