package health

import (
	"github.com/labstack/echo/v4"
	lang "gitlab.snapp.ir/platform/snapp_object_store/langs/en"
	"net/http"
)

type ApplicationHealth struct {
	Status string `json:"status"`
}

func HandleHealth(c echo.Context) error {
	var appHealth ApplicationHealth
	appHealth.Status = lang.ApplicationHealthy
	return c.JSON(http.StatusOK, appHealth)
}
