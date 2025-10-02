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
	"net/url"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
)

// setupTestData creates test data and returns a cleanup function
func setupTestData(t *testing.T, db *bun.DB) func() {
	// Clean up any existing test data first
	db.NewDelete().Model((*models.Fixture)(nil)).Where("sport_id = ?", 1).Exec(context.Background())
	db.NewDelete().Model((*struct {
		bun.BaseModel `bun:"teams"`
		ID            int `bun:"id,pk,autoincrement"`
	})(nil)).Where("sport_id = ?", 1).Exec(context.Background())
	db.NewDelete().Model((*struct {
		bun.BaseModel `bun:"sports"`
		ID            int `bun:"id,pk,autoincrement"`
	})(nil)).Where("id = ?", 1).Exec(context.Background())

	// Create sport
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
	_, err := db.NewInsert().Model(&sport).Exec(context.Background())
	assert.NoError(t, err, "Error inserting sport")

	// Create teams
	teams := []struct {
		bun.BaseModel `bun:"teams"`
		ID            int    `bun:"id,pk,autoincrement" json:"id"`
		Name          string `bun:"name" json:"name"`
		SportId       int    `bun:"sport_id" json:"sport_id"`
		IsActive      bool   `bun:"is_active" json:"is_active"`
	}{
		{ID: 1, Name: "Team 1", SportId: 1, IsActive: true},
		{ID: 2, Name: "Team 2", SportId: 1, IsActive: true},
		{ID: 3, Name: "Team 3", SportId: 1, IsActive: true},
		{ID: 4, Name: "Team 4", SportId: 1, IsActive: true},
	}

	for _, team := range teams {
		_, err = db.NewInsert().Model(&team).Exec(context.Background())
		assert.NoError(t, err, "Error inserting team")
	}

	// Use deterministic base time for consistent testing
	baseTime := time.Unix(1609459200, 0).UTC() // Jan 1, 2021 00:00:00 UTC

	// Create test fixtures with deterministic times
	fixtures := []models.Fixture{
		{
			SportID:  1,
			TeamID1:  1,
			TeamID2:  2,
			DateTime: baseTime, // Base time (2021-01-01 00:00:00 UTC)
			Status:   "scheduled",
		},
		{
			SportID:  1,
			TeamID1:  2,
			TeamID2:  3,
			DateTime: baseTime.Add(24 * time.Hour), // Base time + 24 hours (2021-01-02 00:00:00 UTC)
			Status:   "scheduled",
		},
		{
			SportID:  1,
			TeamID1:  3,
			TeamID2:  4,
			DateTime: baseTime.Add(-48 * time.Hour), // Base time - 48 hours (2020-12-30 00:00:00 UTC)
			Status:   "scheduled",
		},
	}

	// Insert test fixtures
	for _, fixture := range fixtures {
		_, err := db.NewInsert().Model(&fixture).Exec(context.Background())
		assert.NoError(t, err, "Error inserting fixture")
	}

	// Return cleanup function
	return func() {
		db.NewDelete().Model((*models.Fixture)(nil)).Where("sport_id = ?", 1).Exec(context.Background())
		db.NewDelete().Model((*struct {
			bun.BaseModel `bun:"teams"`
			ID            int `bun:"id,pk,autoincrement"`
		})(nil)).Where("sport_id = ?", 1).Exec(context.Background())
		db.NewDelete().Model((*struct {
			bun.BaseModel `bun:"sports"`
			ID            int `bun:"id,pk,autoincrement"`
		})(nil)).Where("id = ?", 1).Exec(context.Background())
	}
}

func Test_GetFixturesByTimeRange(t *testing.T) {
	cfg := config.GetConfig()
	app, err := application.New(cfg)
	assert.NoError(t, err, "Failed to create application")
	app.DB = utils.GetDatabase()

	// Use deterministic base time for consistent testing
	baseTime := time.Unix(1609459200, 0).UTC() // Jan 1, 2021 00:00:00 UTC

	e := echo.New()

	// Add app to context
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("app", app)
			return next(c)
		}
	})

	RegisterRoutes(e, app)

	// Setup test data and defer cleanup
	teardown := setupTestData(t, app.DB)
	defer teardown()

	tests := []struct {
		name           string
		startTime      string
		endTime        string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "test_get_fixtures_by_time_range: gets fixtures from [base time, base+24h]",
			startTime:      baseTime.Format(time.RFC3339),
			endTime:        baseTime.Add(24 * time.Hour).Format(time.RFC3339),
			expectedStatus: http.StatusOK,
			expectedCount:  2, // base time + base+24h fixtures (both inclusive)
		},
		{
			name:           "test_get_fixtures_by_time_range_2: gets fixtures from [base-48h, base time]",
			startTime:      baseTime.Add(-48 * time.Hour).Format(time.RFC3339),
			endTime:        baseTime.Format(time.RFC3339),
			expectedStatus: http.StatusOK,
			expectedCount:  2, // base-48h + base time fixtures (both inclusive)
		},
		{
			name:           "test_get_fixtures_by_time_range_3: gets all fixtures from [base-48h, base+24h]",
			startTime:      baseTime.Add(-48 * time.Hour).Format(time.RFC3339),
			endTime:        baseTime.Add(24 * time.Hour).Format(time.RFC3339),
			expectedStatus: http.StatusOK,
			expectedCount:  3, // all three fixtures (inclusive range)
		},
		{
			name:           "start time after end time",
			startTime:      baseTime.Add(24 * time.Hour).Format(time.RFC3339),
			endTime:        baseTime.Format(time.RFC3339),
			expectedStatus: http.StatusBadRequest,
			expectedCount:  0,
		},
		{
			name:           "missing start parameter",
			startTime:      "",
			endTime:        baseTime.Add(24 * time.Hour).Format(time.RFC3339),
			expectedStatus: http.StatusBadRequest,
			expectedCount:  0,
		},
		{
			name:           "missing end parameter",
			startTime:      baseTime.Format(time.RFC3339),
			endTime:        "",
			expectedStatus: http.StatusBadRequest,
			expectedCount:  0,
		},
		{
			name:           "gets fixture exactly at base + 24h with small range",
			startTime:      baseTime.Add(24 * time.Hour).Format(time.RFC3339),
			endTime:        baseTime.Add(25 * time.Hour).Format(time.RFC3339),
			expectedStatus: http.StatusOK,
			expectedCount:  1, // only the base+24h fixture
		},
		{
			name:           "time range with no fixtures",
			startTime:      baseTime.Add(100 * time.Hour).Format(time.RFC3339),
			endTime:        baseTime.Add(101 * time.Hour).Format(time.RFC3339),
			expectedStatus: http.StatusOK,
			expectedCount:  0, // no fixtures in this range
		},
		{
			name:           "invalid start time",
			startTime:      "invalid-time",
			endTime:        baseTime.Add(24 * time.Hour).Format(time.RFC3339),
			expectedStatus: http.StatusBadRequest,
			expectedCount:  0,
		},
		{
			name:           "invalid end time",
			startTime:      baseTime.Format(time.RFC3339),
			endTime:        "invalid-time",
			expectedStatus: http.StatusBadRequest,
			expectedCount:  0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Build URL with proper query parameter encoding
			params := url.Values{}
			if test.startTime != "" {
				params.Add("start", test.startTime)
			}
			if test.endTime != "" {
				params.Add("end", test.endTime)
			}

			requestURL := "/v1/fixtures"
			if len(params) > 0 {
				requestURL += "?" + params.Encode()
			}

			req := httptest.NewRequest(http.MethodGet, requestURL, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/v1/fixtures")
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
