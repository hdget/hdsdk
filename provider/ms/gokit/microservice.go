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
	// servers
	Servers []*ServerConfig `mapstructure:"servers"`
}

type MicroServiceImpl struct {
	Name   string
	Logger types.LogProvider
	Config *MicroServiceConfig
}

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

func (msi MicroServiceImpl) GetName() string {
	return msi.Name
}

func (msi MicroServiceImpl) GetServerConfig(serverType string) *ServerConfig {
	for _, serverConfig := range msi.Config.Servers {
		configServerType := serverConfig.ServerType

		// if we don't specify the `type` in config file
		// if set to be `GRPC_SERVER` by default
		if configServerType == "" {
			configServerType = GRPC_SERVER
		}

		if configServerType == serverType {
			return serverConfig
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
