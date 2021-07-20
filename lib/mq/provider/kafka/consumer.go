package kafka

import (
	"context"
	"github.com/hdget/sdk/types"
	"strings"
)

type KafkaConsumer struct {
	Logger types.LogProvider
	Client *ConsumerClient

	// 如果是consumer group中的handler
	handler *ConsumerGroupHandler
	ctx     context.Context
	cancel  context.CancelFunc
}

var _ types.MqConsumer = (*KafkaConsumer)(nil)

func (k *Kafka) CreateConsumer(name string, processFunc types.MqMsgProcessFunc, args ...map[types.MqOptionType]types.MqOptioner) (types.MqConsumer, error) {
	options := k.GetDefaultOptions()
	if len(args) > 0 {
		options = args[0]
	}

	// 初始化kafka client
	client, err := k.newConsumerClient(name, options)
	if err != nil {
		return nil, err
	}

	// 连接
	err = client.connect(k.Config.Brokers)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &KafkaConsumer{
		Logger: k.Logger,
		Client: client,
		handler: &ConsumerGroupHandler{
			Logger:  k.Logger,
			Process: processFunc,
			ready:   make(chan bool),
		},
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

// Consume 消费消息
func (kc *KafkaConsumer) Consume() {
	go func() {
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := kc.Client.saramaConsumerGroup.Consume(kc.ctx, strings.Split(kc.Client.Config.Topic, ","), kc.handler); err != nil {
				kc.Logger.Error("consume in group", "err", err)
			}

			// check if context was cancelled, signaling that the consumer should stop
			if kc.ctx.Err() != nil {
				return
			}

			// 如果relanance发生，ready需要重新创建
			kc.handler.ready = make(chan bool)
		}
	}()

	<-kc.handler.ready // Await till the consumer handler has been set up
	kc.Logger.Debug("sarama handler in consumer group up and running!...")

	<-kc.ctx.Done()
}

func (kc *KafkaConsumer) Close() {
	kc.cancel()
	kc.Client.close()
}
