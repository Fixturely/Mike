package server

import (
	"fmt"
	"net/http"
	"time"

	"mike/pkg/application"
	"mike/pkg/routes/health"

	"github.com/labstack/echo/v4"
)

func CreateEcho() *echo.Echo {
	e := echo.New()

	return e
}

func RegisterRoutes(app *application.App, e *echo.Echo) error {
	health.RegisterRoutes(e, app)
	return nil
}

// New creates an instance of the HTTP server for an Application
func New(app *application.App) (*http.Server, error) {
	e := CreateEcho()
	if err := RegisterRoutes(app, e); err != nil {
		return nil, fmt.Errorf("failed to register routes: %w", err)
	}

	addr := fmt.Sprintf(":%d", app.Config.ServerPort)
	// In development, bind to all interfaces so it's accessible from host machine
	if app.Config.Environment == "development" {
		addr = "0.0.0.0" + addr
	}
	srv := &http.Server{
		Addr:         addr,
		Handler:      e,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return srv, nil
}
