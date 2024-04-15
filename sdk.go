package hdsdk

import (
	"github.com/hdget/hdsdk/v2/errdef"
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/hdget/hdsdk/v2/provider/config/viper"
	"github.com/hdget/hdsdk/v2/provider/logger/zerolog"
	"github.com/hdget/hdsdk/v2/types"
	"go.uber.org/fx"
)

type SdkInstance struct {
	app            string
	env            string
	configVar      any  // external config var, if set then need load config stuff into it
	debug          bool // debug mode
	configProvider intf.ConfigProvider
	logger         intf.LoggerProvider
	db             intf.DbProvider
	graph          intf.GraphProvider
	redis          intf.RedisProvider
}

var (
	_instance *SdkInstance
)

// New 默认包括logger
func New(app, env string, options ...Option) *SdkInstance {
	if _instance == nil {
		_instance = &SdkInstance{
			app: app,
			env: env,
		}
	}

	for _, option := range options {
		option(_instance)
	}

	return _instance
}

// Initialize all kinds of capability
func (i *SdkInstance) Initialize(capabilities ...*types.Capability) error {
	configInitialized := false
	loggerInitialized := false

	fxOptions := make([]fx.Option, 0)
	for _, c := range capabilities {
		switch c.Category {
		case types.ProviderCategoryConfig:
			fxOptions = append(fxOptions,
				fx.Provide(func() *types.ConfigArgument { return &types.ConfigArgument{App: i.app, Env: i.env} }),
				c.Module,
				fx.Populate(&_instance.configProvider),
			)
			// mark config provider had been initialized
			configInitialized = true
		case types.ProviderCategoryLogger:
			fxOptions = append(fxOptions, c.Module, fx.Populate(&_instance.logger))
			// mark logger provider had been initialized
			loggerInitialized = true
		case types.ProviderCategoryDb:
			fxOptions = append(fxOptions, c.Module, fx.Populate(&_instance.db))
		case types.ProviderCategoryRedis:
			fxOptions = append(fxOptions, c.Module, fx.Populate(&_instance.redis))
		case types.ProviderCategoryGraph:
			fxOptions = append(fxOptions, c.Module, fx.Populate(&_instance.graph))
		default:
			return errdef.ErrInvalidCapability
		}
	}

	// if config provider not initialized
	if !configInitialized {
		configProvider, err := viper.New(&types.ConfigArgument{App: i.app, Env: i.env})
		if err != nil {
			return err
		}

		fxOptions = append(fxOptions,
			fx.Provide(func() intf.ConfigProvider { return configProvider }),
			fx.Populate(&_instance.configProvider),
		)
	}

	// if logger provider is not initialized, use default logger
	if !loggerInitialized {
		fxOptions = append(fxOptions, zerolog.Capability.Module, fx.Populate(&_instance.logger))
	}

	// in product mode disable fx internal logger
	if !i.debug {
		fxOptions = append(fxOptions, fx.NopLogger)
	}

	_ = fx.New(
		fxOptions...,
	)

	// if configVar set then load config stuff to configVar
	if i.configVar != nil && _instance.configVar != nil {
		return _instance.configProvider.Unmarshal(i.configVar)
	}

	return nil
}
