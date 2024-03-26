package ws

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hdget/hdsdk/types"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"time"
)

type baseServer struct {
	*http.Server
	engine       *gin.Engine
	logger       types.LogProvider
	params       *ServerParam
	routerGroups map[string]*gin.RouterGroup
}

const (
	gracefulShutdownWaitTime = 3 * time.Second
)

func newBaseServer(logger types.LogProvider, address string, options ...ServerOption) *baseServer {
	s := &baseServer{
		Server: &http.Server{
			Addr:    address,
			Handler: getDefaultGinEngine(logger),
		},
		engine: getDefaultGinEngine(logger),
		logger: logger,
		params: defaultServerParams,
	}

	for _, option := range options {
		option(s.params)
	}

	return s
}

func (w *baseServer) AddPublicRouterGroup(middlewares []gin.HandlerFunc, routes []*Route) error {
	return w.addRouterGroup(w.params.publicRouterGroup.Name, w.params.publicRouterGroup.UrlPrefix, middlewares, routes)
}

func (w *baseServer) AddProtectRouterGroup(middlewares []gin.HandlerFunc, routes []*Route) error {
	return w.addRouterGroup(w.params.protectRouterGroup.Name, w.params.protectRouterGroup.UrlPrefix, middlewares, routes)
}

func (w *baseServer) addRouterGroup(name, urlPrefix string, middlewares []gin.HandlerFunc, routes []*Route) error {
	if _, exists := w.routerGroups[name]; exists {
		return errors.Wrapf(ErrDuplicateRouterGroup, "name: %s", name)
	}

	// new router group
	w.routerGroups[name] = w.engine.Group(urlPrefix, middlewares...)

	routeMap := make(map[string]struct{})
	for _, route := range routes {
		// 先检查是否有重复的路由
		k := fmt.Sprintf("%d_%s", route.Method, route.Path)
		if _, exist := routeMap[k]; exist {
			return fmt.Errorf("duplicate route, url: %s, method: %d", route.Path, route.Method)
		}

		// 添加到router group
		switch route.Method {
		case Get:
			w.routerGroups[name].GET(route.Path, route.Handler)
		case Post:
			w.routerGroups[name].POST(route.Path, route.Handler)
		case Delete:
			w.routerGroups[name].DELETE(route.Path, route.Handler)
		}
	}
	return nil
}

func (w *baseServer) AddRoutes(routes []*Route) error {
	routeMap := make(map[string]struct{})
	for _, route := range routes {
		// 先检查是否有重复的路由
		k := fmt.Sprintf("%d_%s", route.Method, route.Path)
		if _, exist := routeMap[k]; exist {
			return fmt.Errorf("duplicate route, url: %s, method: %d", route.Path, route.Method)
		}

		// 添加到router group
		switch route.Method {
		case Get:
			w.engine.GET(route.Path, route.Handler)
		case Post:
			w.engine.POST(route.Path, route.Handler)
		case Delete:
			w.engine.DELETE(route.Path, route.Handler)
		}
	}
	return nil
}

// SetMode set gin to specific mode
func (w *baseServer) SetMode(mode string) {
	gin.SetMode(mode)
}

func (w *baseServer) shutdown(err error) {
	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownWaitTime)
	defer cancel()
	if e := w.Shutdown(ctx); e != nil {
		w.logger.Fatal("http server shutdown", "reason", err, "error", e)
	}
}

func getDefaultGinEngine(logger types.LogProvider) *gin.Engine {
	// new router
	engine := gin.New()

	// set gin to logout to stdout and file
	gin.DefaultWriter = io.MultiWriter(logger.GetStdLogger().Writer())

	// add basic middleware
	engine.Use(
		gin.Recovery(),
		newLoggerMiddleware(logger), // logger middleware
	)

	return engine
}
