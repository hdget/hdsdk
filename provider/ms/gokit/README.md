# MicroService

Create and run microservice with help of [hdsdk](https://github.com/hdget/hdsdk) and [hdkit](https://github.com/hdget/hdkit)

It provides out-of-the-box `GRPC` and `HTTP` functionality with only small piece of config.
 
## Table of Contents
- [Usage](#usage)
    - [use hdkit create project](#use-hdkit-create-project-boilerplate)
- [Config](#config)
- [Example](#example)
- [FAQ](#faq)

## Usage

### Use `hdkit` create project boilerplate

Please refer to [hdkit](https://github.com/hdget/hdkit) to create project boilerplate

> Directory `autogen` are files generated automatically each time, all files under this directory will be overwritten if run `hdkit` again

## Config

 ```
[sdk.service]
    [sdk.service.default]
        name = "testservice"
        [[sdk.service.default.servers]]
            type = "grpc"
            address = "0.0.0.0:12345"
            middlewares=["circuitbreak", "ratelimit"]
        [[sdk.service.default.servers]]
            type = "http"
            address = "0.0.0.0:23456"
            ...
    [[sdk.service.items]]
        name = "httpservice"
        [[sdk.service.items.servers]]
            address = "0.0.0.0:12345"
            middlewares=["circuitbreak", "ratelimit"]
            ...        
```

There are two kinds of service config here:
- default: default service, it can be get by `sdk.MicroService.My()`
- items:   service identified by `name`, it can be get by `sdk.MicroService.By(name)`
- type:  
  - grpc: specify server is Grpc server, if not specified `type`, then `grpc` taken as default
  - http: specify server is Http server
  
> One service could have multiple server, like: grpc, http, please make sure `s` exists in `servers`
> Except `default` config, other `items` config must be `[[` and `]]` wrapped

### Grpc server
please set `type=grpc` under `sdk.service.items.server`, now there are two middlewares supported:
- circuitbreak
- ratelimit

### Http server
still under development

## Example
```
svc := &grpc.SearchServiceImpl{}
grpcServer := hdsdk.MicroService.By("testservice").CreateGrpcServer()

endpoints := &grpc.GrpcEndpoints{
	SearchEndpoint: grpcTransport.CreateHandler(svc, &grpc.SearchHandler{}),
	HelloEndpoint:  grpcTransport.CreateHandler(svc, &grpc.HelloHandler{}),
}

pb.RegisterSearchServiceServer(grpcServer.GetServer(), endpoints)

err = hdsdk.Initialize(&conf)
if err != nil {
    utils.Fatal("sdk initialize", "err", err)
}

var group parallel.Group
group.Add(
    func() error {
        return grpcServer.Run()
    },
    func(err error) {
        grpcServer.Close()
    },
)
group.Run()
```

More details please refer to [microservice example](https://github.com/hdget/hdsdk-examples/tree/main/microservice)