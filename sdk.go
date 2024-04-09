package hdsdk

import (
	"github.com/hdget/hdsdk/v1/intf"
	"github.com/hdget/hdsdk/v1/provider/config/viper"
	"github.com/hdget/hdsdk/v1/provider/logger/zerolog"
	"go.uber.org/fx"
)

//
//// LoadConfig 将配置文件中的内容加载到configVar中
//func LoadConfig(configVar any) error {
//	if configLoader == nil {
//		return errors.New("please initialize sdk first")
//	}
//	return configLoader.LoadLocal(&configVar)
//}
//
//func LoadRemoteConfig(configVar any) error {
//	if configLoader == nil {
//		return errors.New("please initialize sdk first")
//	}
//
//	// 加载SDKConfig
//	sdkConfiger, err := configLoader.GetSDKConfig()
//	if err != nil {
//		return err
//	}
//
//	sdkConfiger.GetEtcdConfig()
//
//	return configLoader.LoadRemote(&configVar)
//}

// Initialize 初始化SDK
func Initialize(app, env string, providers ...fx.Option) error {
	// 所有的provider依赖configProvider和loggerProvider
	fxOptions := []fx.Option{
		fx.Provide(func() (intf.ConfigProvider, error) { return viper.New(app, env) }),
		zerolog.Module,
	}

	fxOptions = append(fxOptions, providers...)
	_ = fx.New(
		fxOptions...,
	)

	return nil
}

// InitializeWithConfig 初始化SDK
func InitializeWithConfig(configProvider intf.ConfigProvider, providers ...fx.Option) error {
	// 所有的provider依赖configProvider和loggerProvider
	fxOptions := []fx.Option{
		fx.Provide(configProvider),
		zerolog.Module,
	}

	fxOptions = append(fxOptions, providers...)
	_ = fx.New(
		fxOptions...,
	)

	return nil
}

//
//// Initialize 初始化SDK
//func Initialize(app, env string, options ...config2.Option) error {
//	// 初始化configLoader
//	configLoader = config2.NewConfigLoader(app, env, options...)
//
//	// 加载SDKConfig
//	sdkConfiger, err := configLoader.GetSDKConfig()
//	if err != nil {
//		return err
//	}
//
//	// 默认加载LoggerProvider
//	fxOptions := []fx.Option{
//		fx.Provide(func() intf.SdkConfiger { return sdkConfiger }),
//		logger.FxModule,
//		fx.Populate(&Logger),
//	}
//
//	if len(sdkConfiger.GetMysqlConfig()) > 0 {
//		fxOptions = append(fxOptions, db.FxModule, fx.Populate(&Mysql))
//	}
//
//	if len(sdkConfiger.GetRedisConfig()) > 0 {
//		fxOptions = append(fxOptions, cache.FxModule, fx.Populate(&Redis))
//	}
//
//	if len(sdkConfiger.GetNeo4jConfig()) > 0 {
//		fxOptions = append(fxOptions, graph.FxModule, fx.Populate(&Neo4j))
//	}
//
//	_ = fx.New(
//		fxOptions...,
//	)
//
//	return nil
//}
