package health

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/snapp-incubator/S3-Panel/internal/messages"
)

type ApplicationHealth struct {
	Status string `json:"status"`
}

func HandleHealth(c echo.Context) error {
	var appHealth ApplicationHealth
	appHealth.Status = messages.ApplicationHealthy
	return c.JSON(http.StatusOK, appHealth)
}
