package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mike/config"
	"mike/pkg/routes/teams/soccer/models"
	"mike/utils"
	"net/http"
	"time"
)

func fetchTeams() {
	cfg := config.GetConfig()

	db := utils.GetDatabase()

	// Premier League ID 39 for 2024 season
	leagueID := 39
	season := "2025"

	url := fmt.Sprintf("%s/teams?league=%d&season=%s", cfg.API.FootballAPIURL, leagueID, season)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

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

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("API request failed with status: %d, response: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	// Parse and pretty print JSON
	var teamsResponse models.TeamsResponse
	if err := json.Unmarshal(body, &teamsResponse); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Convert to team object
	teams := teamsResponse.ToTeams()

	// Save to database
	err = models.InsertTeams(db, teams)
	if err != nil {
		log.Fatalf("Error saving teams to database: %v", err)
	}
}

func main() {
	fetchTeams()
}
