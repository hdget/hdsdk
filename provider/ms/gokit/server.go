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

type BaseGokitServer struct {
	Name        string
	Config      *ServerConfig
	Logger      types.LogProvider
	Middlewares []*MsMiddleware
	ctx         context.Context
	cancel      context.CancelFunc
}

// Close 关闭GrpcServer
func (s *BaseGokitServer) Close() {
	s.cancel()
}
