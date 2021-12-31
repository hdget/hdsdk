package ws

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/hdget/hdsdk"
	"github.com/hdget/hdsdk/lib/err"
	"github.com/hdget/hdsdk/types"
	"github.com/hdget/hdsdk/utils/parallel"
	"net/http"
	"syscall"
	"time"
)

type HttpServer struct {
	*http.Server
	router       *gin.Engine
	routerGroups map[string]*gin.RouterGroup
}

type HttpMethod int

const (
	Get HttpMethod = iota
	Post
	Delete
)
const waitTime = 3 * time.Second

var (
	ErrDuplicateRouterGroup = err.New("duplicate router group")
	ErrRouterGroupNotFound  = err.New("router group not found")
)

func NewHttpServer(logger types.LogProvider, address string) *HttpServer {
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

// SetReleaseMode set gin to release mode
func SetReleaseMode() {
	gin.SetMode(gin.ReleaseMode)
}

// SetDebugMode set gin to debug mode
func SetDebugMode() {
	gin.SetMode(gin.DebugMode)
}

// SetTestMode set gin to test mode
func SetTestMode() {
	gin.SetMode(gin.TestMode)
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

	if err := group.Run(); err != nil {
		hdsdk.Logger.Error("http server quit", "error", err)
	}
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

func addToRouterGroup(routerGroup *gin.RouterGroup, method HttpMethod, path string, handler gin.HandlerFunc) {
	switch method {
	case Get:
		routerGroup.GET(path, handler)
	case Post:
		routerGroup.POST(path, handler)
	case Delete:
		routerGroup.DELETE(path, handler)
	}
}

func addToRouter(router *gin.Engine, method HttpMethod, path string, handler gin.HandlerFunc) {
	switch method {
	case Get:
		router.GET(path, handler)
	case Post:
		router.POST(path, handler)
	case Delete:
		router.DELETE(path, handler)
	}
}
