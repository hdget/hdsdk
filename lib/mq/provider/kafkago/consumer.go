package kafkago

import (
	"github.com/hdget/sdk/types"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	Logger  types.LogProvider
	Client  *ConsumerClient
	Process types.MqMsgProcessFunc
	Buffers []kafka.Message
}

var _ types.MqConsumer = (*KafkaConsumer)(nil)

func (k *Kafka) CreateConsumer(name string, processFunc types.MqMsgProcessFunc, args ...map[types.MqOptionType]types.MqOptioner) (types.MqConsumer, error) {
	options := k.GetDefaultOptions()
	if len(args) > 0 {
		options = args[0]
	}

	// 初始化kafka client
	client, err := k.NewConsumerClient(name, options)
	if err != nil {
		return nil, err
	}

	// 连接
	err = client.connect(k.Config.Brokers)
	if err != nil {
		return nil, err
	}

	c := &KafkaConsumer{
		Logger:  k.Logger,
		Client:  client,
		Process: processFunc,
		Buffers: make([]kafka.Message, 0),
	}

	return c, nil
}

// Consume 消费消息
func (kc *KafkaConsumer) Consume() {
	countRetry := 0
	for {
		// reader会自动重连
		msg, err := kc.Client.Reader.FetchMessage(kc.Client.ctx)
		if err != nil {
			kc.Logger.Error("kafkago fetch message", "name", kc.Client.Name, "retry", countRetry, "err", err)
			break
		}

		ret := kc.Process(msg.Value)
		switch ret {
		case types.Ack:
			err := kc.Client.Reader.CommitMessages(kc.Client.ctx, msg)
			if err != nil {
				kc.Logger.Error("kafkago commit message", "name", kc.Client.Name, "retry", countRetry, "err", err)
			}
		case types.Next:
			kc.Buffers = append(kc.Buffers, msg)
		case types.BatchAck:
			err := kc.Client.Reader.CommitMessages(kc.Client.ctx, kc.Buffers...)
			if err != nil {
				kc.Logger.Error("kafkago batch commit messages", "name", kc.Client.Name, "retry", countRetry, "err", err)
			}
		default:
			// do nothing
		}
	}
}

func (kc *KafkaConsumer) Close() {
	kc.Client.Reader.Close()
}
