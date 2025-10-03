package fixtures

import (
	"mike/pkg/application"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, app *application.App) {
	e.POST("/v1/fixtures/daterange", GetFixturesByTimeRange)
}
