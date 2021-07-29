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
	Address     string   `mapstructure:"address"`
	Transports  []string `mapstructure:"transports"`
	Middlewares []string `mapstructure:"middlewares"`
}

type GokitServer struct {
	Name        string
	Logger      types.LogProvider
	Config      *GokitServerConfig
	grpcServer  *grpc.Server
	Middlewares []endpoint.Middleware
	Tracer      *Tracer
	Options     []kitgrpc.ServerOption

	ctx    context.Context
	cancel context.CancelFunc
}

type RegisterFunc func(grpcServer *grpc.Server, concreteService interface{})

var _ types.MsServer = (*GokitServer)(nil)

// CreateServer 创建微服务server
func (msi MicroServiceImpl) CreateServer() types.MsServer {
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
	mdws := make([]endpoint.Middleware, 0)
	for _, mdwName := range msi.Config.Server.Middlewares {
		f := NewMdwFunctions[mdwName]
		if f != nil {
			mdws = append(mdws, f(msi.Config))
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &GokitServer{
		Logger:      msi.Logger,
		Name:        msi.Name,
		grpcServer:  grpcServer,
		Options:     serverOptions,
		Config:      msi.Config.Server,
		Middlewares: mdws,
		Tracer:      msi.Tracer,
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (gs *GokitServer) GetGrpcServer() *grpc.Server {
	return gs.grpcServer
}

// Run 运行GrpcServer
func (gs *GokitServer) Run() {
	var group parallel.Group
	{
		listener, err := net.Listen("tcp", gs.Config.Address)
		if err != nil {
			gs.Logger.Fatal("new server listener", "err", err)
		}
		group.Add(func() error {
			return gs.grpcServer.Serve(listener)
		}, func(error) {
			listener.Close()
			gs.grpcServer.Stop()
		})
	}
	{
		// 添加信用监听
		group.Add(parallel.SignalActor(gs.ctx, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT))
	}

	gs.Logger.Debug("microservice is running", "name", gs.Name, "address", gs.Config.Address)
	err := group.Run()
	if err != nil {
		gs.Logger.Debug("microservice exit", "err", err)
	}
}

// Close 关闭GrpcServer
func (gs *GokitServer) Close() {
	gs.cancel()
}
