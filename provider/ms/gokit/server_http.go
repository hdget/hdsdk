package gokit

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/hdget/hdsdk/types"
	"github.com/hdget/hdsdk/utils/parallel"
	"net"
	"net/http"
	"syscall"
)

type GokitHttpServer struct {
	BaseServerManager
	Options []kithttp.ServerOption
}

type errorWrapper struct {
	Error string `json:"error"`
}

var _ types.HttpServerManager = (*GokitHttpServer)(nil)

// NewHttpServerManager 创建微服务server
func (msi MicroServiceImpl) NewHttpServerManager() types.HttpServerManager {
	// set serverOptions
	serverOptions := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(HttpErrorEncoder),
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(msi.Logger)),
	}

	// 添加中间件
	mdws := make([]*MsMiddleware, 0)
	serverConfig := msi.GetServerConfig(HTTP)
	if serverConfig != nil {
		for _, mdwName := range serverConfig.Middlewares {
			newFunc := NewMdwFunctions[mdwName]
			if newFunc != nil {
				mdws = append(mdws, newFunc(msi.Config))
			}
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &GokitHttpServer{
		BaseServerManager: BaseServerManager{
			Config:      serverConfig,
			Logger:      msi.Logger,
			Name:        msi.Name,
			Middlewares: mdws,
			ctx:         ctx,
			cancel:      cancel,
		},
		Options: serverOptions,
	}
}

// RunServer 运行GrpcServer
func (s *GokitHttpServer) RunServer(handlers map[string]http.Handler) error {
	var group parallel.Group
	{
		// The HTTP listener mounts the Go kit HTTP handler we created.
		httpListener, err := net.Listen("tcp", s.Config.Address)
		if err != nil {
			return err
		}

		// new http handler
		m := http.NewServeMux()
		for url, handler := range handlers {
			m.Handle(url, handler)
		}

		group.Add(func() error {
			return http.Serve(httpListener, m)
		}, func(error) {
			httpListener.Close()
		})
	}
	{
		// 添加信用监听
		group.Add(parallel.SignalActor(s.ctx, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT))
	}

	s.Logger.Debug("microservice is running", "name", s.Name, "address", s.Config.Address)
	err := group.Run()
	if err != nil {
		s.Logger.Debug("microservice exit", "err", err)
	}
	return err
}

// CreateHandler 创建Http Transport的Handler
func (s *GokitHttpServer) CreateHandler(concreteService interface{}, ap types.HttpAspect) http.Handler {
	// 将具体的service和middleware串联起来
	endpoints := ap.MakeEndpoint(concreteService)
	for _, m := range s.Middlewares {
		if m.Middleware != nil {
			endpoints = m.Middleware(endpoints)
		}

		if len(m.InjectFunctions) == 0 {
			continue
		}

		if injectFunc := m.InjectFunctions[HTTP]; injectFunc != nil {
			_, serverOptions := injectFunc(s.Logger, ap.GetMethodName())
			for _, option := range serverOptions {
				if svrOption, ok := option.(kithttp.ServerOption); ok {
					s.Options = append(s.Options, svrOption)
				}
			}
		}
	}

	return kithttp.NewServer(
		endpoints,
		ap.ServerDecodeRequest,
		ap.ServerEncodeResponse,
		s.Options...,
	)
}

//nolint:errcheck
func HttpErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}
