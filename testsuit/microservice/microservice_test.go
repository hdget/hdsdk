package microservice

import (
	"bytes"
	"github.com/hdget/hdsdk"
	"github.com/hdget/hdsdk/testsuit/microservice/pb"
	"github.com/hdget/hdsdk/testsuit/microservice/service/grpc"
	"github.com/hdget/hdsdk/utils"
	"github.com/hdget/hdsdk/utils/parallel"
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
	[sdk.service]
		[[sdk.service.items]]
			name = "testservice"
			[[sdk.service.items.servers]]
				type = "grpc"
				address = "0.0.0.0:12345"
				middlewares=["trace", "circuitbreak", "ratelimit"]
`

type MicroServiceTestConf struct {
	hdsdk.Config `mapstructure:",squash"`
}

// nolint:errcheck
func TestMicroService(t *testing.T) {
	v := hdsdk.LoadConfig("demo", "local", "")

	// merge config from string
	v.MergeConfig(bytes.NewReader(utils.StringToBytes(TEST_CONFIG_GOKIT_MICROSERVICE)))

	// 将配置信息转换成对应的数据结构
	var conf MicroServiceTestConf
	err := v.Unmarshal(&conf)
	if err != nil {
		utils.Fatal("unmarshal demo conf", "err", err)
	}

	err = hdsdk.Initialize(&conf)
	if err != nil {
		utils.Fatal("sdk initialize", "err", err)
	}

	ms := hdsdk.MicroService.By("testservice")
	if ms == nil {
		utils.Fatal("get microservice instance", "err", err)
	}

	server := ms.CreateGrpcServer()
	if ms == nil {
		utils.Fatal("get grpc server", "err", err)
	}

	// 必须手动注册服务实现
	svc := &grpc.SearchServiceImpl{}
	endpoints := &grpc.GrpcEndpoints{
		SearchEndpoint: server.CreateHandler(svc, &grpc.SearchHandler{}),
		HelloEndpoint:  server.CreateHandler(svc, &grpc.HelloHandler{}),
	}

	pb.RegisterSearchServiceServer(server.GetServer(), endpoints)

	var group parallel.Group
	group.Add(
		func() error {
			return server.Run()
		},
		func(err error) {
			server.Close()
		},
	)
	err = group.Run()
	if err != nil {
		utils.Fatal("microservice exit", "err", err)
	}
}
