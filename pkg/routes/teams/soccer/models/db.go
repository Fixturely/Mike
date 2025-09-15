package models

import (
	"context"
	"mike/pkg/routes/teams/models"

	"github.com/uptrace/bun"
)

const SoccerSportId = 5

// API Response structures for teams belonging to a league and season for soccer
type TeamsResponse struct {
	Get        string `json:"get"`
	Parameters struct {
		League string `json:"league"`
		Season string `json:"season"`
	} `json:"parameters"`
	Errors  []interface{} `json:"errors"`
	Results int           `json:"results"`
	Paging  struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"paging"`
	Response []TeamData `json:"response"`
}

type TeamData struct {
	Team struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Code     string `json:"code"`
		Country  string `json:"country"`
		Founded  int    `json:"founded"`
		National bool   `json:"national"`
		Logo     string `json:"logo"`
	} `json:"team"`
	Venue struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Address  string `json:"address"`
		City     string `json:"city"`
		Capacity int    `json:"capacity"`
		Surface  string `json:"surface"`
		Image    string `json:"image"`
	} `json:"venue"`
}

func (t *TeamData) ToTeam() *models.Team {
	return &models.Team{
		Name:        t.Team.Name,
		SportId:     SoccerSportId,
		Description: "", // API doesn't provide description
		ImageURL:    t.Team.Logo,
		IsActive:    true,
		ApiID:       t.Team.ID,
	}
}

func (tr *TeamsResponse) ToTeams() []models.Team {
	teams := make([]models.Team, 0, len(tr.Response))
	for _, teamData := range tr.Response {
		teams = append(teams, *teamData.ToTeam())
	}
	return teams
}

func InsertTeams(db *bun.DB, teams []models.Team) error {
	if len(teams) == 0 {
		return nil
	}

	_, err := db.NewInsert().
		Model(&teams).
		On("CONFLICT (name, sport_id) DO NOTHING").
		Exec(context.Background())
	return err
}
