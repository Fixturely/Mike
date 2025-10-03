package models

import (
	"context"
	"mike/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
)

// setupModelTestData creates test data for model tests and returns a cleanup function
func setupModelTestData(t *testing.T, db *bun.DB) func() {
	// Clean up any existing test data first
	_, _ = db.NewDelete().Model((*Fixture)(nil)).Where("sport_id = ?", 1).Exec(context.Background())
	_, _ = db.NewDelete().Model((*struct {
		bun.BaseModel `bun:"teams"`
		ID            int `bun:"id,pk,autoincrement"`
	})(nil)).Where("sport_id = ?", 1).Exec(context.Background())
	_, _ = db.NewDelete().Model((*struct {
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
	assert.NoError(t, err)

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
		assert.NoError(t, err)
	}

	// Capture base time for deterministic test timing
	baseTime := time.Unix(1609459200, 0).UTC() // Jan 1, 2021 00:00:00 UTC

	// Insert test fixtures using deterministic times
	fixtures := []Fixture{
		{
			SportID:  1,
			TeamID1:  1,
			TeamID2:  2,
			DateTime: baseTime, // Base time
		},
		{
			SportID:  1,
			TeamID1:  2,
			TeamID2:  3,
			DateTime: baseTime.Add(24 * time.Hour), // Base time + 24 hours
		},
		{
			SportID:  1,
			TeamID1:  3,
			TeamID2:  4,
			DateTime: baseTime.Add(-48 * time.Hour), // Base time - 48 hours
		},
	}

	// Insert test fixtures
	for _, fixture := range fixtures {
		_, err := db.NewInsert().Model(&fixture).Exec(context.Background())
		assert.NoError(t, err)
	}

	// Return cleanup function
	return func() {
		_, _ = db.NewDelete().Model((*Fixture)(nil)).Where("sport_id = ?", 1).Exec(context.Background())
		_, _ = db.NewDelete().Model((*struct {
			bun.BaseModel `bun:"teams"`
			ID            int `bun:"id,pk,autoincrement"`
		})(nil)).Where("sport_id = ?", 1).Exec(context.Background())
		_, _ = db.NewDelete().Model((*struct {
			bun.BaseModel `bun:"sports"`
			ID            int `bun:"id,pk,autoincrement"`
		})(nil)).Where("id = ?", 1).Exec(context.Background())
	}
}

func Test_GetFixturesByTimeRange(t *testing.T) {
	db := utils.GetDatabase()

	// Setup test data and defer cleanup
	teardown := setupModelTestData(t, db)
	defer teardown()

	// Capture base time for deterministic test timing
	baseTime := time.Unix(1609459200, 0).UTC() // Jan 1, 2021 00:00:00 UTC

	// Test cases
	tests := []struct {
		name                 string
		startTime            time.Time
		endTime              time.Time
		expectedFixtureCount int
		wantError            bool
	}{
		{
			name:                 "gets fixtures from base time to base+24h",
			startTime:            baseTime,
			endTime:              baseTime.Add(time.Hour * 24),
			expectedFixtureCount: 2, // base time (2021-01-01) + base+24h (2021-01-02)
			wantError:            false,
		},
		{
			name:                 "gets fixture exactly at base + 24h",
			startTime:            baseTime.Add(time.Hour * 24),
			endTime:              baseTime.Add(time.Hour * 25), // Small range to avoid overlaps
			expectedFixtureCount: 1,                            // Only the base+24h fixture
			wantError:            false,
		},
		{
			name:                 "gets fixtures from base-48h to base time",
			startTime:            baseTime.Add(-time.Hour * 48),
			endTime:              baseTime,
			expectedFixtureCount: 2, // base-48h (2020-12-30) + base time (2021-01-01)
			wantError:            false,
		},
		{
			name:                 "gets all fixtures spanning entire range",
			startTime:            baseTime.Add(-time.Hour * 48),
			endTime:              baseTime.Add(time.Hour * 24),
			expectedFixtureCount: 3, // All three fixtures
			wantError:            false,
		},
		{
			name:                 "empty time range",
			startTime:            time.Time{},
			endTime:              time.Time{},
			expectedFixtureCount: 0,
			wantError:            false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fixtures, err := GetFixturesByTimeRange(db, test.startTime, test.endTime)
			if test.wantError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			assert.Equal(t, test.expectedFixtureCount, len(fixtures))
		})
	}
}
