package hdsdk

import (
	"github.com/hdget/hdsdk/v1/core/config"
	"github.com/hdget/hdsdk/v1/core/logger"
	"github.com/hdget/hdsdk/v1/intf"
	"github.com/hdget/hdsdk/v1/provider/cache"
	"github.com/hdget/hdsdk/v1/provider/db"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

var (
	configLoader intf.ConfigLoader
	Logger       intf.Logger
	Mysql        intf.DbProvider
	Redis        intf.RedisProvider
)

// LoadConfig 将配置文件中的内容加载到configVar中
func LoadConfig(configVar any) error {
	if configLoader == nil {
		return errors.New("please initialize sdk first")
	}
	return configLoader.Load(&configVar)
}

// Initialize 初始化SDK
func Initialize(app, env string, options ...config.Option) error {
	_ = fx.New(
		//fx.NopLogger,
		config.FxModule,
		logger.FxModule,
		db.FxModule,
		cache.FxModule,
		fx.Provide(func() config.Params {
			return config.Params{
				App:     app,
				Env:     env,
				Options: options,
			}
		}),
		fx.Populate(&configLoader),
		fx.Populate(&Logger),
		fx.Populate(&Mysql),
		fx.Populate(&Redis),
	)

	return nil
}
