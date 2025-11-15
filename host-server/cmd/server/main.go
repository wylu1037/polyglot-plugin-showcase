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
		// Configuration
		fx.Provide(func() (*config.Config, error) {
			return config.Load("")
		}),

		// Database
		fx.Provide(database.NewDatabase),
		fx.Invoke(database.AutoMigrate),

		// Echo Web Server
		fx.Provide(bootstrap.NewEchoApp),

		// Plugin System
		fx.Provide(
			plugin.ProvideRegistry,
			plugin.ProvideManager,
			plugin.ProvidePluginDir,
		),
		fx.Invoke(plugin.AutoLoadPlugins),

		// Modules
		plugins.Module,

		// Start Web Server
		fx.Invoke(bootstrap.Start),
	)

	app.Run()
}
