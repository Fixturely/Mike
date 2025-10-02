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
	startTime := c.QueryParam("start")
	endTime := c.QueryParam("end")

	// Validate that required query parameters are provided
	if startTime == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing required query parameter: start"})
	}
	if endTime == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing required query parameter: end"})
	}

	startTimeTime, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid start time"})
	}
	endTimeTime, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid end time"})
	}

	// Validate that start time is not after end time
	if startTimeTime.After(endTimeTime) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Start time cannot be after end time"})
	}

	fixtures, err := models.GetFixturesByTimeRange(app.DB, startTimeTime, endTimeTime)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get fixtures"})
	}
	return c.JSON(http.StatusOK, fixtures)
}
