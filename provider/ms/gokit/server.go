package gokit

import (
	"context"
	"github.com/hdget/hdsdk/types"
)

// ServerConfig 服务端配置
type ServerConfig struct {
	Transport   string   `mapstructure:"transport"`
	Address     string   `mapstructure:"address"`
	Middlewares []string `mapstructure:"middlewares"`
}

type BaseServerManager struct {
	Name        string
	Config      *ServerConfig
	Logger      types.LogProvider
	Middlewares []*MsMiddleware
	ctx         context.Context
	cancel      context.CancelFunc
}

type errorHandler struct {
	Logger types.LogProvider
}

// Close 关闭GrpcServer
func (s *BaseServerManager) Close() {
	s.cancel()
}

func (h errorHandler) Handle(ctx context.Context, err error) {
	h.Logger.Error("encounter error", "error", err)
}
