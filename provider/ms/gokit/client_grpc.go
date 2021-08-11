package gokit

import (
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/hdget/hdsdk/types"
)

// GokitClientConfig 客户端配置
type GokitClientConfig struct {
	Name         string   `mapstructure:"name"`
	ExchangeName string   `mapstructure:"exchange_name"`
	ExchangeType string   `mapstructure:"exchange_type"`
	QueueName    string   `mapstructure:"queue_name"`
	RoutingKeys  []string `mapstructure:"routing_keys"`
}

type GokitClient struct {
	Logger  types.LogProvider
	Options []kitgrpc.ClientOption
}

var _ types.MsGrpcClient = (*GokitClient)(nil)

//
//// CreateGrpcClient producer的名字和route中的名字对应
//func (msi *MicroServiceImpl) CreateGrpcClient() types.MsGrpcClient {
//	clientOptions := make([]kitgrpc.ClientOption, 0)
//	if msi.Tracer != nil {
//		clientOptions = append(clientOptions, kitzipkin.GRPCClientTrace(msi.Tracer.ZipkinTracer))
//	}
//
//	return &GokitClient{
//		Logger:  msi.Logger,
//		Options: clientOptions,
//	}
//}
