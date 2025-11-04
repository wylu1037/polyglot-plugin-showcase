package pluginstore

import (
	"github.com/wylu1037/polyglot-plugin-host-server/app/modules/plugin_store/controller"
	"github.com/wylu1037/polyglot-plugin-host-server/app/modules/plugin_store/repository"
	"github.com/wylu1037/polyglot-plugin-host-server/app/modules/plugin_store/service"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewRoute),
	fx.Provide(repository.NewPluginStoreRepository),
	fx.Provide(service.NewPluginStoreService),
	fx.Provide(controller.NewPluginStoreController),
)
