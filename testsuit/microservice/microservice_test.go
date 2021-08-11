package microservice

import (
	"bytes"
	"context"
	"fmt"
	"github.com/hdget/hdsdk"
	gengrpc "github.com/hdget/hdsdk/testsuit/microservice/autogen/grpc"
	"github.com/hdget/hdsdk/testsuit/microservice/autogen/pb"
	"github.com/hdget/hdsdk/testsuit/microservice/service"
	"github.com/hdget/hdsdk/utils"
	"github.com/hdget/hdsdk/utils/parallel"
	"google.golang.org/grpc"
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
			[[sdk.service.items.clients]]
				type = "grpc"
				address = "0.0.0.0:12345"
				middlewares=["trace", "circuitbreak", "ratelimit"]
			[[sdk.service.items.servers]]
				type = "grpc"
				address = "0.0.0.0:12345"
				middlewares=["trace", "circuitbreak", "ratelimit"]
`

type MicroServiceTestConf struct {
	hdsdk.Config `mapstructure:",squash"`
}

// nolint:errcheck
func TestMsServer(t *testing.T) {
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

	manager := ms.NewGrpcServerManager()
	if manager == nil {
		utils.Fatal("new grpc server manager", "err", err)
	}

	// 必须手动注册服务实现
	svc := &service.SearchServiceImpl{}
	handlers := &gengrpc.Handlers{
		SearchHandler: manager.CreateHandler(svc, &gengrpc.SearchAspect{}),
		HelloHandler:  manager.CreateHandler(svc, &gengrpc.HelloAspect{}),
	}

	pb.RegisterSearchServiceServer(manager.GetServer(), handlers)

	var group parallel.Group
	group.Add(
		func() error {
			return manager.RunServer()
		},
		func(err error) {
			manager.Close()
		},
	)
	err = group.Run()
	if err != nil {
		utils.Fatal("microservice exit", "err", err)
	}
}

// nolint:errcheck
func TestMsClient(t *testing.T) {
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

	conn, err := grpc.DialContext(context.Background(), "0.0.0.0:12345", grpc.WithInsecure())
	if err != nil {
		utils.Fatal("create grpc connection", "err", err)
	}
	defer conn.Close()

	client, err := gengrpc.NewClient(conn, "testservice")
	if err != nil {
		utils.Fatal("create grpc client", "err", err)
	}

	result, err := client.Search(context.Background(), &pb.SearchRequest{Request: "ddd"})
	if err != nil {
		utils.Fatal("call method", "method", "Search", "err", err)
	}
	fmt.Println(result)
}
