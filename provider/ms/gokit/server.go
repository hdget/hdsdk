package gokit

import (
	"context"
	"github.com/hdget/hdsdk/types"
)

// ServerConfig 服务端配置
type ServerConfig struct {
	ServerType  string   `mapstructure:"type"`
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

const (
	GRPC_SERVER = "grpc"
	HTTP_SERVER = "http"
)

// Close 关闭GrpcServer
func (bgs *BaseGokitServer) Close() {
	bgs.cancel()
}
