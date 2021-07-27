package types

import "google.golang.org/grpc"

type RegisterFunc func(grpcServer *grpc.Server, concreteService interface{})

// MsProvider MS: microservice
type MsProvider interface {
	By(name string) MicroService
}

type MicroService interface {
	GetName() string
	CreateGrpcServer(registerFunc RegisterFunc, concreteService interface{}) MsServer
	CreateGrpcClient() MsClient
}

type MsServer interface {
	Run()
}

type MsClient interface {
}

// message queue provider
const (
	_                 SdkType = SdkCategoryMs + iota
	SdkTypeMsGokit         // 基于Gokit的微服务能力
)
