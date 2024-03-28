package ws

import (
	"context"
	"errors"
	"github.com/hdget/hdsdk/types"
	"github.com/hdget/hdutils/parallel"
	"net/http"
	"syscall"
)

type httpServerImpl struct {
	*baseServer
}

func NewHttpServer(logger types.LogProvider, address string, options ...ServerOption) (WebServer, error) {
	return &httpServerImpl{
		baseServer: newBaseServer(logger, address, options...),
	}, nil
}

func (w httpServerImpl) Run() {
	var group parallel.Group
	{
		group.Add(w.ListenAndServe, w.shutdown)
	}
	{
		group.Add(
			parallel.SignalActor(
				context.Background(),
				syscall.SIGINT,
				syscall.SIGQUIT,
				syscall.SIGTERM,
				syscall.SIGKILL,
			),
		)
	}

	// 默认endless服务器会监听如下信号:
	// syscall.SIGHUP,syscall.SIGUSR1,syscall.SIGUSR2,syscall.SIGINT,syscall.SIGTERM,syscall.SIGTSTP
	// 接收到syscall.SIGHUP信号将触发`fork/restart`实现优雅重启(kill -1 pid会发送SIGHUP信号）
	// 接收到syscall.SIGINT或syscall.SIGTERM信号将触发优雅关机
	// 接收到syscall.SIGUSR2信号将触发HammerTime
	if err := group.Run(); err != nil && errors.Is(err, http.ErrServerClosed) {
		w.logger.Error("http server quit", "error", err)
	}
}
