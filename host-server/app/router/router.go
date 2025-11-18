package router

import (
	"github.com/wylu1037/polyglot-plugin-host-server/app/modules/plugins"
)

type Router struct {
	plugins *plugins.Route
}

func NewRouter(
	plugins *plugins.Route,
) *Router {
	return &Router{
		plugins: plugins,
	}
}

func (r *Router) Register() {
	r.plugins.Register()
}
