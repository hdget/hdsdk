package hdsdk

import (
	"github.com/hdget/hdsdk/v1/config"
	"github.com/hdget/hdsdk/v1/errdef"
	"github.com/hdget/hdsdk/v1/intf"
	"github.com/hdget/hdsdk/v1/provider/logger/zerolog"
	"github.com/hdget/hdutils"
	"go.uber.org/fx"
)

type SdkInstance struct {
	configLoader intf.ConfigLoader
	fxOptions    []fx.Option
	logger       intf.LoggerProvider
	db           intf.DbProvider
	graph        intf.GraphProvider
	cache        intf.CacheProvider
}

var (
	_instance *SdkInstance
)

// New 默认包括logger
func New(app, env string) *SdkInstance {
	if _instance == nil {
		_instance = newWithAppEnv(app, env)

		// 所有的provider依赖configProvider和loggerProvider
		_instance.fxOptions = []fx.Option{
			fx.NopLogger,
			fx.Provide(func() intf.ConfigLoader { return _instance.configLoader }), // 提供configProvider
			zerolog.Capability.Module,      // 提供zerolog capability
			fx.Populate(&_instance.logger), // 初始化logger
		}
	}
	return _instance
}

// Empty 不包括logger
func Empty(app, env string) *SdkInstance {
	if _instance == nil {
		_instance = newWithAppEnv(app, env)

		// 所有的provider依赖configProvider和loggerProvider
		_instance.fxOptions = []fx.Option{
			fx.NopLogger,
			fx.Provide(func() intf.ConfigLoader { return _instance.configLoader }), // 提供configProvider
		}
	}
	return _instance
}

func WithConfigLoader(configLoader intf.ConfigLoader) *SdkInstance {
	if _instance != nil {
		_instance.configLoader = configLoader
		return _instance
	}

	_instance = newWithConfigLoader(configLoader)
	return _instance
}

func (i *SdkInstance) LoadConfig(configVar any) *SdkInstance {
	if i != nil {
		_ = i.configLoader.Unmarshal(configVar)
	}
	return i
}

// Initialize 初始化指定的能力
func (i *SdkInstance) Initialize(capabilities ...*intf.Capability) error {
	for _, c := range capabilities {
		switch c.Category {
		case intf.ProviderCategoryLogger:
			i.fxOptions = append(i.fxOptions, c.Module, fx.Populate(&_instance.logger))
		case intf.ProviderCategoryDb:
			i.fxOptions = append(i.fxOptions, c.Module, fx.Populate(&_instance.db))
		case intf.ProviderCategoryCache:
			i.fxOptions = append(i.fxOptions, c.Module, fx.Populate(&_instance.cache))
		case intf.ProviderCategoryNeo4j:
			i.fxOptions = append(i.fxOptions, c.Module, fx.Populate(&_instance.graph))
		default:
			return errdef.ErrInvalidCapability
		}
	}

	_ = fx.New(
		i.fxOptions...,
	)

	return nil
}

func newWithAppEnv(app, env string) *SdkInstance {
	// 初始化configLoader
	configLoader, err := config.New(app, env)
	if configLoader == nil {
		hdutils.LogFatal("new config loader", "err", err)
	}

	return newWithConfigLoader(configLoader)
}

func newWithConfigLoader(configLoader intf.ConfigLoader) *SdkInstance {
	return &SdkInstance{
		configLoader: configLoader,
	}
}
