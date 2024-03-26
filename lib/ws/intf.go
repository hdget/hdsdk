package ws

import "github.com/gin-gonic/gin"

type WebServer interface {
	Run()
	AddPublicRouterGroup(middlewares []gin.HandlerFunc, routes []*Route) error
	AddProtectRouterGroup(middlewares []gin.HandlerFunc, routes []*Route) error
	AddRoutes(routes []*Route) error
	SetMode(mode string)
}

type Route struct {
	Method  HttpMethod
	Path    string
	Handler gin.HandlerFunc
}
