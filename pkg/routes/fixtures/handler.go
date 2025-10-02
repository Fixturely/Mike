package fixtures

import (
	"mike/pkg/application"
	"mike/pkg/routes/fixtures/models"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func GetFixturesByTimeRange(c echo.Context) error {
	app := c.Get("app").(*application.App)
	startTime := c.Param("startTime")
	endTime := c.Param("endTime")

	startTimeTime, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid start time"})
	}
	endTimeTime, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid end time"})
	}
	fixtures, err := models.GetFixturesByTimeRange(app.DB, startTimeTime, endTimeTime)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get fixtures"})
	}
	return c.JSON(http.StatusOK, fixtures)
}
