package pluginstore

import (
	"github.com/labstack/echo/v4"
	"github.com/wylu1037/polyglot-plugin-host-server/app/modules/plugin_store/controller"
)

type Route struct {
	app        *echo.Echo
	controller controller.PluginStoreController
}

func NewRoute(
	app *echo.Echo,
	controller controller.PluginStoreController,
) *Route {
	return &Route{
		app:        app,
		controller: controller,
	}
}

func (r *Route) Register() {

}
