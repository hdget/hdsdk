// Package gokit
package gokit

import (
	"github.com/hdget/sdk/types"
)

// MicroServiceConfig 微服务配置
type MicroServiceConfig struct {
	Name   string             `mapstructure:"name"`
	Client *GokitClientConfig `mapstructure:"client"`
	Server *GokitServerConfig `mapstructure:"server"`

	// middleware
	Trace        *TraceConfig        `mapstructure:"trace"`        // 链路追踪
	CircuitBreak *CircuitBreakConfig `mapstructure:"circuitbreak"` // 熔断
	RateLimit    *RateLimitConfig    `mapstructure:"ratelimit"`    // 限流
}

type MicroServiceImpl struct {
	Name   string
	Logger types.LogProvider
	Config *MicroServiceConfig
	Tracer *Tracer
}

var _ types.MicroService = (*MicroServiceImpl)(nil)

func NewMicroService(logger types.LogProvider, config *MicroServiceConfig) (types.MicroService, error) {
	err := validateConfig(config)
	if err != nil {
		return nil, err
	}

	tracer, err := newTracer(config)
	if err != nil {
		return nil, err
	}

	return &MicroServiceImpl{
		Logger: logger,
		Tracer: tracer,
		Name:   config.Name,
		Config: config,
	}, nil
}

func (msi MicroServiceImpl) GetName() string {
	return msi.Name
}

// 校验配置
func validateConfig(config *MicroServiceConfig) error {
	if config == nil {
		return types.ErrEmptyConfig
	}
	if config.Name == "" {
		return types.ErrInvalidConfig
	}
	return nil
}
