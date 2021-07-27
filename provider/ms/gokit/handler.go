package gokit

import (
	"context"
	"github.com/go-kit/kit/endpoint"
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

//func (msi MicroServiceImpl) NewGrpcHandler(genService interface{}, h ServiceHandler) *kitgrpc.Server {
//	tracerOptions := append(msi.Tracer.Options,
//		kitgrpc.ServerBefore(
//			opentracing.GRPCToContext(s.Tracer.OpenTracer, h.GetName(), logger),
//	),
//	)
//	endpoints := h.MakeEndpoint(svc)
//
//	for _, m := range s.Middlewares {
//		endpoints = m(endpoints)
//	}
//
//	return grpctransport.NewServer(
//		endpoints,
//		h.DecodeGrpcRequest,
//		h.EncodeDomainReply,
//		tracerOptions...,
//	)
//}
