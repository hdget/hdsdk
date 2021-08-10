package grpc

import (
	"context"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/hdget/hdsdk/testsuit/microservice/pb"
)

type GrpcEndpoints struct {
	SearchEndpoint kitgrpc.Handler
	HelloEndpoint  kitgrpc.Handler
}

func (e GrpcEndpoints) Search(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error) {
	_, resp, err := e.SearchEndpoint.ServeGRPC(ctx, request)
	return resp.(*pb.SearchResponse), err
}

func (e GrpcEndpoints) Hello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloResponse, error) {
	_, resp, err := e.HelloEndpoint.ServeGRPC(ctx, request)
	return resp.(*pb.HelloResponse), err
}
