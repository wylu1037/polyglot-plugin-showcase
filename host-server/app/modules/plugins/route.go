package plugins

import (
	"github.com/labstack/echo/v4"
	"github.com/wylu1037/polyglot-plugin-host-server/app/modules/plugins/controller"
)

type Route struct {
	app        *echo.Echo
	controller controller.PluginController
}

func NewRoute(
	app *echo.Echo,
	controller controller.PluginController,
) *Route {
	return &Route{
		app:        app,
		controller: controller,
	}
}

func (r *Route) Register() {
	api := r.app.Group("/api/plugins")

	api.POST("/install", r.controller.InstallPlugin)
	api.GET("", r.controller.ListPlugins)
	api.GET("/:id", r.controller.GetPlugin)
	api.POST("/:id/activate", r.controller.ActivatePlugin)
	api.POST("/:id/deactivate", r.controller.DeactivatePlugin)
	api.DELETE("/:id", r.controller.UninstallPlugin)
	api.POST("/:id/call", r.controller.CallPlugin)
}
