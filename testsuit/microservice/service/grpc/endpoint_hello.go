package grpc

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/hdget/sdk/testsuit/microservice/pb"
)

type HelloHandler struct{}

func (h HelloHandler) GetName() string {
	return "hello"
}

func (h HelloHandler) MakeEndpoint(svc interface{}) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return svc.(*SearchServiceImpl).Hello(ctx, request.(*pb.HelloRequest))
	}
}

func (h HelloHandler) ServerDecodeRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	panic("implement me")
}

func (h HelloHandler) ServerEncodeResponse(ctx context.Context, response interface{}) (interface{}, error) {
	panic("implement me")
}
