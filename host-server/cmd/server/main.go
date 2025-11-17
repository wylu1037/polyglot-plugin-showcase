package main

import (
	"github.com/wylu1037/polyglot-plugin-host-server/app/database"
	"github.com/wylu1037/polyglot-plugin-host-server/app/modules/plugins"
	"github.com/wylu1037/polyglot-plugin-host-server/config"
	"github.com/wylu1037/polyglot-plugin-host-server/internal/bootstrap"
	"github.com/wylu1037/polyglot-plugin-host-server/internal/plugin"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(config.Load),
		fx.Provide(database.NewDatabase),
		fx.Provide(
			plugin.ProvideRegistry,
			plugin.ProvideManager,
			plugin.ProvidePluginDir,
		),
		fx.Provide(bootstrap.NewEchoApp),
		plugins.Module,
		fx.Invoke(database.AutoMigrate),
		fx.Invoke(plugin.AutoLoadPlugins),
		fx.Invoke(bootstrap.Start),
	)

	app.Run()
}
