package types

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

// MsProvider MS: microservice
type MsProvider interface {
	By(name string) MicroService
}

type MicroService interface {
	GetName() string
	CreateServer() MsServer
	CreateClient() MsClient
}

type MsServer interface {
	GetGrpcServer() *grpc.Server
	CreateEndpointServer(concreteService interface{}, eh EndpointHandler) *kitgrpc.Server
	Run()
}

type MsClient interface {
}

//type EndpointHandler interface {
//	GetName() string
//
//	// MakeEndpoint 解析request, 调用服务函数, 封装成endpoint
//	MakeEndpoint(svc interface{}) endpoint.Endpoint
//
//	// ClientEncodeRequest client side convert domain request to grpc request
//	ClientEncodeRequest(ctx context.Context, request interface{}) (interface{}, error)
//	// ClientDecodeResponse client side decode grpc response to domain response
//	ClientDecodeResponse(ctx context.Context, grpcReply interface{}) (interface{}, error)
//	// ServerDecodeRequest server side convert grpc request to domain request
//	ServerDecodeRequest(ctx context.Context, grpcReq interface{}) (interface{}, error)
//	// ServerEncodeResponse server side encode response to domain response
//	ServerEncodeResponse(ctx context.Context, response interface{}) (interface{}, error)
//}

type EndpointHandler interface {
	GetName() string
	// MakeEndpoint 解析request, 调用服务函数, 封装成endpoint
	MakeEndpoint(svc interface{}) endpoint.Endpoint
	// ServerDecodeRequest server side convert grpc request to domain request
	ServerDecodeRequest(ctx context.Context, grpcReq interface{}) (interface{}, error)
	// ServerEncodeResponse server side encode response to domain response
	ServerEncodeResponse(ctx context.Context, response interface{}) (interface{}, error)
}

// message queue provider
const (
	_              SdkType = SdkCategoryMs + iota
	SdkTypeMsGokit         // 基于Gokit的微服务能力
)
