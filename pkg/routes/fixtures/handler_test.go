package fixtures

import (
	"context"
	"encoding/json"
	"mike/config"
	"mike/pkg/application"
	"mike/pkg/routes/fixtures/models"
	"mike/utils"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
)

func Test_GetFixturesByTimeRange(t *testing.T) {
	cfg := config.GetConfig()
	app, err := application.New(cfg)
	assert.NoError(t, err, "Failed to create application")
	app.DB = utils.GetDatabase()

	e := echo.New()

	// Add app to context
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("app", app)
			return next(c)
		}
	})

	RegisterRoutes(e, app)

	// Clean up any existing test data
	app.DB.NewDelete().Model((*models.Fixture)(nil)).Where("sport_id = ?", 1).Exec(context.Background())
	app.DB.NewDelete().Model((*struct {
		bun.BaseModel `bun:"teams"`
		ID            int `bun:"id,pk,autoincrement"`
	})(nil)).Where("sport_id = ?", 1).Exec(context.Background())
	app.DB.NewDelete().Model((*struct {
		bun.BaseModel `bun:"sports"`
		ID            int `bun:"id,pk,autoincrement"`
	})(nil)).Where("id = ?", 1).Exec(context.Background())

	// Create test data
	sport := struct {
		bun.BaseModel `bun:"sports"`
		ID            int    `bun:"id,pk,autoincrement" json:"id"`
		Name          string `bun:"name" json:"name"`
		Description   string `bun:"description" json:"description"`
		IsActive      bool   `bun:"is_active" json:"is_active"`
	}{
		ID:          1,
		Name:        "Test Sport",
		Description: "Test sport for fixtures",
		IsActive:    true,
	}
	_, err = app.DB.NewInsert().Model(&sport).Exec(context.Background())
	assert.NoError(t, err, "Error inserting sport")

	teams := []struct {
		bun.BaseModel `bun:"teams"`
		ID            int    `bun:"id,pk,autoincrement" json:"id"`
		Name          string `bun:"name" json:"name"`
		SportId       int    `bun:"sport_id" json:"sport_id"`
		IsActive      bool   `bun:"is_active" json:"is_active"`
	}{
		{ID: 1, Name: "Team 1", SportId: 1, IsActive: true},
		{ID: 2, Name: "Team 2", SportId: 1, IsActive: true},
	}

	for _, team := range teams {
		_, err = app.DB.NewInsert().Model(&team).Exec(context.Background())
		assert.NoError(t, err, "Error inserting team")
	}

	// Create test fixtures
	fixtures := []models.Fixture{
		{
			SportID:  1,
			TeamID1:  1,
			TeamID2:  2,
			DateTime: time.Now(),
			Status:   "scheduled",
		},
	}

	for _, fixture := range fixtures {
		_, err := app.DB.NewInsert().Model(&fixture).Exec(context.Background())
		assert.NoError(t, err, "Error inserting fixture")
	}

	tests := []struct {
		name           string
		startTime      string
		endTime        string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "valid time range",
			startTime:      time.Now().Add(-time.Hour).Format(time.RFC3339),
			endTime:        time.Now().Add(time.Hour).Format(time.RFC3339),
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:           "invalid start time",
			startTime:      "invalid-time",
			endTime:        time.Now().Add(time.Hour).Format(time.RFC3339),
			expectedStatus: http.StatusBadRequest,
			expectedCount:  0,
		},
		{
			name:           "invalid end time",
			startTime:      time.Now().Add(-time.Hour).Format(time.RFC3339),
			endTime:        "invalid-time",
			expectedStatus: http.StatusBadRequest,
			expectedCount:  0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/v1/fixtures/"+test.startTime+"/"+test.endTime, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/v1/fixtures/:startTime/:endTime")
			c.SetParamNames("startTime", "endTime")
			c.SetParamValues(test.startTime, test.endTime)
			c.Set("app", app)

			err := GetFixturesByTimeRange(c)
			assert.NoError(t, err, "Handler should not return error")

			assert.Equal(t, test.expectedStatus, rec.Code, "HTTP status code should match expected")

			if test.expectedStatus == http.StatusOK {
				var result []models.Fixture
				err = json.Unmarshal(rec.Body.Bytes(), &result)
				assert.NoError(t, err, "Response should be valid JSON")
				assert.Len(t, result, test.expectedCount, "Number of fixtures should match expected")
			}
		})
	}
}
