package grpc

import (
	"context"
	"github.com/hdget/hdsdk/testsuit/microservice/pb"
)

type SearchServiceImpl struct {
}

func (s SearchServiceImpl) Hello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Response: "hello world",
	}, nil
}
func (s SearchServiceImpl) Search(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error) {
	return &pb.SearchResponse{
		Response: "search response",
	}, nil
}
