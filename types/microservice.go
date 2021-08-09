package types

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	kithttp "github.com/go-kit/kit/transport/http"
	"google.golang.org/grpc"
	"net/http"
)

// MsProvider MS: microservice
type MsProvider interface {
	My() MicroService
	By(name string) MicroService
}

type MicroService interface {
	GetName() string
	CreateGrpcServer() MsGrpcServer
	CreateHttpServer() MsHttpServer
}

type MsGrpcServer interface {
	GetServer() *grpc.Server
	CreateHandler(concreteService interface{}, ge GrpcEndpoint) *kitgrpc.Server
	Run() error
	Close()
}

type MsGrpcClient interface {
}

type MsHttpServer interface {
	CreateHandler(concreteService interface{}, he HttpEndpoint) *kithttp.Server
	Run(handlers map[string]*kithttp.Server) error
	Close()
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

// GrpcEndpoint Grpc端点的实现
type GrpcEndpoint interface {
	GetName() string
	// MakeEndpoint 解析request, 调用服务函数, 封装成endpoint
	MakeEndpoint(svc interface{}) endpoint.Endpoint
	// ServerDecodeRequest server side convert grpc request to domain request
	ServerDecodeRequest(ctx context.Context, request interface{}) (interface{}, error)
	// ServerEncodeResponse server side encode response to domain response
	ServerEncodeResponse(ctx context.Context, response interface{}) (interface{}, error)
}

// HttpEndpoint HttpEndpoint
type HttpEndpoint interface {
	GetName() string
	// MakeEndpoint 解析request, 调用服务函数, 封装成endpoint
	MakeEndpoint(svc interface{}) endpoint.Endpoint
	// ServerDecodeRequest server side convert grpc request to domain request
	ServerDecodeRequest(context.Context, *http.Request) (interface{}, error)
	// ServerEncodeResponse server side encode response to domain response
	ServerEncodeResponse(context.Context, http.ResponseWriter, interface{}) error
}

// message queue provider
const (
	_              SdkType = SdkCategoryMs + iota
	SdkTypeMsGokit         // 基于Gokit的微服务能力
)
