package main

import (
	"context"

	"github.com/hippeus/appbase/pkg/demoapp"
	"github.com/hippeus/appbase/pkg/logger"
)

type App struct {
	cfg Config
	logger.Logger
	web demoapp.Service
}

func (a *App) Run(ctx context.Context) error {
	a.Info("Starting application...")
	a.web.Name = a.cfg.App.Name
	a.web.Host = a.cfg.HTTP.Host
	a.web.Port = a.cfg.HTTP.Port
	a.web.Logger = a.Logger

	return a.web.Run(ctx)
}

func main() {
	var app App
	app.cfg = getConfig()
	app.Logger = logger.New(app.cfg.Logging).WithPrefix(app.cfg.App.Name)
	ctx := context.Background()
	if err := app.Run(ctx); err != nil {
		app.Errorf("application terminated unexpectedly: %v", err)
	}
}
