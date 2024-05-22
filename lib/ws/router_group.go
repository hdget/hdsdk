package ws

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

type RouterGroup struct {
	ginRouterGroup *gin.RouterGroup
	UrlPrefix      string
}

func (b baseServer) NewRouterGroup(urlPrefix string) *RouterGroup {
	return &RouterGroup{
		ginRouterGroup: b.engine.Group(urlPrefix),
		UrlPrefix:      urlPrefix,
	}
}

func (rg *RouterGroup) Use(middlewares ...gin.HandlerFunc) *RouterGroup {
	rg.ginRouterGroup.Use(middlewares...)
	return rg
}

func (rg *RouterGroup) AddRoute(routes ...*Route) error {
	routeMap := make(map[string]struct{})
	for _, route := range routes {
		// 先检查是否有重复的路由
		k := fmt.Sprintf("%s_%s", route.Method, route.Path)
		if _, exist := routeMap[k]; exist {
			return fmt.Errorf("duplicate route, url: %s, method: %s", route.Path, route.Method)
		}

		// 添加到router group
		switch strings.ToUpper(route.Method) {
		case "GET":
			rg.ginRouterGroup.GET(route.Path, route.Handler)
		case "POST":
			rg.ginRouterGroup.POST(route.Path, route.Handler)
		}
	}
	return nil
}
