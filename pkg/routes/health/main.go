package health

import (
	"mike/pkg/application"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, app *application.App) {
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
}
