package bootstrap

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/wylu1037/polyglot-plugin-host-server/app/common"
	"github.com/wylu1037/polyglot-plugin-host-server/config"
	"go.uber.org/fx"
)

type WebServerParams struct {
	fx.In
	Echo   *echo.Echo
	Config *config.Config
}

func Start(lc fx.Lifecycle, p WebServerParams) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				addr := p.Config.GetServerAddr()
				fmt.Printf("Starting server on %s\n", addr)
				if err := p.Echo.Start(addr); err != nil {
					fmt.Printf("Server error: %v\n", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Shutting down server...")
			return p.Echo.Shutdown(ctx)
		},
	})
}

func NewEchoApp(config *config.Config) *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	// Set custom error handler
	e.HTTPErrorHandler = common.CustomHTTPErrorHandler

	e.Use(middleware.Logger(), middleware.Recover(), middleware.CORS())

	return e
}
