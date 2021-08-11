package gokit

import (
	"context"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/hdget/hdsdk/types"
	"github.com/hdget/hdsdk/utils/parallel"
	"net"
	"net/http"
	"syscall"
)

type GokitHttpServer struct {
	BaseGokitServer
	Options []kithttp.ServerOption
}

var _ types.MsHttpServer = (*GokitHttpServer)(nil)

// CreateHttpServer 创建微服务server
func (msi MicroServiceImpl) CreateHttpServer() types.MsHttpServer {
	// set serverOptions
	serverOptions := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(msi.Logger)),
	}

	// 添加中间件
	mdws := make([]*MsMiddleware, 0)
	serverConfig := msi.GetServerConfig(HTTP_SERVER)
	for _, mdwName := range serverConfig.Middlewares {
		newFunc := NewMdwFunctions[mdwName]
		if newFunc != nil {
			mdws = append(mdws, newFunc(msi.Config))
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &GokitHttpServer{
		BaseGokitServer: BaseGokitServer{
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

// Run 运行GrpcServer
func (s *GokitHttpServer) Run(handlers map[string]*kithttp.Server) error {
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
func (s *GokitHttpServer) CreateHandler(concreteService interface{}, ep types.HttpEndpoint) *kithttp.Server {
	// 将具体的service和middleware串联起来
	endpoints := ep.MakeEndpoint(concreteService)
	for _, m := range s.Middlewares {
		if m.Middleware != nil {
			endpoints = m.Middleware(endpoints)
		}

		if len(m.InjectFunctions) > 0 {
			injectFunc := m.InjectFunctions[HTTP_SERVER]
			if injectFunc != nil {
				_, serverOptions := injectFunc(s.Logger, ep.GetName())
				for _, option := range serverOptions {
					svrOption, ok := option.(kithttp.ServerOption)
					if ok {
						s.Options = append(s.Options, svrOption)
					}
				}
			}
		}
	}

	return kithttp.NewServer(
		endpoints,
		ep.ServerDecodeRequest,
		ep.ServerEncodeResponse,
		s.Options...,
	)
}
