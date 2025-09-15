package models

import "github.com/uptrace/bun"

type Sport struct {
	bun.BaseModel `bun:"sports"`
	ID            int    `bun:"id,pk,autoincrement" json:"id"`
	Name          string `bun:"name" json:"name"`
	Description   string `bun:"description" json:"description"`
	ImageURL      string `bun:"image_url" json:"image_url"`
	IsActive      bool   `bun:"is_active" json:"is_active"`
}
