package intf

import (
	"context"
	"github.com/hdget/hdsdk/v2/provider/mq"
)

// MessageQueueProvider
// 相同name的多个订阅者如果订阅同一个topic,则只有一个订阅者会收到消息
// 不同name的多个订阅者果订阅同一个topic,则所有订阅者都会收到消息
type MessageQueueProvider interface {
	Provider
	// Publisher 如果要使用PublishDelay接口的时候name必须设置
	Publisher(name string, args ...*mq.PublisherOption) (MessageQueuePublisher, error)
	// Subscriber 如果要使用SubscribeDelay接口的时候name必须设置
	Subscriber(name string, args ...*mq.SubscriberOption) (MessageQueueSubscriber, error)
}

type MessageQueuePublisher interface {
	// Publish publishes provided messages to given topic.
	//
	// Publish can be synchronous or asynchronous - it depends on the implementation.
	//
	// Most publishers implementations don't support atomic publishing of messages.
	// This means that if publishing one of the messages fails, the next messages will not be published.
	//
	// Publish must be thread safe.
	Publish(topic string, messages [][]byte, delaySeconds ...int64) error
	// Close should flush unsent messages, if publisher is async.
	Close() error
}

type MessageQueueSubscriber interface {
	// Subscribe returns output channel with messages from provided topic.
	// Channel is closed, when Close() was called on the subscriber.
	//
	// To receive the next message, `Ack()` must be called on the received message.
	// If message processing failed and message should be redelivered `Nack()` should be called.
	//
	// When provided ctx is cancelled, subscriber will close subscribe and close output channel.
	// Provided ctx is set to all produced messages.
	// When Nack or Ack is called on the message, context of the message is canceled.
	Subscribe(ctx context.Context, topic string) (<-chan *mq.Message, error)
	// Close closes all subscriptions with their output channels and flush offsets etc. when needed.
	Close() error
}
