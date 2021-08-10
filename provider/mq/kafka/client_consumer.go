package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/hdget/hdsdk/types"
)

type ConsumerClient struct {
	*BaseClient
	Option *ConsumeOption
	Config *ConsumerConfig

	saramaConfig        *sarama.Config
	saramaClient        sarama.Client
	saramaConsumerGroup sarama.ConsumerGroup
}

func (k *Kafka) newConsumerClient(name string, options map[types.MqOptionType]types.MqOptioner) (*ConsumerClient, error) {
	// 获取匹配的路由配置
	config := k.getConsumerConfig(name)
	if config == nil {
		return nil, fmt.Errorf("no matched consumer config for: %s", name)
	}

	cc := &ConsumerClient{
		BaseClient: k.newBaseClient(name, options),
		Config:     config,
		Option:     getConsumeOption(options),
	}

	cc.saramaConfig = cc.getSaramaConfig()
	return cc, nil
}

// connect 连接kafka, 生成consumerGroup
func (cc *ConsumerClient) connect(brokers []string) error {
	saramaClient, err := sarama.NewClient(brokers, cc.saramaConfig)
	if err != nil {
		return err
	}
	cc.saramaClient = saramaClient

	consumerGroup, err := sarama.NewConsumerGroupFromClient(cc.Config.GroupId, saramaClient)
	if err != nil {
		return err
	}
	cc.saramaConsumerGroup = consumerGroup

	return nil
}

// 获取sarama配置
func (cc *ConsumerClient) getSaramaConfig() *sarama.Config {
	saramaConfig := sarama.NewConfig()

	// 固定版本号
	saramaConfig.Version = sarama.V0_11_0_2

	// 检查是否需要创建安全SASL认证
	if cc.Config.User != "" && cc.Config.Password != "" {
		// SASL
		saramaConfig.Net.SASL.Enable = true
		saramaConfig.ClientID = cc.Config.Name
		// user = <消费组的账号>-<消费组ID>
		saramaConfig.Net.SASL.User = fmt.Sprintf("%s-%s", cc.Config.User, cc.Config.GroupId)
		saramaConfig.Net.SASL.Password = cc.Config.Password
	}

	// consume options
	saramaConfig.Consumer.Return.Errors = cc.Option.ReturnErrors
	saramaConfig.Consumer.Offsets.Initial = cc.Option.InitialOffset
	saramaConfig.Consumer.Offsets.AutoCommit.Enable = cc.Option.AutoCommit

	return saramaConfig
}

func (cc *ConsumerClient) close() {
	if cc.saramaConsumerGroup != nil {
		cc.saramaConsumerGroup.Close()
	}

	if cc.saramaClient != nil {
		cc.saramaClient.Close()
	}
}
