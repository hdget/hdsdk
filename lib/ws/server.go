package ws

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hdget/hdsdk/v2/intf"
	"io"
	"net/http"
	"strings"
)

type WebServer interface {
	Start() error
	Stop() error
	SetMode(mode string)
	GracefulStop(ctx context.Context) error
	AddRoutes(routes []*Route) error
	NewRouterGroup(urlPrefix string) *RouterGroup
}

type baseServer struct {
	*http.Server
	engine *gin.Engine
	logger intf.LoggerProvider
	params *ServerParam
}

func newBaseServer(logger intf.LoggerProvider, address string, options ...ServerOption) *baseServer {
	engine := getDefaultGinEngine(logger)
	s := &baseServer{
		Server: &http.Server{
			Addr:    address,
			Handler: engine,
		},
		engine: engine,
		logger: logger,
		params: defaultServerParams,
	}

	for _, option := range options {
		option(s.params)
	}

	return s
}

func (b baseServer) Stop() error {
	if err := b.Server.Close(); err != nil {
		return err
	}
	return nil
}

func (b baseServer) GracefulStop(ctx context.Context) error {
	if err := b.Server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

func (b baseServer) AddRoutes(routes []*Route) error {
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
			b.engine.GET(route.Path, route.Handler)
		case "POST":
			b.engine.POST(route.Path, route.Handler)
		}
	}
	return nil
}

// SetMode set ws to specific mode
func (b baseServer) SetMode(mode string) {
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
