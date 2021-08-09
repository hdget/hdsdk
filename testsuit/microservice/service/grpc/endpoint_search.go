package grpc

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/hdget/sdk/testsuit/microservice/pb"
)

// put into handler file
type SearchHandler struct{}

func (h SearchHandler) GetName() string {
	return "search"
}

func (s SearchHandler) MakeEndpoint(svc interface{}) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return svc.(*SearchServiceImpl).Search(ctx, request.(*pb.SearchRequest))
	}
}

func (s SearchHandler) ServerDecodeRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	return grpcReq.(*pb.SearchRequest), nil
}

func (s SearchHandler) ServerEncodeResponse(ctx context.Context, response interface{}) (interface{}, error) {
	return response.(*pb.SearchResponse), nil
}
