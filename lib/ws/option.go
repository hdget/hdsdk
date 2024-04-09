package ws

import "time"

type ServerOption func(param *ServerParam)

type ServerParam struct {
	publicRouterGroup        *routerGroupParam
	protectRouterGroup       *routerGroupParam
	gracefulShutdownWaitTime time.Duration
}

type routerGroupParam struct {
	Name      string
	UrlPrefix string
}

var (
	defaultServerParams = &ServerParam{
		publicRouterGroup: &routerGroupParam{
			Name:      "public",
			UrlPrefix: "/public",
		},
		protectRouterGroup: &routerGroupParam{
			Name:      "protect",
			UrlPrefix: "/api",
		},
		gracefulShutdownWaitTime: 10 * time.Second,
	}
)

func WithGracefulShutdownWaitTime(waitTime time.Duration) ServerOption {
	return func(param *ServerParam) {
		param.gracefulShutdownWaitTime = waitTime
	}
}
