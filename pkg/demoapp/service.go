package demoapp

import (
	"context"
	"fmt"
	"sync"

	restv1 "github.com/hippeus/appbase/pkg/demoapp/rest"
	"github.com/hippeus/appbase/pkg/logger"
	"github.com/hippeus/appbase/pkg/spa"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Service struct {
	Name   string
	Host   string
	Port   uint32
	Logger logger.Logger

	once sync.Once
	e    *echo.Echo
	ui   *spa.SPA
}

func (svc *Service) Run(ctx context.Context) error {
	svc.once.Do(func() {
		svc.e = echo.New()
		if svc.Host == "" {
			svc.Host = "0.0.0.0"
		}
		if svc.Port == 0 {
			svc.Port = 80
		}
		if svc.Name == "" {
			svc.Name = "webapp"
		}
	})

	svc.e.HideBanner = true

	// Top level middlewares
	mw := []echo.MiddlewareFunc{
		middleware.NonWWWRedirect(),
		middleware.Recover(),
	}

	if svc.Logger != nil {
		svc.Logger = svc.Logger.WithPrefix("http")
		mw = append(mw, middleware.LoggerWithConfig(middleware.LoggerConfig{
			Output: logger.NewEchoMiddlewareLogger(svc.Logger),
			Format: `{"host": "${host}", "path": "${path}", "method": "${method}", "status_code": ${status}, "latency": "${latency_human}"}`,
		}))
	} else {
		svc.Logger = logger.NOP{}
	}

	svc.e.Use(mw...)

	restv1.RegisterHandlers(svc.e, &restv1.Handler{
		Logger:  svc.Logger,
		AppName: svc.Name,
	})

	svc.ui = &spa.SPA{
		Logger:    svc.Logger,
		MountPath: "/*",
	}
	svc.ui.MountToEchoRouter(svc.e)

	return svc.e.Start(fmt.Sprintf("%s:%d", svc.Host, svc.Port))
}
