package gokit

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/tracing/opentracing"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/hdget/sdk/types"
	"github.com/hdget/sdk/utils/parallel"
	"net"
	"net/http"
	"syscall"
)

type GokitHttpServer struct {
	BaseGokitServer
	Options []kithttp.ServerOption
}

var _ types.MsHttpServer = (*GokitHttpServer)(nil)

// CreateServer 创建微服务server
func (msi MicroServiceImpl) CreateHttpServer() types.MsHttpServer {
	// set serverOptions
	serverOptions := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(msi.Logger)),
		// Zipkin GRPC Trace can either be instantiated per gRPC method with a
		// provided operation name or a global tracing service can be instantiated
		// without an operation name and fed to each Go kit gRPC server as a GrpcServiceServerConfig.
		// In the latter case, the operation name will be the endpoint's grpc method
		// path if used in combination with the Go kit gRPC Interceptor.
		kitzipkin.HTTPServerTrace(msi.Tracer.ZipkinTracer),
	}

	// 添加中间件
	mdws := make([]endpoint.Middleware, 0)
	serverConfig := msi.GetServerConfig("http")
	for _, mdwName := range serverConfig.Middlewares {
		f := NewMdwFunctions[mdwName]
		if f != nil {
			mdws = append(mdws, f(msi.Config))
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &GokitHttpServer{
		BaseGokitServer: BaseGokitServer{
			Config:      serverConfig,
			Logger:      msi.Logger,
			Name:        msi.Name,
			Middlewares: mdws,
			Tracer:      msi.Tracer,
			ctx:         ctx,
			cancel:      cancel,
		},
		Options: serverOptions,
	}
}

// Run 运行GrpcServer
func (ghs *GokitHttpServer) Run(handlers map[string]*kithttp.Server) error {
	var group parallel.Group
	{
		// The HTTP listener mounts the Go kit HTTP handler we created.
		httpListener, err := net.Listen("tcp", ghs.Config.Address)
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
		group.Add(parallel.SignalActor(ghs.ctx, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT))
	}

	ghs.Logger.Debug("microservice is running", "name", ghs.Name, "address", ghs.Config.Address)
	err := group.Run()
	if err != nil {
		ghs.Logger.Debug("microservice exit", "err", err)
	}
	return err
}

// CreateHandler 创建Http Transport的Handler
func (ghs *GokitHttpServer) CreateHandler(concreteService interface{}, h types.HttpEndpoint) *kithttp.Server {
	// 将具体的service和middleware串联起来
	endpoints := h.MakeEndpoint(concreteService)
	for _, m := range ghs.Middlewares {
		endpoints = m(endpoints)
	}

	// 添加tracer到ServerBefore
	options := append(ghs.Options,
		kithttp.ServerBefore(
			opentracing.HTTPToContext(ghs.Tracer.OpenTracer, h.GetName(), ghs.Logger),
		),
	)

	return kithttp.NewServer(
		endpoints,
		h.ServerDecodeRequest,
		h.ServerEncodeResponse,
		options...,
	)
}
