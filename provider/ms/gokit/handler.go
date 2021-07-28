package gokit

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/tracing/opentracing"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
)

type ServiceHandler interface {
	GetName() string
	MakeEndpoint(svc interface{}) endpoint.Endpoint
	// Reply message
	Reply() interface{}

	// EncodeDomainRequest client side: convert domain request to grpc request
	EncodeDomainRequest(_ context.Context, request interface{}) (interface{}, error)
	DecodeGrpcReply(_ context.Context, grpcReply interface{}) (interface{}, error)

	// DecodeGrpcRequest server side: convert grpc request to domain request
	DecodeGrpcRequest(_ context.Context, grpcReq interface{}) (interface{}, error)
	EncodeDomainReply(_ context.Context, response interface{}) (interface{}, error)
}

func (gs *GokitServer) NewGrpcHandler(concreteService interface{}, h ServiceHandler) *kitgrpc.Server {
	// 将具体的service和middleware串联起来
	endpoints := h.MakeEndpoint(concreteService)
	for _, m := range gs.Middlewares {
		endpoints = m(endpoints)
	}

	// 添加tracer到ServerBefore
	options := append(gs.Options,
		kitgrpc.ServerBefore(
			opentracing.GRPCToContext(gs.Tracer.OpenTracer, h.GetName(), gs.Logger),
		),
	)

	return kitgrpc.NewServer(
		endpoints,
		h.DecodeGrpcRequest,
		h.EncodeDomainReply,
		options...,
	)
}
