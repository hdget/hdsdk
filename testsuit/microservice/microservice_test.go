package microservice

import (
	"bytes"
	"github.com/hdget/sdk"
	"github.com/hdget/sdk/testsuit/microservice/pb"
	"github.com/hdget/sdk/testsuit/microservice/service/grpc"
	"github.com/hdget/sdk/types"
	"github.com/hdget/sdk/utils"
	"github.com/hdget/sdk/utils/parallel"
	"testing"
)

const TEST_CONFIG_GOKIT_MICROSERVICE = `
[sdk]
	[sdk.log]
        filename = "demo.log"
		[sdk.log.rotate]
			# 最大保存时间7天(单位hour)
        	max_age = 168
        	# 日志切割时间间隔24小时（单位hour)
        	rotation_time=24
	[sdk.ms]
		[[sdk.ms.items]]
			name = "testservice"
			[sdk.ms.items.trace]
				url = "http://192.168.0.114:9411/api/v2/spans"
			[[sdk.ms.items.servers]]
				type = "grpc"
				address = "0.0.0.0:12345"
				middlewares=["circuitbreak"]
			[[sdk.ms.items.servers]]
				type = "http"
				address = "0.0.0.0:23456"
				middlewares=["ratelimit"]
`

type MicroServiceTestConf struct {
	sdk.Config `mapstructure:",squash"`
}

// nolint:errcheck
func TestMicroService(t *testing.T) {
	v := sdk.LoadConfig("demo", "local", "")

	// merge config from string
	v.MergeConfig(bytes.NewReader(utils.StringToBytes(TEST_CONFIG_GOKIT_MICROSERVICE)))

	// 将配置信息转换成对应的数据结构
	var conf MicroServiceTestConf
	err := v.Unmarshal(&conf)
	if err != nil {
		utils.Fatal("unmarshal demo conf", "err", err)
	}

	err = sdk.Initialize(&conf)
	if err != nil {
		utils.Fatal("sdk initialize", "err", err)
	}

	var group parallel.Group
	grpcTransport := getGrpcTransport()
	group.Add(
		func() error {
			return grpcTransport.Run()
		},
		func(err error) {
			grpcTransport.Close()
		},
	)
	group.Run()
}

func getGrpcTransport() types.MsGrpcServer {
	// 必须手动注册服务实现
	svc := &grpc.SearchServiceImpl{}
	grpcTransport := sdk.MicroService.By("testservice").CreateGrpcServer()
	endpoints := &grpc.GrpcEndpoints{
		SearchEndpoint: grpcTransport.CreateHandler(svc, &grpc.SearchHandler{}),
		HelloEndpoint:  grpcTransport.CreateHandler(svc, &grpc.HelloHandler{}),
	}
	pb.RegisterSearchServiceServer(grpcTransport.GetServer(), endpoints)
	return grpcTransport
}
