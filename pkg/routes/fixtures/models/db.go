package models

import "time"

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
