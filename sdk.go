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
	"sync"
)

type SdkInstance struct {
	option *optionObject

	configProvider intf.ConfigProvider
	logger         intf.LoggerProvider
	db             intf.DbProvider
	sqlxDb         intf.SqlxDbProvider
	dbBuilder      intf.DbBuilderProvider
	redis          intf.RedisProvider
	mq             intf.MessageQueueProvider
	//graph          intf.GraphProvider
}

var (
	_instance *SdkInstance
	once      sync.Once
)

func New(app, env string, options ...Option) *SdkInstance {
	once.Do(
		func() {
			v, err := newInstance(app, env, options...)
			if err != nil {
				logger.Fatal("new sdk instance", "err", err)
			}
			_instance = v
		},
	)
	return _instance
}

func HasInitialized() bool {
	return _instance != nil
}

func GetInstance() *SdkInstance {
	return _instance
}

func (i *SdkInstance) LoadConfig(configVar any) *SdkInstance {
	if i.configProvider == nil {
		logger.Fatal("config provider not initialized")
	}

	err := i.configProvider.Unmarshal(configVar)
	if err != nil {
		logger.Fatal("unmarshal to config variable", "err", err)
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
		case intf.ProviderCategoryDbBuilder: // will removed in the future
			fxOptions = append(fxOptions, c.Module, fx.Populate(&_instance.dbBuilder))
		case intf.ProviderCategoryRedis:
			fxOptions = append(fxOptions, c.Module, fx.Populate(&_instance.redis))
		//case intf.ProviderCategoryGraph:
		//	fxOptions = append(fxOptions, c.Module, fx.Populate(&_instance.graph))
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
	if !i.option.debug {
		fxOptions = append(fxOptions, fx.NopLogger)
	}

	err := fx.New(fxOptions...).Start(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func newInstance(app, env string, options ...Option) (*SdkInstance, error) {
	sdkOption := defaultSdkOption
	for _, apply := range options {
		apply(sdkOption)
	}

	var viperOptions []viper.Option
	if sdkOption.configFilePath != "" {
		viperOptions = append(viperOptions, viper.WithConfigFile(sdkOption.configFilePath))
	}

	configProvider, err := viper.New(app, env, viperOptions...)
	if err != nil {
		return nil, err
	}

	return &SdkInstance{
		option:         sdkOption,
		configProvider: configProvider,
	}, nil
}
