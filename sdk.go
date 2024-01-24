package hdsdk

import (
	"github.com/hdget/hdsdk/v1/config"
	"github.com/hdget/hdsdk/v1/intf"
	"github.com/hdget/hdsdk/v1/provider/cache"
	"github.com/hdget/hdsdk/v1/provider/db"
	"github.com/hdget/hdsdk/v1/provider/graph"
	"github.com/hdget/hdsdk/v1/provider/logger"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

var (
	configLoader intf.ConfigLoader
	Logger       intf.LoggerProvider
	Mysql        intf.DbProvider
	Redis        intf.RedisProvider
	Neo4j        intf.GraphProvider
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
	// 初始化configLoader
	configLoader = config.NewConfigLoader(app, env, options...)

	// 加载SDKConfig
	sdkConfiger, err := config.NewSdkConfiger(configLoader)
	if err != nil {
		return err
	}

	// 默认加载LoggerProvider
	fxOptions := []fx.Option{
		fx.Provide(func() intf.SdkConfiger { return sdkConfiger }),
		logger.FxModule,
		fx.Populate(&Logger),
	}

	if len(sdkConfiger.GetMysqlConfig()) > 0 {
		fxOptions = append(fxOptions, db.FxModule, fx.Populate(&Mysql))
	}

	if len(sdkConfiger.GetRedisConfig()) > 0 {
		fxOptions = append(fxOptions, cache.FxModule, fx.Populate(&Redis))
	}

	if len(sdkConfiger.GetNeo4jConfig()) > 0 {
		fxOptions = append(fxOptions, graph.FxModule, fx.Populate(&Neo4j))
	}

	_ = fx.New(
		fxOptions...,
	)

	return nil
}
