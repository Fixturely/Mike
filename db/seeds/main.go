package main

import (
	"context"
	"log"
	"mike/pkg/routes/sports/models"
	"mike/utils"
	"time"

	"github.com/uptrace/bun"
)

func getDatabase() *bun.DB {
	return utils.GetDatabase()
}

func main() {
	log.Println("Starting database seeding...")
	db := getDatabase()

	// Run all seed functions
	SeedSports(db)
	SeedTeams(db)
	SeedFixtures(db)

	log.Println("All seeding completed successfully!")
}

func SeedSports(db *bun.DB) {
	log.Println("Seeding Sports")
	sports := []string{
		"Baseball",
		"Basketball",
		"Football",
		"Hockey",
		"Soccer",
	}

	for _, sport := range sports {
		_, err := db.NewInsert().Model(&models.Sport{Name: sport}).Exec(context.Background())
		if err != nil {
			log.Printf("Error inserting sport %s: %v", sport, err)
		}
	}
	log.Println("Sports seeding completed")
}

func SeedTeams(db *bun.DB) {
	log.Println("Seeding Teams")
	// Sample teams for each sport
	teams := []struct {
		Name    string
		SportID int
	}{
		// Baseball teams (Sport ID 1)
		{"New York Yankees", 1},
		{"Boston Red Sox", 1},
		{"Los Angeles Dodgers", 1},
		{"Chicago Cubs", 1},
		{"San Francisco Giants", 1},

		// Basketball teams (Sport ID 2)
		{"Los Angeles Lakers", 2},
		{"Boston Celtics", 2},
		{"Golden State Warriors", 2},
		{"Chicago Bulls", 2},
		{"Miami Heat", 2},

		// Football teams (Sport ID 3)
		{"New England Patriots", 3},
		{"Kansas City Chiefs", 3},
		{"Green Bay Packers", 3},
		{"Dallas Cowboys", 3},
		{"Pittsburgh Steelers", 3},

		// Hockey teams (Sport ID 4)
		{"Boston Bruins", 4},
		{"Toronto Maple Leafs", 4},
		{"Montreal Canadiens", 4},
		{"New York Rangers", 4},
		{"Detroit Red Wings", 4},

		// Soccer teams (Sport ID 5)
		{"Manchester United", 5},
		{"Real Madrid", 5},
		{"Barcelona", 5},
		{"Bayern Munich", 5},
		{"Liverpool", 5},
	}

	for _, team := range teams {
		_, err := db.NewInsert().Model(&Team{
			Name:    team.Name,
			SportID: team.SportID,
		}).Exec(context.Background())
		if err != nil {
			log.Printf("Error inserting team %s: %v", team.Name, err)
		}
	}
	log.Println("Teams seeding completed")
}

