package sdk

import (
	"github.com/hdget/sdk/types"
)

type Config struct {
	Sdk *types.SdkConfigItem `mapstructure:"sdk"`
}

var _ types.Configer = (*Config)(nil)

// GetMysqlConfig 获取数据库配置
func (c *Config) GetMysqlConfig() interface{} {
	if c == nil || c.Sdk == nil {
		return nil
	}

	return c.Sdk.Mysql
}

// GetRedisConfig 获取缓存配置
func (c *Config) GetRedisConfig() interface{} {
	if c == nil || c.Sdk == nil {
		return nil
	}

	return c.Sdk.Redis
}

// GetLogConfig 获取日志配置
func (c *Config) GetLogConfig() interface{} {
	if c == nil || c.Sdk == nil {
		return nil
	}

	return c.Sdk.Log
}

// GetRabbitmqConfig 获取消息队列配置
func (c *Config) GetRabbitmqConfig() interface{} {
	if c == nil || c.Sdk == nil {
		return nil
	}

	return c.Sdk.RabbitMq
}

// GetKafkaConfig 获取Kafka消息队列配置
func (c *Config) GetKafkaConfig() interface{} {
	if c == nil || c.Sdk == nil {
		return nil
	}

	return c.Sdk.Kafka
}

// GetNosqlConfig 获取非SQL配置
func (c *Config) GetNosqlConfig() interface{} {
	if c == nil || c.Sdk == nil {
		return nil
	}

	return c.Sdk.Kv
}

// GetKvConfig 获取KV配置
func (c *Config) GetKvConfig() interface{} {
	if c == nil || c.Sdk == nil {
		return nil
	}
	return c.Sdk.Kv
}
