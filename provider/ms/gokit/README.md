# MicroService

Create and run microservice with help of [hdsdk](https://github.com/hdget/hdsdk) and [hdkit](https://github.com/hdget/hdkit)

It provides out-of-the-box `GRPC` and `HTTP` functionality with only small piece of config.
 
## Table of Contents
- [Usage](#usage)
    - [use hdkit create project](#use-hdkit-create-project-boilerplate)
- [Config](#config)
    - [middleware trace config](#middleware-trace-config)
    - [middleware circuit break config](#middleware-circuit-break-config)
    - [middleware rate limit config](#middleware-rate-limit-config)
    - [server config](#server-config)
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
        [[sdk.service.default.clients]]
            transport = "grpc"
            address = "0.0.0.0:12345"
            middlewares=["circuitbreak", "ratelimit"]
        [[sdk.service.default.servers]]
            transport = "grpc"
            address = "0.0.0.0:12345"
            middlewares=["circuitbreak", "ratelimit"]
        [[sdk.service.default.servers]]
            transport = "http"
            address = "0.0.0.0:23456"
            middlewares=["trace", "circuitbreak", "ratelimit"]
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
  
> One service could have multiple server, like: grpc, http, please make sure `s` exists in `servers`
> Except `default` config, other `items` config must be `[[` and `]]` wrapped

### middleware trace config
trace config can be ignored, which will be used default trace config instead

```
[sdk.service]
	[[sdk.service.items]]
		name = "testservice"
	    [sdk.service.items.trace]
	      url = "http://localhost:9411/api/v2/spans"
	      address = "localhost:80"	
```

- url: zipkin http report url
- address: open tracer address

### middleware circuit break config
circuitbreak config can be ignored, which will be used default circuit break config instead

```
[sdk.service]
	[[sdk.service.items]]
		name = "testservice"
	    [sdk.service.items.circuitbreak]
	        requests = 10,
	        interval = 0
		    timeout = 60
		    max_requests =  100
		    failure_ratio = 1.0
```

- requests: successive requests
- interval: what's the period it reset timeout counter when in shutdown status
- timeout: time between half open status to open status
- max_requests: when in half open status, if max requests is 0, only one request can be allowed
- failure_ratio: request failure rate

### middleware rate limit config
ratelimit config can be ignored, which will be used default rate limit config instead

```
[sdk.service]
	[[sdk.service.items]]
		name = "testservice"
	    [sdk.service.items.ratelimit]
	        limit = 30
		    burst = 50
```

- limit: how many requests allowed in one second
- burst: max burst requests 

### server config
```
[sdk.service]
	[[sdk.service.items]]
		name = "testservice"
    [[sdk.service.items.servers]]
	    type = "http"
	    address = "0.0.0.0:12345"
		middlewares=["trace", "circuitbreak", "ratelimit"]
```

- type: what's the server type to serve the service, now supports: `grpc` and `http`. If not specified this item, then `grpc` taken as default
- address: what's the server listen address
- middlewares: service middlewares, now it supports:
  - trace
  - circuitbreak
  - ratelimit

  > Note: if specified middleware type here, then it inherits the config from above 

## http client

It can access http transport by using following url:

`/<service name>/<method name>` all use snake case format

The parameters sent by json body, it supports `GET` and `POST`

## Example
```
svc := &grpc.SearchServiceImpl{}
manager := hdsdk.MicroService.By("testservice").NewGrpcServerManager()

endpoints := &grpc.GrpcEndpoints{
	SearchEndpoint: manager.CreateHandler(svc, &grpc.SearchHandler{}),
	HelloEndpoint:  manager.CreateHandler(svc, &grpc.HelloHandler{}),
}

pb.RegisterSearchServiceServer(manager.GetServer(), endpoints)

err = hdsdk.Initialize(&conf)
if err != nil {
    utils.Fatal("sdk initialize", "err", err)
}

var group parallel.Group
group.Add(
    func() error {
        return manager.RunServer()
    },
    func(err error) {
        manager.Close()
    },
)

err := group.Run()
if err != nil {
  utils.Fatal("grpc server exist", "err", err)
}
```

More details please refer to [microservice example](https://github.com/hdget/hdsdk-examples/tree/main/microservice)