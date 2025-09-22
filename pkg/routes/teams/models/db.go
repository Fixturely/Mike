// Generic teams model for all sports
package models

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun"
)

type Team struct {
	bun.BaseModel `bun:"teams"`
	ID            int    `bun:"id,pk,autoincrement" json:"id"`
	Name          string `bun:"name" json:"name"`
	SportId       int    `bun:"sport_id" json:"sport_id"`
	Description   string `bun:"description" json:"description"`
	ImageURL      string `bun:"image_url" json:"image_url"`
	IsActive      bool   `bun:"is_active" json:"is_active"`
	ApiID         int    `bun:"api_id" json:"api_id"`
}

func CheckTeamExists(db *bun.DB, teamId int) (bool, error) {
	var team Team
	err := db.NewSelect().Model(&team).Where("id = ?", teamId).Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func GetTeamDetails(db *bun.DB, teamId int) (Team, error) {
	var team Team
	err := db.NewSelect().Model(&team).Where("id = ?", teamId).Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Team{}, errors.New("team not found")
		}
		return Team{}, err
	}
	return team, nil
}
