package types

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"net/http"
)

// MsProvider MS: microservice
type MsProvider interface {
	My() MicroService
	By(name string) MicroService
}

type MicroService interface {
	NewGrpcServerManager() GrpcServerManager
	NewGrpcClientManager() GrpcClientManager

	NewHttpServerManager() HttpServerManager
}

type GrpcServerManager interface {
	GetServer() *grpc.Server
	CreateHandler(concreteService interface{}, ap GrpcAspect) *kitgrpc.Server
	RunServer() error // run server
	Close()           // shutdown server
}

type HttpServerManager interface {
	CreateHandler(concreteService interface{}, ap HttpAspect) http.Handler
	RunServer(handlers map[string]http.Handler) error // run server
	Close()                                           // shutdown server
}

type GrpcClientManager interface {
	CreateConnection(args ...grpc.DialOption) (*grpc.ClientConn, error)
	CreateEndpoint(conn *grpc.ClientConn, ap GrpcAspect) endpoint.Endpoint
}

// GrpcAspect Grpc切面
type GrpcAspect interface {
	GetServiceName() string
	GetMethodName() string

	// MakeEndpoint 解析request, 调用服务函数, 封装成endpoint
	MakeEndpoint(svc interface{}) endpoint.Endpoint
	// ServerDecodeRequest server side convert grpc request to domain request
	ServerDecodeRequest(ctx context.Context, request interface{}) (interface{}, error)
	// ServerEncodeResponse server side encode response to domain response
	ServerEncodeResponse(ctx context.Context, response interface{}) (interface{}, error)

	// GetGrpcReplyType 获取Grpc返回类型
	GetGrpcReplyType() interface{}
	// ClientEncodeRequest server side convert grpc request to domain request
	ClientEncodeRequest(ctx context.Context, request interface{}) (interface{}, error)
	// ClientDecodeResponse server side convert grpc request to domain request
	ClientDecodeResponse(ctx context.Context, request interface{}) (interface{}, error)
}

// HttpAspect HttpAspect
type HttpAspect interface {
	GetMethodName() string
	// MakeEndpoint 解析request, 调用服务函数, 封装成endpoint
	MakeEndpoint(svc interface{}) endpoint.Endpoint
	// ServerDecodeRequest server side convert grpc request to domain request
	ServerDecodeRequest(ctx context.Context, request *http.Request) (interface{}, error)
	// ServerEncodeResponse server side encode response to domain response
	ServerEncodeResponse(context.Context, http.ResponseWriter, interface{}) error
}

// message queue provider
const (
	_              SdkType = SdkCategoryMs + iota
	SdkTypeMsGokit         // 基于Gokit的微服务能力
)
