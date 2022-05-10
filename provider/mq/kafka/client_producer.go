package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/hdget/hdsdk/types"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type ProducerClient struct {
	*BaseClient
	Parameter *ProducerParameter
	Option    *PublishOption

	saramaClient sarama.Client
	saramaConfig *sarama.Config
	handler      sarama.SyncProducer
}

type ProducerParameter struct {
	Topics []string `mapstructure:"topics"`
}

func (k *Kafka) newProducerClient(parameters map[string]interface{}, options ...types.MqOptioner) (*ProducerClient, error) {
	producerParams, err := parseProducerParameter(parameters)
	if err != nil {
		return nil, err
	}

	client := k.newBaseClient(options...)

	pc := &ProducerClient{
		BaseClient: client,
		Parameter:  producerParams,
		Option:     GetPublishOption(client.options),
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

func parseProducerParameter(params map[string]interface{}) (*ProducerParameter, error) {
	var producerParams ProducerParameter
	err := mapstructure.Decode(params, &producerParams)
	if err != nil {
		return nil, err
	}

	if len(producerParams.Topics) == 0 {
		return nil, errors.New("invalid parameter")
	}

	return &producerParams, nil
}
