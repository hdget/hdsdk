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
	debug bool // debug mode

}

var (
	_instance *SdkInstance
)

func New(app, env string, options ...Option) *SdkInstance {
	if _instance == nil {
		_instance = &SdkInstance{}
	}

	for _, option := range options {
		option(_instance)
	}

	var err error
	_configProvider, err = viper.New(app, env)
	if err != nil {
		hdutils.LogFatal("new default config provider", "err", err)
	}

	return _instance
}

func (i *SdkInstance) LoadConfig(configVar any) *SdkInstance {
	if _configProvider != nil {
		// if config provider is already provided in New, ignore it
		err := _configProvider.Unmarshal(configVar)
		if err != nil {
			hdutils.LogError("unmarshal to config var", "err", err)
		}
	}
	return i
}

// Initialize all kinds of capability
func (i *SdkInstance) Initialize(capabilities ...*intf.Capability) error {
	if _configProvider == nil {
		return errdef.ErrConfigProviderNotReady
	}

	loggerInitialized := false
	fxOptions := []fx.Option{
		fx.Provide(func() intf.ConfigProvider { return _configProvider }),
	}
	for _, c := range capabilities {
		switch c.Category {
		case intf.ProviderCategoryLogger:
			fxOptions = append(fxOptions, c.Module, fx.Populate(&_logger))
			// mark logger provider had been initialized
			loggerInitialized = true
		case intf.ProviderCategoryDb:
			fxOptions = append(fxOptions, c.Module, fx.Populate(&_db))
		case intf.ProviderCategoryRedis:
			fxOptions = append(fxOptions, c.Module, fx.Populate(&_redis))
		case intf.ProviderCategoryGraph:
			fxOptions = append(fxOptions, c.Module, fx.Populate(&_graph))
		case intf.ProviderCategoryMq:
			fxOptions = append(fxOptions, c.Module, fx.Populate(&_mq))
		default:
			return errors.Wrapf(errdef.ErrInvalidCapability, "capability: %s", c.Name)
		}
	}

	// if logger provider is not initialized, use default logger
	if !loggerInitialized {
		fxOptions = append(fxOptions, zerolog.Capability.Module, fx.Populate(&_logger))
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
