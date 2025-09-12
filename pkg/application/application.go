package application

import (
	"mike/config"
	"mike/utils"

	"github.com/uptrace/bun"
)

type App struct {
	Config *config.Config
	DB     *bun.DB
}

func New(cfg *config.Config) (*App, error) {
	// Initialize database connection
	db := utils.GetDatabase()
	app := &App{
		Config: cfg,
		DB:     db,
	}
	return app, nil
}
