package hdsdk

import (
	"github.com/hdget/hdsdk/v2/errdef"
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/hdget/hdsdk/v2/provider/config/viper"
	"github.com/hdget/hdsdk/v2/provider/logger/zerolog"
	"github.com/hdget/hdutils"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

type SdkInstance struct {
	debug          bool // debug mode
	configProvider intf.ConfigProvider
	logger         intf.LoggerProvider
	db             intf.DbProvider
	graph          intf.GraphProvider
	redis          intf.RedisProvider
	mq             intf.MqProvider
}

var (
	_instance *SdkInstance
)

// New 默认包括logger
func New(options ...Option) *SdkInstance {
	if _instance == nil {
		_instance = &SdkInstance{}
	}

	for _, option := range options {
		option(_instance)
	}

	return _instance
}

func (i *SdkInstance) LoadConfig(configVar any) *SdkInstance {
	if i.configProvider != nil {
		// if config provider is already provided in New, ignore it
		err := i.configProvider.Unmarshal(configVar)
		if err != nil {
			hdutils.LogError("unmarshal to config var", "err", err)
		}
	}
	return i
}

func (i *SdkInstance) UseDefaultConfigProvider(app, env string) *SdkInstance {
	var err error
	i.configProvider, err = viper.New(app, env)
	if err != nil {
		hdutils.LogError("new default config provider", "err", err)
	}
	return i
}

// Initialize all kinds of capability
func (i *SdkInstance) Initialize(capabilities ...*intf.Capability) error {
	if i.configProvider == nil {
		return errdef.ErrConfigProviderNotReady
	}

	loggerInitialized := false
	fxOptions := []fx.Option{
		fx.Provide(func() intf.ConfigProvider { return i.configProvider }),
	}
	for _, c := range capabilities {
		switch c.Category {
		case intf.ProviderCategoryLogger:
			fxOptions = append(fxOptions, c.Module, fx.Populate(&_instance.logger))
			// mark logger provider had been initialized
			loggerInitialized = true
		case intf.ProviderCategoryDb:
			fxOptions = append(fxOptions, c.Module, fx.Populate(&_instance.db))
		case intf.ProviderCategoryRedis:
			fxOptions = append(fxOptions, c.Module, fx.Populate(&_instance.redis))
		case intf.ProviderCategoryGraph:
			fxOptions = append(fxOptions, c.Module, fx.Populate(&_instance.graph))
		case intf.ProviderCategoryMq:
			fxOptions = append(fxOptions, c.Module, fx.Populate(&_instance.mq))
		default:
			return errors.Wrapf(errdef.ErrInvalidCapability, "capability: %s", c.Name)
		}
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

	return nil
}
