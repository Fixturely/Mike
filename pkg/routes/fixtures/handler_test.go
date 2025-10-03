package fixtures

import (
	"bytes"
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

func setupTestApp(t *testing.T) (*application.App, *echo.Echo) {
	cfg := config.GetConfig()
	app, err := application.New(cfg)
	assert.NoError(t, err, "Failed to create application")
	app.DB = utils.GetDatabase()

	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("app", app)
			return next(c)
		}
	})
	RegisterRoutes(e, app)
	return app, e
}

// setupTestData creates test data and returns a cleanup function
func setupTestData(t *testing.T, db *bun.DB) func() {
	// Clean up any existing test data first
	_, _ = db.NewDelete().Model((*models.Fixture)(nil)).Where("sport_id IN (SELECT id FROM sports WHERE name = ?)", "Test Sport Handler").Exec(context.Background())
	_, _ = db.NewDelete().Model((*struct {
		bun.BaseModel `bun:"teams"`
		ID            int `bun:"id,pk,autoincrement"`
	})(nil)).Where("sport_id IN (SELECT id FROM sports WHERE name = ?)", "Test Sport Handler").Exec(context.Background())
	_, _ = db.NewDelete().Model((*struct {
		bun.BaseModel `bun:"sports"`
		ID            int `bun:"id,pk,autoincrement"`
	})(nil)).Where("name = ?", "Test Sport Handler").Exec(context.Background())

	// Create sport without specifying ID to let auto-increment handle it
	sport := struct {
		bun.BaseModel `bun:"sports"`
		Name          string `bun:"name" json:"name"`
		Description   string `bun:"description" json:"description"`
		IsActive      bool   `bun:"is_active" json:"is_active"`
	}{
		Name:        "Test Sport Handler",
		Description: "Test sport for fixtures handler tests",
		IsActive:    true,
	}
	_, err := db.NewInsert().Model(&sport).Exec(context.Background())
	assert.NoError(t, err, "Error inserting sport")

	// Get the inserted sport ID
	var sportID int
	err = db.NewSelect().Model((*struct {
		bun.BaseModel `bun:"sports"`
		ID            int `bun:"id,pk,autoincrement"`
	})(nil)).Where("name = ?", "Test Sport Handler").Scan(context.Background(), &sportID)
	assert.NoError(t, err)

	// Create teams without specifying IDs
	teams := []struct {
		bun.BaseModel `bun:"teams"`
		Name          string `bun:"name" json:"name"`
		SportId       int    `bun:"sport_id" json:"sport_id"`
		IsActive      bool   `bun:"is_active" json:"is_active"`
	}{
		{Name: "Team 1 Handler", SportId: sportID, IsActive: true},
		{Name: "Team 2 Handler", SportId: sportID, IsActive: true},
		{Name: "Team 3 Handler", SportId: sportID, IsActive: true},
		{Name: "Team 4 Handler", SportId: sportID, IsActive: true},
	}

	var teamIDs []int
	for _, team := range teams {
		_, err = db.NewInsert().Model(&team).Exec(context.Background())
		assert.NoError(t, err, "Error inserting team")

		var teamID int
		err = db.NewSelect().Model((*struct {
			bun.BaseModel `bun:"teams"`
			ID            int `bun:"id,pk,autoincrement"`
		})(nil)).Where("name = ? AND sport_id = ?", team.Name, sportID).Scan(context.Background(), &teamID)
		assert.NoError(t, err)
		teamIDs = append(teamIDs, teamID)
	}

	// Use deterministic base time for consistent testing
	baseTime := time.Unix(1609459200, 0).UTC() // Jan 1, 2021 00:00:00 UTC

	// Create test fixtures with deterministic times and dynamic IDs
	fixtures := []models.Fixture{
		{
			SportID:  sportID,
			TeamID1:  teamIDs[0],
			TeamID2:  teamIDs[1],
			DateTime: baseTime, // Base time (2021-01-01 00:00:00 UTC)
			Status:   "scheduled",
		},
		{
			SportID:  sportID,
			TeamID1:  teamIDs[1],
			TeamID2:  teamIDs[2],
			DateTime: baseTime.Add(24 * time.Hour), // Base time + 24 hours (2021-01-02 00:00:00 UTC)
			Status:   "scheduled",
		},
		{
			SportID:  sportID,
			TeamID1:  teamIDs[2],
			TeamID2:  teamIDs[3],
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
		_, _ = db.NewDelete().Model((*models.Fixture)(nil)).Where("sport_id IN (SELECT id FROM sports WHERE name = ?)", "Test Sport Handler").Exec(context.Background())
		_, _ = db.NewDelete().Model((*struct {
			bun.BaseModel `bun:"teams"`
			ID            int `bun:"id,pk,autoincrement"`
		})(nil)).Where("sport_id IN (SELECT id FROM sports WHERE name = ?)", "Test Sport Handler").Exec(context.Background())
		_, _ = db.NewDelete().Model((*struct {
			bun.BaseModel `bun:"sports"`
			ID            int `bun:"id,pk,autoincrement"`
		})(nil)).Where("name = ?", "Test Sport Handler").Exec(context.Background())
	}
}

func Test_GetFixturesByTimeRange(t *testing.T) {
	app, e := setupTestApp(t)

	// Use deterministic base time for consistent testing
	baseTime := time.Unix(1609459200, 0).UTC() // Jan 1, 2021 00:00:00 UTC

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
			// Create JSON request body
			requestBody := TimeRangeRequest{
				Start: test.startTime,
				End:   test.endTime,
			}

			jsonBody, err := json.Marshal(requestBody)
			assert.NoError(t, err, "Should be able to marshal request body")

			req := httptest.NewRequest(http.MethodPost, "/v1/fixtures/daterange", bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/v1/fixtures/daterange")
			c.Set("app", app)

			err = GetFixturesByTimeRange(c)
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

func Test_GetFixturesByTimeRange_AcceptsDifferentTimeFormats(t *testing.T) {
	app, e := setupTestApp(t)

	tests := []struct {
		name        string
		startTime   string
		endTime     string
		wantError   bool
		description string
	}{
		{
			name:        "RFC3339 format",
			startTime:   "2021-01-01T00:00:00Z",
			endTime:     "2021-01-02T00:00:00Z",
			wantError:   false,
			description: "accepts RFC3339 format",
		},
		{
			name:        "PostgreSQL timestamp format",
			startTime:   "2021-01-01 00:00:00.000000",
			endTime:     "2021-01-02 00:00:00.000000",
			wantError:   false,
			description: "accepts PostgreSQL timestamp format",
		},
		{
			name:        "Simple datetime format",
			startTime:   "2021-01-01 00:00:00",
			endTime:     "2021-01-02 00:00:00",
			wantError:   false,
			description: "accepts simple datetime format",
		},
		{
			name:        "Date only format",
			startTime:   "2021-01-01",
			endTime:     "2021-01-02",
			wantError:   false,
			description: "accepts date only format",
		},
		{
			name:        "Invalid format",
			startTime:   "invalid-date",
			endTime:     "2021-01-02",
			wantError:   true,
			description: "rejects invalid format",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create JSON request body
			requestBody := TimeRangeRequest{
				Start: test.startTime,
				End:   test.endTime,
			}

			jsonBody, err := json.Marshal(requestBody)
			assert.NoError(t, err, "Should be able to marshal request body")

			req := httptest.NewRequest(http.MethodPost, "/v1/fixtures/daterange", bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/v1/fixtures/daterange")
			c.Set("app", app)

			err = GetFixturesByTimeRange(c)

			if test.wantError {
				assert.Equal(t, http.StatusBadRequest, rec.Code, test.description)
			} else {
				assert.NoError(t, err, test.description)
				assert.Equal(t, http.StatusOK, rec.Code, test.description)
			}
		})
	}
}
