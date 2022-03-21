package ws

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/hdget/hdsdk"
	"github.com/hdget/hdsdk/types"
	"github.com/hdget/hdsdk/utils/parallel"
	"github.com/pkg/errors"
	"net/http"
	"syscall"
)

type WebServer interface {
	Run()
	NewRouterGroup(name, relativePath string, handlers ...gin.HandlerFunc) (*gin.RouterGroup, error)
	CreateRouterGroup(name, relativePath string, handlers ...gin.HandlerFunc) error
	GetRouterGroup(groupName string) *gin.RouterGroup
	AddRoutes(routes []*Route)
	AddGroupRoutes(groupName string, routes []*Route) error
}

type HttpServer struct {
	*http.Server
	router       *gin.Engine
	routerGroups map[string]*gin.RouterGroup
}

func NewHttpServer(logger types.LogProvider, address string) WebServer {
	router := NewRouter(logger)
	return &HttpServer{
		Server: &http.Server{
			Addr:    address,
			Handler: router,
		},
		router:       router,
		routerGroups: make(map[string]*gin.RouterGroup),
	}
}

// Run http server
func (srv *HttpServer) Run() {
	var group parallel.Group
	{
		group.Add(srv.ListenAndServe, srv.shutdown)
	}
	{
		group.Add(
			parallel.SignalActor(
				context.Background(),
				syscall.SIGINT,
				syscall.SIGQUIT,
				syscall.SIGTERM,
				syscall.SIGKILL,
			),
		)
	}

	if err := group.Run(); err != nil && err != http.ErrServerClosed {
		hdsdk.Logger.Error("http server quit", "error", err)
	}
}

// NewRouterGroup new a gin router group
func (srv *HttpServer) NewRouterGroup(name, relativePath string, handlers ...gin.HandlerFunc) (*gin.RouterGroup, error) {
	err := srv.CreateRouterGroup(name, relativePath)
	if err != nil {
		return nil, errors.Wrapf(err, "create router group, name: %s, prefix: %s", name, relativePath)

	}
	return srv.GetRouterGroup(name), nil
}

// CreateRouterGroup create a gin router group
func (srv *HttpServer) CreateRouterGroup(name, relativePath string, handlers ...gin.HandlerFunc) error {
	_, exist := srv.routerGroups[name]
	if exist {
		return ErrDuplicateRouterGroup
	}

	srv.routerGroups[name] = srv.router.Group(relativePath, handlers...)
	return nil
}

func (srv *HttpServer) GetRouterGroup(groupName string) *gin.RouterGroup {
	return srv.routerGroups[groupName]
}

// AddRoutes add routes from Route slice
func (srv *HttpServer) AddRoutes(routes []*Route) {
	for _, r := range routes {
		addToRouter(srv.router, r.Method, r.Path, r.Handler)
	}
}

func (srv *HttpServer) AddGroupRoutes(groupName string, routes []*Route) error {
	routerGroup := srv.GetRouterGroup(groupName)
	if routerGroup == nil {
		return ErrRouterGroupNotFound
	}

	for _, r := range routes {
		addToRouterGroup(routerGroup, r.Method, r.Path, r.Handler)
	}
	return nil
}

func (srv *HttpServer) shutdown(err error) {
	ctx, cancel := context.WithTimeout(context.Background(), waitTime)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		hdsdk.Logger.Fatal("http server shutdown", "error", err)
	}
}
