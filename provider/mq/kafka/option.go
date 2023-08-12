package kafka

import (
	"github.com/IBM/sarama"
	"github.com/hdget/hdsdk/types"
)

type PublishOption struct {
	// NoResponse:0,  doesn't send any response, the TCP ACK is all you get
	// WaitForLocal:1,  waits for only the local commit to succeed before responding
	// WaitForAll:-1,  waits for all in-sync replicas to commit before responding
	RequiredAcks sarama.RequiredAcks
	// The total number of times to retry sending a message (default 3)
	RetryMax int
	// If enabled, successfully delivered messages will be returned on the successes channel
	ReturnSuccess bool
	// If enabled, successfully delivered messages will be returned on the successes channel
	ReturnError bool
}

type ConsumeOption struct {
	InitialOffset int64 // 如果之前没有offset提交，选择offset的策略
	// broker等待Consumer.Fetch.Min的最长时间，
	// 如果没取到足够Consumer.Fetch.Min, 等待MaxWaitTime后也会返回
	// MaxWaitTime     time.Duration
	// 是否在消费时有任何错误都会返回到Errors通道
	ReturnErrors bool
	// 是否自动提交
	AutoCommit bool
}

var (
	defaultPublishOption = &PublishOption{
		RequiredAcks:  sarama.WaitForLocal, // 需要本地commit成功
		RetryMax:      3,                   // 发送失败3次开始重试
		ReturnSuccess: true,                // 默认需要确认成功发送
		ReturnError:   true,
	}
	defaultConsumeOption = &ConsumeOption{
		// OffsetNewest:-1, 下一条生产到partition的消息的offset
		// OffsetOldest:-2, Broker的partition上最老的offset
		InitialOffset: -1,
		ReturnErrors:  true,
		// 默认禁止自动提交
		AutoCommit: false,
	}
)

func (q PublishOption) GetType() types.MqOptionType {
	return types.MqOptionPublish
}

func (q ConsumeOption) GetType() types.MqOptionType {
	return types.MqOptionConsume
}

func GetPublishOption(options map[types.MqOptionType]types.MqOptioner) *PublishOption {
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