func SeedFixtures(db *bun.DB) {
	log.Println("Seeding Fixtures")

	// Sample fixtures with realistic matchups
	fixtures := []struct {
		SportID  int
		TeamID1  int
		TeamID2  int
		DateTime time.Time
		Details  map[string]interface{}
		Status   string
	}{
		// Baseball fixtures (Teams 1-5)
		{1, 1, 2, time.Date(2024, 3, 15, 19, 0, 0, 0, time.UTC), map[string]interface{}{"venue": "Yankee Stadium", "season": "2024"}, "scheduled"},
		{1, 3, 4, time.Date(2024, 3, 16, 20, 0, 0, 0, time.UTC), map[string]interface{}{"venue": "Dodger Stadium", "season": "2024"}, "scheduled"},
		{1, 5, 1, time.Date(2024, 3, 17, 19, 30, 0, 0, time.UTC), map[string]interface{}{"venue": "Oracle Park", "season": "2024"}, "scheduled"},

		// Basketball fixtures (Teams 6-10)
		{2, 6, 7, time.Date(2024, 3, 20, 20, 0, 0, 0, time.UTC), map[string]interface{}{"venue": "Crypto.com Arena", "season": "2024"}, "scheduled"},
		{2, 8, 9, time.Date(2024, 3, 21, 19, 30, 0, 0, time.UTC), map[string]interface{}{"venue": "Chase Center", "season": "2024"}, "scheduled"},
		{2, 10, 6, time.Date(2024, 3, 22, 20, 0, 0, 0, time.UTC), map[string]interface{}{"venue": "FTX Arena", "season": "2024"}, "scheduled"},

		// Football fixtures (Teams 11-15)
		{3, 11, 12, time.Date(2024, 3, 25, 13, 0, 0, 0, time.UTC), map[string]interface{}{"venue": "Gillette Stadium", "season": "2024"}, "scheduled"},
		{3, 13, 14, time.Date(2024, 3, 26, 16, 0, 0, 0, time.UTC), map[string]interface{}{"venue": "Lambeau Field", "season": "2024"}, "scheduled"},
		{3, 15, 11, time.Date(2024, 3, 27, 20, 0, 0, 0, time.UTC), map[string]interface{}{"venue": "AT&T Stadium", "season": "2024"}, "scheduled"},

		// Hockey fixtures (Teams 16-20)
		{4, 16, 17, time.Date(2024, 3, 30, 19, 0, 0, 0, time.UTC), map[string]interface{}{"venue": "TD Garden", "season": "2024"}, "scheduled"},
		{4, 18, 19, time.Date(2024, 3, 31, 19, 30, 0, 0, time.UTC), map[string]interface{}{"venue": "Bell Centre", "season": "2024"}, "scheduled"},
		{4, 20, 16, time.Date(2024, 4, 1, 19, 0, 0, 0, time.UTC), map[string]interface{}{"venue": "Madison Square Garden", "season": "2024"}, "scheduled"},

		// Soccer fixtures (Teams 21-25)
		{5, 21, 22, time.Date(2024, 4, 5, 15, 0, 0, 0, time.UTC), map[string]interface{}{"venue": "Old Trafford", "season": "2024"}, "scheduled"},
		{5, 23, 24, time.Date(2024, 4, 6, 16, 0, 0, 0, time.UTC), map[string]interface{}{"venue": "Camp Nou", "season": "2024"}, "scheduled"},
		{5, 25, 21, time.Date(2024, 4, 7, 14, 30, 0, 0, time.UTC), map[string]interface{}{"venue": "Allianz Arena", "season": "2024"}, "scheduled"},
	}

	for _, fixture := range fixtures {
		_, err := db.NewInsert().Model(&Fixture{
			SportID:  fixture.SportID,
			TeamID1:  fixture.TeamID1,
			TeamID2:  fixture.TeamID2,
			DateTime: fixture.DateTime,
			Details:  fixture.Details,
			Status:   fixture.Status,
		}).Exec(context.Background())
		if err != nil {
			log.Printf("Error inserting fixture: %v", err)
		}
	}
	log.Println("Fixtures seeding completed")
}

// Team model for seeding
type Team struct {
	ID          int    `bun:"id,pk,autoincrement" json:"id"`
	Name        string `bun:"name" json:"name"`
	SportID     int    `bun:"sport_id" json:"sport_id"`
	Description string `bun:"description" json:"description"`
	ImageURL    string `bun:"image_url" json:"image_url"`
	IsActive    bool   `bun:"is_active" json:"is_active"`
}

// Fixture model for seeding
type Fixture struct {
	ID       int                    `bun:"id,pk,autoincrement" json:"id"`
	SportID  int                    `bun:"sport_id" json:"sport_id"`
	TeamID1  int                    `bun:"team_id_1" json:"team_id_1"`
	TeamID2  int                    `bun:"team_id_2" json:"team_id_2"`
	DateTime time.Time              `bun:"date_time" json:"date_time"`
	Details  map[string]interface{} `bun:"details" json:"details"`
	Status   string                 `bun:"status" json:"status"`
}
