package ws

import "time"

type ServerOption func(param *ServerParam)

type ServerParam struct {
	gracefulShutdownWaitTime time.Duration
}

var (
	defaultServerParams = &ServerParam{
		gracefulShutdownWaitTime: 10 * time.Second,
	}
)

func WithGracefulShutdownWaitTime(waitTime time.Duration) ServerOption {
	return func(param *ServerParam) {
		param.gracefulShutdownWaitTime = waitTime
	}
}
