package ws

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strings"
)

type WebServer interface {
	Start() error
	Stop() error
	GracefulStop(ctx context.Context) error
	AddPublicRouterGroup(middlewares []gin.HandlerFunc, routes []*Route) error
	AddProtectRouterGroup(middlewares []gin.HandlerFunc, routes []*Route) error
	AddRoutes(routes []*Route) error
	SetMode(mode string)
}

type baseServer struct {
	*http.Server
	engine       *gin.Engine
	logger       intf.LoggerProvider
	params       *ServerParam
	routerGroups map[string]*gin.RouterGroup
}

func newBaseServer(logger intf.LoggerProvider, address string, options ...ServerOption) *baseServer {
	s := &baseServer{
		Server: &http.Server{
			Addr:    address,
			Handler: getDefaultGinEngine(logger),
		},
		engine:       getDefaultGinEngine(logger),
		logger:       logger,
		params:       defaultServerParams,
		routerGroups: make(map[string]*gin.RouterGroup),
	}

	for _, option := range options {
		option(s.params)
	}

	return s
}

func (w baseServer) AddPublicRouterGroup(middlewares []gin.HandlerFunc, routes []*Route) error {
	return w.addRouterGroup(w.params.publicRouterGroup.Name, w.params.publicRouterGroup.UrlPrefix, middlewares, routes)
}

func (w baseServer) AddProtectRouterGroup(middlewares []gin.HandlerFunc, routes []*Route) error {
	return w.addRouterGroup(w.params.protectRouterGroup.Name, w.params.protectRouterGroup.UrlPrefix, middlewares, routes)
}

func (w baseServer) Stop() error {
	if err := w.Server.Close(); err != nil {
		return err
	}
	return nil
}

func (w baseServer) GracefulStop(ctx context.Context) error {
	if err := w.Server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

func (w baseServer) AddRoutes(routes []*Route) error {
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
			w.engine.GET(route.Path, route.Handler)
		case "POST":
			w.engine.POST(route.Path, route.Handler)
		}
	}
	return nil
}

// SetMode set ws to specific mode
func (w baseServer) SetMode(mode string) {
	gin.SetMode(mode)
}

func getDefaultGinEngine(logger intf.LoggerProvider) *gin.Engine {
	// new router
	engine := gin.New()

	// set ws to logout to stdout and file
	gin.DefaultWriter = io.MultiWriter(logger.GetStdLogger().Writer())

	// add basic middleware
	engine.Use(
		gin.Recovery(),
		newLoggerMiddleware(logger), // logger middleware
	)

	return engine
}

func (w baseServer) addRouterGroup(name, urlPrefix string, middlewares []gin.HandlerFunc, routes []*Route) error {
	if _, exists := w.routerGroups[name]; exists {
		return errors.Wrapf(ErrDuplicateRouterGroup, "name: %s", name)
	}

	// new router group
	w.routerGroups[name] = w.engine.Group(urlPrefix, middlewares...)

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
			w.routerGroups[name].GET(route.Path, route.Handler)
		case "POST":
			w.routerGroups[name].POST(route.Path, route.Handler)
		}
	}
	return nil
}
