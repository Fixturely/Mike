// Generic teams model for all sports
package models

import "github.com/uptrace/bun"

type Team struct {
	bun.BaseModel `bun:"teams"`
	ID            int64  `bun:"id,pk,autoincrement" json:"id"`
	Name          string `bun:"name" json:"name"`
	SportId       int    `bun:"sport_id" json:"sport_id"`
	Description   string `bun:"description" json:"description"`
	ImageURL      string `bun:"image_url" json:"image_url"`
	IsActive      bool   `bun:"is_active" json:"is_active"`
	ApiID         int    `bun:"api_id" json:"api_id"`
}
