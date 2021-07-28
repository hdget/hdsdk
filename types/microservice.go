package types

import "google.golang.org/grpc"

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
	Run()
}

type MsClient interface {
}

// message queue provider
const (
	_              SdkType = SdkCategoryMs + iota
	SdkTypeMsGokit         // 基于Gokit的微服务能力
)
