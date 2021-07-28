package gokit

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/kit/transport"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/hdget/sdk/types"
	"github.com/hdget/sdk/utils/parallel"
	"google.golang.org/grpc"
	"net"
	"syscall"
)

// GokitServerConfig 服务端配置
type GokitServerConfig struct {
	Url string `mapstructure:"url"`
	Transports  []string `mapstructure:"transports"`
	Middlewares []string `mapstructure:"middlewares"`
}

type GokitServer struct {
	Name               string
	Logger             types.LogProvider
	Config             *GokitServerConfig
	GrpcServer         *grpc.Server
	Middlewares 		[]endpoint.Middleware
	Tracer             *Tracer
	Options  	  	   []kitgrpc.ServerOption

	ctx                context.Context
	cancel             context.CancelFunc
}



var _ types.MsServer = (*GokitServer)(nil)

// CreateGrpcServer 将具体的服务实现注册到RPC服务中去，然后返回MicroServiceServer
// registerFunc是从proto文件生成的RegisterXxxServer函数
// concreteService是具体服务的实现结构
func (msi MicroServiceImpl) CreateGrpcServer(registerFunc types.RegisterFunc, concreteService interface{}) types.MsServer {
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))

	registerFunc(grpcServer, concreteService)

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
	mdws := make([]endpoint.Middleware, 0)
	for _, mdwName := range msi.Config.Server.Middlewares {
		f := NewMdwFunctions[mdwName]
		if f != nil {
			mdws = append(mdws, f(msi.Config))
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &GokitServer{
		GrpcServer: grpcServer,
		Options: serverOptions,
		Middlewares: mdws,
		Tracer: msi.Tracer,
		ctx: ctx,
		cancel: cancel,
	}
}

// Run 运行GrpcServer
func (gs *GokitServer) Run() {
	var group parallel.Group
	{
		listener, err := net.Listen("tcp", gs.Config.Url)
		if err != nil {
			gs.Logger.Fatal("new server listener" , "err", err)
		}
		group.Add(func() error {
			return gs.GrpcServer.Serve(listener)
		}, func(error) {
			listener.Close()
			gs.GrpcServer.Stop()
		})
	}
	{
		// 添加信用监听
		group.Add(parallel.SignalActor(gs.ctx, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT))
	}

	err := group.Run()
	if err != nil {
		gs.Logger.Debug("microservice exit", "err", err)
	}
	gs.Logger.Debug("microservice is running", "name", gs.Name)
}

// Close 关闭GrpcServer
func (gs *GokitServer) Close() {
	gs.cancel()
}
