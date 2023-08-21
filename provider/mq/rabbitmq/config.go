package rabbitmq

import (
	"github.com/hdget/hdsdk/types"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type MqProviderConfig struct {
	Default *MqConfig   `mapstructure:"default"` // 缺省
	Items   []*MqConfig `mapstructure:"items"`
}

// ConsumerConfig 客户端配置
//type ConsumerConfig struct {
//	Name         string   `mapstructure:"name"`
//	ExchangeName string   `mapstructure:"exchange_name"`
//	ExchangeType string   `mapstructure:"exchange_type"`
//	QueueName    string   `mapstructure:"queue_name"`
//	RoutingKeys  []string `mapstructure:"routing_keys"`
//}

// ProducerConfig 发送端配置
//type ProducerConfig struct {
//	//Name         string `mapstructure:"name"`
//	//ExchangeName string `mapstructure:"exchange_name"`
//	//ExchangeType string `mapstructure:"exchange_type"`
//}

// MqConfig amqp://user:pass@host:10000/vhost
type MqConfig struct {
	Name     string `mapstructure:"name"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Vhost    string `mapstructure:"vhost"`
	//Consumers []*ConsumerConfig `mapstructure:"consumers"`
	//Producers []*ProducerConfig `mapstructure:"producers"`
}

var (
	errInvalidProducerParam = errors.New("invalid producer params")
)

// 校验Mq配置
func validateMqConfig(providerType string, conf *MqConfig) error {
	if conf == nil {
		return types.ErrEmptyConfig
	}

	if conf.Host == "" || conf.Username == "" || conf.Password == "" || conf.Port == 0 {
		return types.ErrInvalidConfig
	}

	// extra provider需要提供name
	if providerType == types.ProviderTypeOther && conf.Name == "" {
		return types.ErrInvalidConfig
	}

	return nil
}

// 解析Mq配置
func parseConfig(rootConfiger types.Configer) (*MqProviderConfig, error) {
	if rootConfiger == nil {
		return nil, types.ErrEmptyConfig
	}

	data := rootConfiger.GetRabbitmqConfig()
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

//func (rmq *RabbitMq) getConsumerConfig(name string) (*ConsumerConfig, error) {
//	var found *ConsumerConfig
//	for _, conf := range rmq.Config.Consumers {
//		if conf.Name == name {
//			found = conf
//			break
//		}
//	}
//
//	if found == nil {
//		return nil, ErrConsumerConfigNotFound
//	}
//
//	if found.QueueName == "" ||
//		found.Name == "" ||
//		!utils.StringSliceContains(SupportedExchangeTypes, found.ExchangeType) {
//		return nil, ErrInvalidConsumerConfig
//	}
//
//	return found, nil
//}

//func (rmq *RabbitMq) getProducerConfig(name string) (*ProducerConfig, error) {
//	var found *ProducerConfig
//	for _, conf := range rmq.Config.Producers {
//		if conf.Name == name {
//			found = conf
//			break
//		}
//	}
//
//	if found == nil {
//		return nil, ErrProducerConfigNotFound
//	}
//
//	return found, nil
//}
