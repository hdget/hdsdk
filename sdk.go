package hdsdk

import (
	"context"
	"github.com/hdget/hdsdk/v2/errdef"
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/hdget/hdsdk/v2/provider/config/viper"
	"github.com/hdget/hdsdk/v2/provider/logger/zerolog"
	"github.com/hdget/hdutils/logger"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

type SdkInstance struct {
	debug          bool // debug mode
	configProvider intf.ConfigProvider
	logger         intf.LoggerProvider
	db             intf.DbProvider
	sqlxDb         intf.SqlxDbProvider
	graph          intf.GraphProvider
	redis          intf.RedisProvider
	mq             intf.MqProvider
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
	_instance.configProvider, err = viper.New(app, env)
	if err != nil {
		logger.LogFatal("new default config provider", "err", err)
	}

	return _instance
}

func HasInitialized() bool {
	return _instance != nil
}

func (i *SdkInstance) LoadConfig(configVar any) *SdkInstance {
	if i.configProvider != nil {
		// if config provider is already provided in New, ignore it
		err := i.configProvider.Unmarshal(configVar)
		if err != nil {
			logger.LogError("unmarshal to config var", "err", err)
		}
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
		case intf.ProviderCategoryDbSqlx: // will removed in the future
			fxOptions = append(fxOptions, c.Module, fx.Populate(&_instance.sqlxDb))
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

	err := fx.New(fxOptions...).Start(context.Background())
	if err != nil {
		return err
	}

	return nil
}
