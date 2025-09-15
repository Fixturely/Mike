package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mike/config"
	"mike/pkg/routes/fixtures/soccer/models"
	"mike/utils"
	"net/http"
	"time"
)

// Cleaned shape aligned with our fixtures schema (without DB IDs yet)
type CleanFixture struct {
	SportID  int             `json:"sport_id"`
	HomeTeam string          `json:"home_team"`
	AwayTeam string          `json:"away_team"`
	DateTime time.Time       `json:"date_time"`
	Status   string          `json:"status"`
	Details  json.RawMessage `json:"details"`
}

func fetchFixtures() {
	cfg := config.GetConfig()
	db := utils.GetDatabase()

	leagueID := 39
	season := "2025"

	url := fmt.Sprintf("%s/fixtures?league=%d&season=%s", cfg.API.FootballAPIURL, leagueID, season)

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Add("X-RapidAPI-Key", cfg.API.FootballAPIKey)
	req.Header.Add("X-RapidAPI-Host", "v3.football.api-sports.io")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Fatalf("API request failed with status: %d, response: %s", resp.StatusCode, string(body))
	}

	var fixturesResp models.FixturesResponse
	if err := json.Unmarshal(body, &fixturesResp); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	fixtures := fixturesResp.ToFixtures()

	err = models.InsertFixtures(db, fixtures)
	if err != nil {
		log.Fatalf("Error inserting fixtures: %v", err)
	}

}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func main() {
	fetchFixtures()
}
