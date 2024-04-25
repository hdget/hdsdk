package rabbitmq

import (
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

// direct: exchange --> routingKey --> q1,q2...
// topic: exchange --> pattern match --> q1,q2...
// fanout: exchange --> routingkey not used --> q1, q2

// RoutingKey: publisher send message to which exchange depends on RoutingKey
// BindingKey: exchange bind to which queue depends on BindingKey

type ExchangeKind string

const (
	ExchangeKindDefault  ExchangeKind = "default"        // use default exchange
	ExchangeKindExplicit ExchangeKind = "explicit"       // use customized exchange
	ExchangeKindDelay    ExchangeKind = "x-delayed-type" // delay exchange
)

type ExchangeType string

const (
	ExchangeTypeDirect  ExchangeType = "direct"
	ExchangeTypeTopic   ExchangeType = "topic"
	ExchangeTypeFanout  ExchangeType = "fanout"
	ExchangeTypeHeaders ExchangeType = "headers"
)

type Topology struct {
	exchangeKind ExchangeKind
	exchangeType ExchangeType
	exchangeName string
	queueName    string
	routingKey   string
	bindingKey   string
}

func newTopology(topic string, exchangeKind ExchangeKind, exchangeType ExchangeType) *Topology {
	switch exchangeKind {
	case ExchangeKindDefault: // use default exchange, as default exchange implicitly bind to all queues, routingKey==queueName
		return &Topology{
			exchangeName: "",
			queueName:    topic,
			exchangeKind: exchangeKind,
			exchangeType: exchangeType,
			routingKey:   topic, // queue name
			bindingKey:   topic, // queue name
		}
	case ExchangeKindExplicit: // use custom exchange
		return &Topology{
			exchangeName: topic,
			exchangeKind: exchangeKind,
			exchangeType: exchangeType,
			queueName:    topic,
			routingKey:   "", // queue name
			bindingKey:   "",
		}
	case ExchangeKindDelay:
		return &Topology{
			exchangeName: topic,
			exchangeKind: exchangeKind,
			exchangeType: exchangeType,
			queueName:    topic,
			routingKey:   topic, // queue name
			bindingKey:   "",
		}

	}
	return nil
}

func (t *Topology) declareExchange(amqpChannel *amqp.Channel) error {
	var exchangeArgs amqp.Table
	if t.exchangeKind == ExchangeKindDelay {
		exchangeArgs = amqp.Table{string(t.exchangeKind): string(t.exchangeType)}
	}

	return amqpChannel.ExchangeDeclare(
		t.exchangeName,
		string(t.exchangeType),
		true,
		false,
		false,
		false,
		exchangeArgs,
	)
}

func (t *Topology) declareQueue(amqpChannel *amqp.Channel) error {
	_, err := amqpChannel.QueueDeclare(
		t.queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "cannot declare queue")
	}
	return nil
}

func (t *Topology) bindQueue(amqpChannel *amqp.Channel) error {
	err := amqpChannel.QueueBind(
		t.queueName,
		t.bindingKey,
		t.exchangeName,
		false,
		nil)
	if err != nil {
		return errors.Wrap(err, "cannot bind queue")
	}
	return nil
}
