package intf

import (
	"context"
	"github.com/hdget/hdsdk/v2/provider/mq"
)

type MsgPayload []byte

type MqProvider interface {
	Provider
	NewPublisher() (Publisher, error)
	NewSubscriber() (Subscriber, error)
}

// Publisher is the emitting part of a Pub/Sub.
type Publisher interface {
	// Publish publishes provided messages to given topic.
	//
	// Publish can be synchronous or asynchronous - it depends on the implementation.
	//
	// Most publishers implementations don't support atomic publishing of messages.
	// This means that if publishing one of the messages fails, the next messages will not be published.
	//
	// Publish must be thread safe.
	Publish(topic string, messages []MsgPayload) error
	// Close should flush unsent messages, if publisher is async.
	Close() error
	PublishDelay(topic string, messages []MsgPayload, delaySeconds int64) error
}

// Subscriber is the consuming part of the Pub/Sub.
type Subscriber interface {
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
	SubscribeDelay(ctx context.Context, topic string) (<-chan *mq.Message, error)
	// Close closes all subscriptions with their output channels and flush offsets etc. when needed.
	Close() error
}
