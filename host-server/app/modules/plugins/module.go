package plugins

import (
	"github.com/wylu1037/polyglot-plugin-host-server/app/modules/plugins/controller"
	"github.com/wylu1037/polyglot-plugin-host-server/app/modules/plugins/repository"
	"github.com/wylu1037/polyglot-plugin-host-server/app/modules/plugins/service"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewRoute),
	fx.Provide(repository.NewPluginRepository),
	fx.Provide(service.NewPluginService),
	fx.Provide(controller.NewPluginController),
)
