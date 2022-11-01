package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/hdget/hdsdk/types"
)

type ConsumerGroupHandler struct {
	Logger  types.LogProvider
	ready   chan bool
	Process types.MqMsgProcessFunc
}

// Setup 消费组在执行Consume时会初始化consumer,并依次调用Setup()->ConsumeClaim(), 关闭之前会调用Cleanup()
// 消息的具体消费发生在ConsumeClaim()中，你可以在Setup()中初始化一些东西
func (cgh ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(cgh.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (cgh ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 在ConsumeClaim()中必须执行不断取消息的循环:ConsumerGroupClaim.Messages()
func (cgh ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for msg := range claim.Messages() {
		action := cgh.Process(msg.Value)
		switch action {
		case types.Ack:
			session.MarkMessage(msg, "")
		default:
			// do nothing
		}
	}

	return nil
}
