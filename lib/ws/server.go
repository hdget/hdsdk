package ws

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/hdget/hdsdk"
	"github.com/hdget/hdsdk/types"
	"github.com/hdget/hdsdk/utils/parallel"
	"net/http"
	"syscall"
	"time"
)

type HttpServer struct {
	*http.Server
	router *gin.Engine
}

type HttpMethod int

const (
	Get HttpMethod = iota
	Post
	Delete
)
const waitTime = 3 * time.Second

func NewHttpServer(logger types.LogProvider, address string) HttpServer {
	router := NewRouter(logger)
	return HttpServer{
		Server: &http.Server{
			Addr:    address,
			Handler: router,
		},
		router: router,
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

// AddRoute add route handler
func (srv *HttpServer) AddRoute(method HttpMethod, path string, handler gin.HandlerFunc) {
	switch method {
	case Get:
		srv.router.GET(path, handler)
	case Post:
		srv.router.POST(path, handler)
	case Delete:
		srv.router.DELETE(path, handler)
	}
}

func (srv *HttpServer) shutdown(err error) {
	ctx, cancel := context.WithTimeout(context.Background(), waitTime)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		hdsdk.Logger.Fatal("http server shutdown", "error", err)
	}
}
