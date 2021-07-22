package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/hdget/sdk/types"
)

type ProducerClient struct {
	*BaseClient
	Config *ProducerConfig
	Option *PublishOption

	saramaClient sarama.Client
	saramaConfig *sarama.Config
	handler      sarama.SyncProducer
}

func (k *Kafka) newProducerClient(name string, options map[types.MqOptionType]types.MqOptioner) (*ProducerClient, error) {
	// 获取匹配的路由配置
	config := k.getProducerConfig(name)
	if config == nil {
		return nil, fmt.Errorf("no matched producer config for: %s", name)
	}

	pc := &ProducerClient{
		BaseClient: k.newBaseClient(name, options),
		Config:     config,
		Option:     getPublishOption(options),
	}

	pc.saramaConfig = pc.getSaramaConfig()
	return pc, nil
}

// connect balance策略
func (pc *ProducerClient) connect(brokers []string) error {
	saramaClient, err := sarama.NewClient(brokers, pc.saramaConfig)
	if err != nil {
		return err
	}
	pc.saramaClient = saramaClient

	handler, err := sarama.NewSyncProducerFromClient(saramaClient)
	if err != nil {
		pc.logger.Error("new producer", "err", err)
	}
	pc.handler = handler
	return nil
}

// 获取sarama配置
func (pc *ProducerClient) getSaramaConfig() *sarama.Config {
	saramaConfig := sarama.NewConfig()

	// 固定版本号
	saramaConfig.Version = sarama.V0_11_0_2

	// consume options
	saramaConfig.Producer.RequiredAcks = pc.Option.RequiredAcks // Wait for local commits success to ack the message
	saramaConfig.Producer.Retry.Max = pc.Option.RetryMax        // Retry up to 10 times to produce the message
	saramaConfig.Producer.Return.Successes = pc.Option.ReturnSuccess
	saramaConfig.Producer.Return.Errors = pc.Option.ReturnError

	return saramaConfig
}

func (pc *ProducerClient) close() {
	if pc.handler != nil {
		pc.handler.Close()
	}

	if pc.saramaClient != nil {
		pc.saramaClient.Close()
	}
}
