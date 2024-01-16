package config

import "github.com/hdget/hdsdk/intf"

// 首先必须创建一个继承自hdsdk.Config的配置struct
// e,g:
//
// import "github.com/hdget/hdsdk"
//
// type XXXConfig struct {
//		*hdsdk.Config `mapstructure:",squash"`
// }

// ConfigItem items
type sdkConfig struct {
	Log      map[string]any `mapstructure:"log"`      // 日志配置
	Mysql    map[string]any `mapstructure:"mysql"`    // 数据库配置
	Redis    map[string]any `mapstructure:"redis"`    // 缓存配置
	RabbitMq map[string]any `mapstructure:"rabbitmq"` // rabbitmq消息队列配置
	Etcd     map[string]any `mapstructure:"etcd"`     // etcd Key/Value数据库配置
	Neo4j    map[string]any `mapstructure:"neo4j"`    // neo4j数据库配置
}

const (
	sdkConfigSection = "sdk"
)

func NewConfiger(loader intf.ConfigLoader) (intf.Configer, error) {
	var config sdkConfig
	err := loader.LoadKey(sdkConfigSection, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// GetLogConfig 获取日志配置
func (c *sdkConfig) GetLogConfig() map[string]any {
	return c.Log
}

// GetMysqlConfig 获取数据库配置
func (c *sdkConfig) GetMysqlConfig() map[string]any {
	return c.Mysql
}

// GetRedisConfig 获取缓存配置
func (c *sdkConfig) GetRedisConfig() map[string]any {
	return c.Redis
}

// GetRabbitmqConfig 获取消息队列配置
func (c *sdkConfig) GetRabbitmqConfig() map[string]any {
	return c.RabbitMq
}

// GetEtcdConfig 获取etcd配置
func (c *sdkConfig) GetEtcdConfig() map[string]any {
	return c.Etcd
}

// GetNeo4jConfig 获取图数据库配置
func (c *sdkConfig) GetNeo4jConfig() map[string]any {
	return c.Neo4j
}
