package bootstrap

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/wylu1037/polyglot-plugin-host-server/app/router"
	"github.com/wylu1037/polyglot-plugin-host-server/config"
	"go.uber.org/fx"
)

type WebServerParams struct {
	fx.In
	Router *router.Router
	Config *config.Config
}

func Start(lc fx.Lifecycle, p WebServerParams) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}

func NewEchoApp(config *config.Config) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	return e
}
