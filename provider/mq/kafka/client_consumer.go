package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/hdget/hdsdk/types"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type ConsumerClient struct {
	*BaseClient
	Option              *ConsumeOption
	Parameter           *ConsumerParameter
	saramaConfig        *sarama.Config
	saramaClient        sarama.Client
	saramaConsumerGroup sarama.ConsumerGroup
}

type ConsumerParameter struct {
	Name     string `mapstructure:"name"`
	GroupId  string `mapstructure:"groupId"`
	Topic    string `mapstructure:"topic"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

func (k *Kafka) newConsumerClient(parameters map[string]interface{}, args ...types.MqOptioner) (*ConsumerClient, error) {
	consumerParam, err := parseConsumerParameter(parameters)
	if err != nil {
		return nil, err
	}

	client := k.newBaseClient(args...)
	cc := &ConsumerClient{
		BaseClient: client,
		Option:     getConsumeOption(client.options),
		Parameter:  consumerParam,
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

	consumerGroup, err := sarama.NewConsumerGroupFromClient(cc.Parameter.GroupId, saramaClient)
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
	if cc.Parameter.User != "" && cc.Parameter.Password != "" {
		// SASL
		saramaConfig.Net.SASL.Enable = true
		saramaConfig.ClientID = cc.Parameter.Name
		// user = <消费组的账号>-<消费组ID>
		saramaConfig.Net.SASL.User = fmt.Sprintf("%s-%s", cc.Parameter.User, cc.Parameter.GroupId)
		saramaConfig.Net.SASL.Password = cc.Parameter.Password
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

func parseConsumerParameter(params map[string]interface{}) (*ConsumerParameter, error) {
	var consumerParams ConsumerParameter
	err := mapstructure.Decode(params, &consumerParams)
	if err != nil {
		return nil, err
	}

	if consumerParams.GroupId == "" ||
		consumerParams.Topic == "" ||
		consumerParams.Name == "" ||
		consumerParams.User == "" {
		return nil, errors.New("invalid parameter")
	}

	return &consumerParams, nil
}
