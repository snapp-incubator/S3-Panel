package health

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	Healthy = "Healthy"
)

func HandleHealth() echo.HandlerFunc {
	return func(c echo.Context) error {
		healthData := map[string]string{
			"status": Healthy,
		}
		return c.JSON(http.StatusOK, healthData)
	}
}
