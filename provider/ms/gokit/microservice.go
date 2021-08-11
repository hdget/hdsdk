// Package gokit
package gokit

import (
	"github.com/hdget/hdsdk/types"
)

// MicroServiceConfig 微服务配置
type MicroServiceConfig struct {
	Name string `mapstructure:"name"`
	// middlewares
	Trace        *TraceConfig        `mapstructure:"trace"`        // 链路追踪
	CircuitBreak *CircuitBreakConfig `mapstructure:"circuitbreak"` // 熔断
	RateLimit    *RateLimitConfig    `mapstructure:"ratelimit"`    // 限流
	// clients and servers
	Clients []*ClientConfig `mapstructure:"clients"`
	Servers []*ServerConfig `mapstructure:"servers"`
}

type MicroServiceImpl struct {
	Name   string
	Logger types.LogProvider
	Config *MicroServiceConfig
}

// transport type
const (
	GRPC = "grpc"
	HTTP = "http"
)

var _ types.MicroService = (*MicroServiceImpl)(nil)

func NewMicroService(logger types.LogProvider, config *MicroServiceConfig) (types.MicroService, error) {
	err := validateConfig(config)
	if err != nil {
		return nil, err
	}

	return &MicroServiceImpl{
		Logger: logger,
		Name:   config.Name,
		Config: config,
	}, nil
}

func (msi MicroServiceImpl) GetServerConfig(transport string) *ServerConfig {
	for _, serverConfig := range msi.Config.Servers {
		configTransport := serverConfig.Transport

		// if we don't specify the `type` in config file
		// if set to be `GRPC` by default
		if configTransport == "" {
			configTransport = GRPC
		}

		if configTransport == transport {
			return serverConfig
		}
	}

	return nil
}

func (msi MicroServiceImpl) GetClientConfig(transport string) *ClientConfig {
	for _, clientConfig := range msi.Config.Clients {
		configTransport := clientConfig.Transport

		// if we don't specify the `type` in config file
		// if set to be `GRPC` by default
		if configTransport == "" {
			configTransport = GRPC
		}

		if configTransport == transport {
			return clientConfig
		}
	}

	return nil
}

// 校验配置
func validateConfig(config *MicroServiceConfig) error {
	if config == nil {
		return types.ErrEmptyConfig
	}
	if config.Name == "" || len(config.Servers) == 0 {
		return types.ErrInvalidConfig
	}
	return nil
}
