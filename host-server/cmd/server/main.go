package main

import (
	"github.com/wylu1037/polyglot-plugin-host-server/config"
	"github.com/wylu1037/polyglot-plugin-host-server/internal/bootstrap"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(config.Load),
		fx.Provide(bootstrap.NewEchoApp),
	)

	app.Run()
}
