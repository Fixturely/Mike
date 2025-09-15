package models

import (
	"context"
	"log"
	fixtureModels "mike/pkg/routes/fixtures/models"
	"mike/pkg/routes/teams/models"
	"mike/utils"
	"time"

	"github.com/uptrace/bun"
)

// Raw API response structures for fixtures belonging to a league and season for soccer
type FixturesResponse struct {
	Get        string            `json:"get"`
	Parameters map[string]string `json:"parameters"`
	Errors     []interface{}     `json:"errors"`
	Results    int               `json:"results"`
	Paging     struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"paging"`
	Response []FixtureItem `json:"response"`
}

type FixtureItem struct {
	Fixture struct {
		ID        int    `json:"id"`
		Timezone  string `json:"timezone"`
		Date      string `json:"date"` // ISO8601
		Timestamp int64  `json:"timestamp"`
		Status    struct {
			Long    string      `json:"long"`
			Short   string      `json:"short"`
			Elapsed *int        `json:"elapsed"`
			Extra   interface{} `json:"extra"`
		} `json:"status"`
	} `json:"fixture"`
	League struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Season int    `json:"season"`
	} `json:"league"`
	Teams struct {
		Home struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"home"`
		Away struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"away"`
	} `json:"teams"`
}

func getTeamDetails(teamApiId int) (models.Team, error) {
	db := utils.GetDatabase()

	var team models.Team
	err := db.NewSelect().Model(&team).Where("api_id = ?", teamApiId).Scan(context.Background())
	if err != nil {
		return models.Team{}, err
	}
	return team, nil
}

func (fixture FixtureItem) ToFixture() fixtureModels.Fixture {
	teamId1, err := getTeamDetails(fixture.Teams.Home.ID)
	if err != nil {
		log.Fatalf("Error getting team details: %v", err)
	}
	teamId2, err := getTeamDetails(fixture.Teams.Away.ID)
	if err != nil {
		log.Fatalf("Error getting team details: %v", err)
	}
	date, err := time.Parse(time.RFC3339, fixture.Fixture.Date)
	if err != nil {
		log.Fatalf("Error parsing date: %v", err)
	}
	return fixtureModels.Fixture{
		SportID:  5, // Soccer
		TeamID1:  int(teamId1.ID),
		TeamID2:  int(teamId2.ID),
		DateTime: date,
		Status:   fixture.Fixture.Status.Short,
		Details: fixtureModels.Details{
			HomeTeam: teamId1.Name,
			AwayTeam: teamId2.Name,
			DateTime: date,
			Status:   fixture.Fixture.Status.Short,
		},
	}
}

func (fixturesResp FixturesResponse) ToFixtures() []fixtureModels.Fixture {
	fixtures := make([]fixtureModels.Fixture, 0, len(fixturesResp.Response))
	for _, fixture := range fixturesResp.Response {
		fixtures = append(fixtures, fixture.ToFixture())
	}
	return fixtures
}

func InsertFixtures(db *bun.DB, fixtures []fixtureModels.Fixture) error {
	if len(fixtures) == 0 {
		return nil
	}
	_, err := db.NewInsert().
		Model(&fixtures).
		On("CONFLICT (sport_id, team_id_1, team_id_2, date_time) DO NOTHING").
		Exec(context.Background())
	return err
}
