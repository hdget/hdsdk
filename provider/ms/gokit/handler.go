package gokit

import (
	"github.com/go-kit/kit/tracing/opentracing"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/hdget/sdk/types"
)

func (gs *GokitServer) CreateEndpointServer(concreteService interface{}, eh types.EndpointHandler) *kitgrpc.Server {
	// 将具体的service和middleware串联起来
	endpoints := eh.MakeEndpoint(concreteService)
	for _, m := range gs.Middlewares {
		endpoints = m(endpoints)
	}

	// 添加tracer到ServerBefore
	options := append(gs.Options,
		kitgrpc.ServerBefore(
			opentracing.GRPCToContext(gs.Tracer.OpenTracer, eh.GetName(), gs.Logger),
		),
	)

	return kitgrpc.NewServer(
		endpoints,
		eh.ServerDecodeRequest,
		eh.ServerEncodeResponse,
		options...,
	)
}
