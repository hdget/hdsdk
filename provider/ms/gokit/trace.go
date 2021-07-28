package gokit

import (
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/reporter"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/pkg/errors"
)

// TraceConfig 服务跟踪选项
type TraceConfig struct {
	Url     string // Zipkin tracing HTTP reporter URL
	Address string
}

type Tracer struct {
	Reporter     reporter.Reporter
	ZipkinTracer *zipkin.Tracer     // zipkin tracer
	OpenTracer   opentracing.Tracer // opentracing tracer
}

var (
	defaultTraceConfig = &TraceConfig{
		Url:     "http://localhost:9411/api/v2/spans",
		Address: "",
	}
	ErrInvalidTraceConfig = errors.New("invalid trace config")
)

func newTracer(config *MicroServiceConfig) (*Tracer, error) {
	// 如果没有配置tracer不做任何事
	traceConfig := config.Trace
	if traceConfig == nil {
		traceConfig = defaultTraceConfig
	}

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
	openTracer := zipkintracer.Wrap(zipkinTracer)

	// optionally set as Global OpenTracing openTracer instance
	opentracing.SetGlobalTracer(openTracer)

	return &Tracer{
		Reporter:     reporter,
		ZipkinTracer: zipkinTracer,
		OpenTracer:   openTracer,
	}, nil
}
