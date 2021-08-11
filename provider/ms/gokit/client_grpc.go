package gokit

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/hdget/hdsdk/types"
	"google.golang.org/grpc"
)

type GrpcClientManager struct {
	*BaseGokitClient
	Options []kitgrpc.ClientOption
}

var _ types.GrpcClientManager = (*GrpcClientManager)(nil)

// NewGrpcClientManager 创建微服务client manager
func (msi MicroServiceImpl) NewGrpcClientManager() types.GrpcClientManager {
	// IMPORTANT: 这里使用MsMiddleware来封装gokit的middleware
	mdws := make([]*MsMiddleware, 0)
	serverConfig := msi.GetServerConfig(GRPC)
	for _, mdwName := range serverConfig.Middlewares {
		newFunc := NewMdwFunctions[mdwName]
		if newFunc != nil {
			mdws = append(mdws, newFunc(msi.Config))
		}
	}

	return &GrpcClientManager{
		BaseGokitClient: &BaseGokitClient{
			Logger:      msi.Logger,
			Middlewares: mdws,
		},
		Options: make([]kitgrpc.ClientOption, 0),
	}
}

func (cm GrpcClientManager) CreateConnection(args ...grpc.DialOption) (*grpc.ClientConn, error) {
	conn, err := grpc.DialContext(context.Background(), cm.Config.Address, args...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (cm GrpcClientManager) CreateEndpoint(conn *grpc.ClientConn, ap types.GrpcAspect) endpoint.Endpoint {
	// 有些中间件只需要处理option
	for _, m := range cm.Middlewares {
		if len(m.InjectFunctions) == 0 {
			continue
		}

		if injectFunc := m.InjectFunctions[GRPC]; injectFunc != nil {
			// 从inject function中获取client options
			clientOptions, _ := injectFunc(cm.Logger, ap.GetServiceName())
			for _, option := range clientOptions {
				clientOption, ok := option.(kitgrpc.ClientOption)
				if ok {
					cm.Options = append(cm.Options, clientOption)
				}
			}
		}
	}

	// client作为endpoints的第一个endpoint
	endpoints := kitgrpc.NewClient(
		conn,
		ap.GetServiceName(),     // service name, e,g: pb.SearchService
		ap.GetMethodName(),      // method name: Search
		ap.ClientEncodeRequest,  // ap.ClientEncodeRequest
		ap.ClientDecodeResponse, // ap.ClientDecodeResponse
		ap.GetGrpcReplyType(),   // ap.Rely()
		cm.Options...,
	).Endpoint()

	// 处理其他中间件
	for _, m := range cm.Middlewares {
		if m.Middleware != nil {
			endpoints = m.Middleware(endpoints)
		}
	}

	return endpoints
}
