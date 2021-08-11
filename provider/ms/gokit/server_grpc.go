package gokit

import (
	"context"
	"github.com/go-kit/kit/transport"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/hdget/hdsdk/types"
	"github.com/hdget/hdsdk/utils/parallel"
	"google.golang.org/grpc"
	"net"
	"syscall"
)

type GokitGrpcServer struct {
	BaseGokitServer
	grpcServer *grpc.Server
	Options    []kitgrpc.ServerOption
}

var _ types.GrpcServerManager = (*GokitGrpcServer)(nil)

// CreateGrpcServer 创建微服务server
func (msi MicroServiceImpl) NewGrpcServerManager() types.GrpcServerManager {
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))

	// global serverOptions
	serverOptions := []kitgrpc.ServerOption{
		kitgrpc.ServerErrorHandler(transport.NewLogErrorHandler(msi.Logger)),
	}

	// 添加中间件
	mdws := make([]*MsMiddleware, 0)
	serverConfig := msi.GetServerConfig(GRPC)
	for _, mdwName := range serverConfig.Middlewares {
		newFunc := NewMdwFunctions[mdwName]
		if newFunc != nil {
			mdws = append(mdws, newFunc(msi.Config))
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &GokitGrpcServer{
		BaseGokitServer: BaseGokitServer{
			Config:      serverConfig,
			Logger:      msi.Logger,
			Name:        msi.Name,
			Middlewares: mdws,
			ctx:         ctx,
			cancel:      cancel,
		},
		grpcServer: grpcServer,
		Options:    serverOptions,
	}

}

func (s *GokitGrpcServer) GetServer() *grpc.Server {
	return s.grpcServer
}

// Run 运行GrpcServer
func (s *GokitGrpcServer) RunServer() error {
	var group parallel.Group
	{
		listener, err := net.Listen("tcp", s.Config.Address)
		if err != nil {
			return err
		}
		group.Add(func() error {
			return s.grpcServer.Serve(listener)
		}, func(error) {
			listener.Close()
			s.grpcServer.Stop()
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

// CreateHandler 创建Grpc Transport的handler
func (s *GokitGrpcServer) CreateHandler(concreteService interface{}, ep types.GrpcAspect) *kitgrpc.Server {
	// 将具体的service和middleware串联起来
	endpoints := ep.MakeEndpoint(concreteService)
	for _, m := range s.Middlewares {
		if m.Middleware != nil {
			endpoints = m.Middleware(endpoints)
		}

		if len(m.InjectFunctions) > 0 {
			injectFunc := m.InjectFunctions[GRPC]
			if injectFunc != nil {
				_, serverOptions := injectFunc(s.Logger, ep.GetServiceName())
				for _, option := range serverOptions {
					svrOption, ok := option.(kitgrpc.ServerOption)
					if ok {
						s.Options = append(s.Options, svrOption)
					}
				}
			}
		}
	}

	return kitgrpc.NewServer(
		endpoints,
		ep.ServerDecodeRequest,
		ep.ServerEncodeResponse,
		s.Options...,
	)
}
