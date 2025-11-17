package main

import (
	"github.com/wylu1037/polyglot-plugin-host-server/app/database"
	"github.com/wylu1037/polyglot-plugin-host-server/app/modules/plugins"
	"github.com/wylu1037/polyglot-plugin-host-server/config"
	"github.com/wylu1037/polyglot-plugin-host-server/internal/bootstrap"
	"github.com/wylu1037/polyglot-plugin-host-server/internal/plugin"
	"go.uber.org/fx"

	_ "github.com/wylu1037/polyglot-plugin-host-server/docs"
)

// @title           Polyglot Plugin Host Server API
// @version         1.0
// @description     A plugin management system that supports dynamic loading and execution of plugins
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:8080
// @BasePath  /
// @schemes http https
// @tag.name plugins
// @tag.description Plugin management operations
func main() {
	app := fx.New(
		fx.Supply(""),
		fx.Provide(config.Load),
		fx.Provide(database.NewDatabase),
		fx.Provide(
			plugin.ProvideRegistry,
			plugin.ProvideManager,
		),
		fx.Provide(bootstrap.NewEchoApp),
		plugins.Module,
		fx.Invoke(database.AutoMigrate),
		fx.Invoke(plugin.AutoLoadPlugins),
		fx.Invoke(bootstrap.Start),
	)

	app.Run()
}
