package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/uptrace/bun"
)

type Sport struct {
	bun.BaseModel `bun:"sports"`
	ID            int64  `bun:"id,pk,autoincrement" json:"id"`
	Name          string `bun:"name" json:"name"`
	Description   string `bun:"description" json:"description"`
	ImageURL      string `bun:"image_url" json:"image_url"`
	IsActive      bool   `bun:"is_active" json:"is_active"`
}

func CheckSportExists(db *bun.DB, sportId int) (bool, error) {
	var sport Sport
	err := db.NewSelect().Model(&sport).Where("id = ?", sportId).Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("Sport not found - in this function")
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func GetSportDetails(db *bun.DB, sportId int) (Sport, error) {
	var sport Sport
	err := db.NewSelect().Model(&sport).Where("id = ?", sportId).Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Sport{}, errors.New("sport not found")
		}
		return Sport{}, err
	}
	return sport, nil
}

func GetAllSports(db *bun.DB) ([]Sport, error) {
	var sports []Sport
	err := db.NewSelect().Model(&sports).Where("is_active = ?", true).Scan(context.Background())
	if err != nil {
		return nil, err
	}
	return sports, nil
}
