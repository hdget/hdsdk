package gokit

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/tracing/opentracing"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
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

var _ types.MsGrpcServer = (*GokitGrpcServer)(nil)

// CreateGrpcServer 创建微服务server
func (msi MicroServiceImpl) CreateGrpcServer() types.MsGrpcServer {
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))

	// set serverOptions
	serverOptions := []kitgrpc.ServerOption{
		kitgrpc.ServerErrorHandler(transport.NewLogErrorHandler(msi.Logger)),
		// Zipkin GRPC Trace can either be instantiated per gRPC method with a
		// provided operation name or a global tracing service can be instantiated
		// without an operation name and fed to each Go kit gRPC server as a GrpcServiceServerConfig.
		// In the latter case, the operation name will be the endpoint's grpc method
		// path if used in combination with the Go kit gRPC Interceptor.
		kitzipkin.GRPCServerTrace(msi.Tracer.ZipkinTracer),
	}

	// 添加中间件
	// 添加中间件
	mdws := make([]endpoint.Middleware, 0)
	serverConfig := msi.GetServerConfig(GRPC_SERVER)
	for _, mdwName := range serverConfig.Middlewares {
		f := NewMdwFunctions[mdwName]
		if f != nil {
			mdws = append(mdws, f(msi.Config))
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &GokitGrpcServer{
		BaseGokitServer: BaseGokitServer{
			Config:      serverConfig,
			Logger:      msi.Logger,
			Name:        msi.Name,
			Middlewares: mdws,
			Tracer:      msi.Tracer,
			ctx:         ctx,
			cancel:      cancel,
		},
		grpcServer: grpcServer,
		Options:    serverOptions,
	}

}

func (ggs *GokitGrpcServer) GetServer() *grpc.Server {
	return ggs.grpcServer
}

// Run 运行GrpcServer
func (ggs *GokitGrpcServer) Run() error {
	var group parallel.Group
	{
		listener, err := net.Listen("tcp", ggs.Config.Address)
		if err != nil {
			return err
		}
		group.Add(func() error {
			return ggs.grpcServer.Serve(listener)
		}, func(error) {
			listener.Close()
			ggs.grpcServer.Stop()
		})
	}
	{
		// 添加信用监听
		group.Add(parallel.SignalActor(ggs.ctx, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT))
	}

	ggs.Logger.Debug("microservice is running", "name", ggs.Name, "address", ggs.Config.Address)
	err := group.Run()
	if err != nil {
		ggs.Logger.Debug("microservice exit", "err", err)
	}
	return err
}

// CreateHandler 创建Grpc Transport的handler
func (ggs *GokitGrpcServer) CreateHandler(concreteService interface{}, ge types.GrpcEndpoint) *kitgrpc.Server {
	// 将具体的service和middleware串联起来
	endpoints := ge.MakeEndpoint(concreteService)
	for _, m := range ggs.Middlewares {
		endpoints = m(endpoints)
	}

	// 添加tracer到ServerBefore
	options := append(ggs.Options,
		kitgrpc.ServerBefore(
			opentracing.GRPCToContext(ggs.Tracer.OpenTracer, ge.GetName(), ggs.Logger),
		),
	)

	return kitgrpc.NewServer(
		endpoints,
		ge.ServerDecodeRequest,
		ge.ServerEncodeResponse,
		options...,
	)
}
