package kafka

import (
	"github.com/hdget/sdk/types"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type MqProviderConfig struct {
	Default *MqConfig   `mapstructure:"default"`
	Items   []*MqConfig `mapstructure:"items"`
}

// ConsumerConfig 客户端配置
type ConsumerConfig struct {
	Name      string `mapstructure:"name"`
	Topic     string `mapstructure:"topic"`
	Partition int    `mapstructure:"partition"`
	GroupId   string `mapstructure:"group_id"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
}

// ProducerConfig 发送端配置
type ProducerConfig struct {
	Name        string   `mapstructure:"name"`
	Topics      []string `mapstructure:"topics"`
	Balance     string   `mapstructure:"balance"`
	Compression string   `mapstructure:"compression"`
}

// MqConfig amqp://user:pass@host:10000/vhost
type MqConfig struct {
	Name      string            `mapstructure:"name"`
	Brokers   []string          `mapstructure:"brokers"`
	Consumers []*ConsumerConfig `mapstructure:"consumers"`
	Producers []*ProducerConfig `mapstructure:"producers"`
}

// 校验Mq配置
func validateMqConfig(providerType string, conf *MqConfig) error {
	if conf == nil {
		return types.ErrEmptyConfig
	}

	return nil
}

// 解析Mq配置
func parseConfig(rootConfiger types.Configer) (*MqProviderConfig, error) {
	data := rootConfiger.GetKafkaConfig()
	if data == nil {
		return nil, types.ErrEmptyConfig
	}

	values, ok := data.(map[string]interface{})
	if !ok {
		return nil, types.ErrInvalidConfig
	}

	var conf MqProviderConfig
	err := mapstructure.Decode(values, &conf)
	if err != nil {
		return nil, errors.Wrap(err, "decode mq config")
	}

	return &conf, nil
}

func (k *Kafka) getConsumerConfig(name string) *ConsumerConfig {
	var found *ConsumerConfig
	for _, conf := range k.Config.Consumers {
		if conf.Name == name {
			found = conf
			break
		}
	}
	return found
}

func (k *Kafka) getProducerConfig(name string) *ProducerConfig {
	var found *ProducerConfig
	for _, conf := range k.Config.Producers {
		if conf.Name == name {
			found = conf
			break
		}
	}
	return found
}
