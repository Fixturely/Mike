package sports

import (
	"mike/pkg/application"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, app *application.App) {
	e.GET("/v1/sport/exists/:id", CheckSportExists)
	e.GET("/v1/sport/details/:id", GetSportDetails)
	e.GET("/v1/sports", GetAllSports)
}
