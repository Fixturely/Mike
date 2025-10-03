package fixtures

import (
	"mike/pkg/application"
	"mike/pkg/routes/fixtures/models"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// TimeRangeRequest represents the JSON request body for date range queries
type TimeRangeRequest struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

func GetFixturesByTimeRange(c echo.Context) error {
	app := c.Get("app").(*application.App)

	// Parse JSON request body
	var req TimeRangeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON in request body"})
	}

	// Validate request fields
	if req.Start == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing required field: start"})
	}
	if req.End == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing required field: end"})
	}

	startTime := req.Start
	endTime := req.End

	// Try multiple date/time formats
	timeFormats := []string{
		time.RFC3339,                  // "2006-01-02T15:04:05Z07:00"
		time.RFC3339Nano,              // "2006-01-02T15:04:05.999999999Z07:00"
		"2006-01-02T15:04:05Z",        // "2006-01-02T15:04:05Z"
		"2006-01-02T15:04:05.999999Z", // "2006-01-02T15:04:05.999999Z"
		"2006-01-02 15:04:05.999999",  // "2006-01-02 15:04:05.999999" (postgres format)
		"2006-01-02 15:04:05",         // "2006-01-02 15:04:05"
		"2006-01-02",                  // "2006-01-02" (date only)
	}

	var startTimeTime time.Time
	var endTimeTime time.Time

	// Parse start time
	for _, format := range timeFormats {
		if t, err := time.Parse(format, startTime); err == nil {
			startTimeTime = t
			break
		}
	}
	if startTimeTime.IsZero() {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid start time format"})
	}

	// Parse end time
	for _, format := range timeFormats {
		if t, err := time.Parse(format, endTime); err == nil {
			endTimeTime = t
			break
		}
	}
	if endTimeTime.IsZero() {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid end time format"})
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
