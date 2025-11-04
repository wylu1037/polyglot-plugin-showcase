package router

import "github.com/labstack/echo/v4"

type Router struct {
	engine *echo.Echo
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Register() {}
