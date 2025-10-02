package models

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

func GetFixturesByTimeRange(db *bun.DB, startTime time.Time, endTime time.Time) ([]Fixture, error) {
	var fixtures []Fixture
	err := db.NewSelect().Model(&fixtures).Where("date_time >= ? AND date_time <= ?", startTime, endTime).Scan(context.Background())
	if err != nil {
		return nil, err
	}
	return fixtures, nil
}

type Fixture struct {
	ID       int       `bun:"id,pk,autoincrement" json:"id"`
	SportID  int       `bun:"sport_id" json:"sport_id"`
	TeamID1  int       `bun:"team_id_1" json:"team_id_1"`
	TeamID2  int       `bun:"team_id_2" json:"team_id_2"`
	DateTime time.Time `bun:"date_time" json:"date_time"`
	Details  Details   `bun:"details" json:"details"`
	Status   string    `bun:"status" json:"status"`
}

type Details struct {
	HomeTeam string    `json:"home_team"`
	AwayTeam string    `json:"away_team"`
	DateTime time.Time `json:"date_time"`
	Status   string    `json:"status"`
}
