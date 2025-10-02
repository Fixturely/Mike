package models

import (
	"context"
	"mike/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
)

func Test_GetFixturesByTimeRange(t *testing.T) {
	db := utils.GetDatabase()

	// Clean up any existing test data
	db.NewDelete().Model((*Fixture)(nil)).Where("sport_id = ?", 1).Exec(context.Background())
	db.NewDelete().Model((*struct {
		bun.BaseModel `bun:"teams"`
		ID            int `bun:"id,pk,autoincrement"`
	})(nil)).Where("sport_id = ?", 1).Exec(context.Background())
	db.NewDelete().Model((*struct {
		bun.BaseModel `bun:"sports"`
		ID            int `bun:"id,pk,autoincrement"`
	})(nil)).Where("id = ?", 1).Exec(context.Background())

	// Create test data - first create sports and teams
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

	// Insert test fixtures
	fixtures := []Fixture{
		{
			SportID:  1,
			TeamID1:  1,
			TeamID2:  2,
			DateTime: time.Now(), // Current time
			Status:   "scheduled",
		},
		{
			SportID:  1,
			TeamID1:  2,
			TeamID2:  3,
			DateTime: time.Now().Add(time.Hour * 24), // 24 hours from now
			Status:   "scheduled",
		},
		{
			SportID:  1,
			TeamID1:  3,
			TeamID2:  4,
			DateTime: time.Now().Add(-time.Hour * 48), // 48 hours behind current time
			Status:   "scheduled",
		},
	}

	// Insert test fixtures
	for _, fixture := range fixtures {
		_, err := db.NewInsert().Model(&fixture).Exec(context.Background())
		assert.NoError(t, err)
	}

	// Test cases
	tests := []struct {
		name                 string
		startTime            time.Time
		endTime              time.Time
		expectedFixtureCount int
		wantError            bool
	}{
		{
			name:                 "test_get_fixtures_by_time_range",
			startTime:            time.Now(),
			endTime:              time.Now().Add(time.Hour * 24),
			expectedFixtureCount: 1,
			wantError:            false,
		},
		{
			name:                 "test_get_fixtures_by_time_range_2",
			startTime:            time.Now().Add(time.Hour * 24),
			endTime:              time.Now().Add(time.Hour * 48),
			expectedFixtureCount: 0,
			wantError:            false,
		},
		{
			name:                 "test_get_fixtures_by_time_range_3",
			startTime:            time.Now().Add(-time.Hour * 48),
			endTime:              time.Now(),
			expectedFixtureCount: 1,
			wantError:            false,
		},
		{
			name:                 "test_get_fixtures_by_time_range_4",
			startTime:            time.Now().Add(-time.Hour * 48),
			endTime:              time.Now().Add(time.Hour * 24),
			expectedFixtureCount: 2,
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
