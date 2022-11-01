package gokit

import (
	kitot "github.com/go-kit/kit/tracing/opentracing"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/reporter"
	"hdsdk/types"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/pkg/errors"
)

// TraceConfig 服务跟踪选项
type TraceConfig struct {
	Url     string `mapstructure:"url"` // Zipkin tracing HTTP reporter URL
	Address string `mapstructure:"address"`
}

type Tracer struct {
	Reporter     reporter.Reporter
	ZipkinTracer *zipkin.Tracer     // zipkin tracer
	OpenTracer   opentracing.Tracer // opentracing tracer
}

var (
	defaultTraceConfig = &TraceConfig{
		Url:     "http://localhost:9411/api/v2/spans",
		Address: "localhost:80",
	}
	ErrInvalidTraceConfig = errors.New("invalid trace config")
)

// NewMdwTrace returns an endpoint.Middleware that acts as a tracer.
// Requests that would exceed the
// maximum request rate are simply rejected with an error.
func NewMdwTrace(config *MicroServiceConfig) *MsMiddleware {
	tracer, err := newTracer(config)
	if err != nil {
		return nil
	}

	return &MsMiddleware{
		InjectFunctions: map[string]InjectFunction{
			GRPC: tracer.getGrpcOptions,
			HTTP: tracer.getHttpOptions,
		},
	}
}

func newTracer(config *MicroServiceConfig) (*Tracer, error) {
	// 如果没有配置tracer不做任何事
	traceConfig := config.getTraceConfig()

	// set up a span reporter
	reporter := zipkinhttp.NewReporter(traceConfig.Url)

	// create our local service endpoint
	endpoint, err := zipkin.NewEndpoint(config.Name, traceConfig.Address)
	if err != nil {
		return nil, err
	}

	// initialize our openTracer
	zipkinTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		return nil, err
	}

	// use zipkin-go-opentracing to wrap our openTracer
	openTracer := zipkinot.Wrap(zipkinTracer)

	// optionally set as Global OpenTracing openTracer instance
	opentracing.SetGlobalTracer(openTracer)

	return &Tracer{
		Reporter:     reporter,
		ZipkinTracer: zipkinTracer,
		OpenTracer:   openTracer,
	}, nil
}

func (m MicroServiceConfig) getTraceConfig() *TraceConfig {
	// 如果没有配置tracer不做任何事
	traceConfig := m.Trace
	if traceConfig == nil {
		traceConfig = defaultTraceConfig
	}

	if traceConfig.Url == "" {
		traceConfig.Url = defaultTraceConfig.Url
	}

	if traceConfig.Address == "" {
		traceConfig.Address = defaultTraceConfig.Address
	}

	return traceConfig
}

func (t *Tracer) getGrpcOptions(logger types.LogProvider, endpointName string) ([]interface{}, []interface{}) {
	clientOptions := []interface{}{
		kitzipkin.GRPCClientTrace(t.ZipkinTracer),
		kitgrpc.ClientBefore(
			kitot.ContextToGRPC(t.OpenTracer, logger),
		),
	}

	serverOptions := []interface{}{
		// Zipkin GRPC Trace can either be instantiated per gRPC method with a
		// provided operation name or a global tracing service can be instantiated
		// without an operation name and fed to each Go kit gRPC server as a GrpcServiceServerConfig.
		// In the latter case, the operation name will be the endpoint's grpc method
		// path if used in combination with the Go kit gRPC Interceptor.
		kitzipkin.GRPCServerTrace(t.ZipkinTracer),
		kitgrpc.ServerBefore(
			kitot.GRPCToContext(t.OpenTracer, endpointName, logger),
		),
	}
	return clientOptions, serverOptions
}

func (t *Tracer) getHttpOptions(logger types.LogProvider, endpointName string) ([]interface{}, []interface{}) {
	clientOptions := []interface{}{
		kitzipkin.HTTPClientTrace(t.ZipkinTracer),
		kithttp.ClientBefore(
			kitot.HTTPToContext(t.OpenTracer, endpointName, logger),
		),
	}

	serverOptions := []interface{}{
		// Zipkin GRPC Trace can either be instantiated per gRPC method with a
		// provided operation name or a global tracing service can be instantiated
		// without an operation name and fed to each Go kit gRPC server as a GrpcServiceServerConfig.
		// In the latter case, the operation name will be the endpoint's grpc method
		// path if used in combination with the Go kit gRPC Interceptor.
		kitzipkin.HTTPServerTrace(t.ZipkinTracer),
		kithttp.ServerBefore(
			kitot.HTTPToContext(t.OpenTracer, endpointName, logger),
		),
	}
	return clientOptions, serverOptions
}
