package gokit

import (
	"github.com/hdget/hdsdk/types"
)

// ClientConfig 客户端配置
type ClientConfig struct {
	Transport   string   `mapstructure:"transport"`
	Address     string   `mapstructure:"address"`
	Middlewares []string `mapstructure:"middlewares"`
}

type BaseGokitClient struct {
	Logger      types.LogProvider
	Config      *ClientConfig
	Middlewares []*MsMiddleware
}
