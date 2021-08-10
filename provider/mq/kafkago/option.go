package kafkago

import (
	"github.com/hdget/hdsdk/types"
	"hash"
)

// mandatory=true
// 当没有队列匹配routingKey, 发布的消息也可能处于不能递交状态
// immediate=true
// 如果在匹配的队列上没有消费者准备好，发布的消息也可能处于不能递交状态
type PublishOption struct {
	HashFunc hash.Hash32
}

type ConsumeOption struct {
	MinBytes       int
	MaxBytes       int
	CommitInterval int // flushes commits to Kafka every second
}

var (
	defaultPublishOption = &PublishOption{
		HashFunc: nil,
	}
	defaultConsumeOption = &ConsumeOption{
		MinBytes:       1,
		MaxBytes:       1e6, // default is 1M
		CommitInterval: 1,
	}
)

func (q PublishOption) GetType() types.MqOptionType {
	return types.MqOptionPublish
}

func (q ConsumeOption) GetType() types.MqOptionType {
	return types.MqOptionConsume
}

func getPublishOption(options map[types.MqOptionType]types.MqOptioner) *PublishOption {
	if len(options) == 0 {
		return defaultPublishOption
	}

	v := options[types.MqOptionPublish]
	if v == nil {
		return defaultPublishOption
	}

	option, ok := v.(*PublishOption)
	if !ok {
		return defaultPublishOption
	}
	return option
}

func getConsumeOption(options map[types.MqOptionType]types.MqOptioner) *ConsumeOption {
	if len(options) == 0 {
		return defaultConsumeOption
	}

	v := options[types.MqOptionConsume]
	if v == nil {
		return defaultConsumeOption
	}

	option, ok := v.(*ConsumeOption)
	if !ok {
		return defaultConsumeOption
	}
	return option
}
