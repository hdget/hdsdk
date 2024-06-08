package rabbitmq

import (
	"fmt"
	"github.com/hdget/hdutils/text"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"strings"
)

// direct: exchange --> RoutingKey --> q1,q2...
// topic: exchange --> pattern match --> q1,q2...
// fanout: exchange --> routingkey not used --> q1, q2
// publisher --- RoutingKey ---> exchange
// subscriber -- queue --- BindingKey --- exchange
// RoutingKey: publisher send message to which exchange depends on RoutingKey
// BindingKey: exchange bind to which queue depends on BindingKey

type ExchangeKind int

const (
	ExchangeKindDefault ExchangeKind = iota // default exchange
	ExchangeKindDelay                       // delay exchange
)

type ExchangeType string

const (
	ExchangeTypeDirect ExchangeType = "direct"
	ExchangeTypeFanout ExchangeType = "fanout"
	exchangeSeparator               = "@"
)

type Topology struct {
	ExchangeKind ExchangeKind
	ExchangeType ExchangeType
	ExchangeName string
	QueueName    string
	RoutingKey   string
	BindingKey   string
}

// newTopology will parse topic string to decide the exchange, queue, routing key ...
// <routingKey>@<exchange>
func newTopology(topic string) (*Topology, error) {
	var result *Topology
	index := strings.LastIndex(topic, exchangeSeparator)
	switch index {
	case -1:
		cleanTopic := text.CleanString(topic)
		if cleanTopic == "" {
			return nil, fmt.Errorf("invalid topic, topic: %s", topic)
		}

		// use default exchange
		result = &Topology{
			ExchangeKind: ExchangeKindDefault,
			ExchangeType: ExchangeTypeDirect,
			QueueName:    cleanTopic,
			RoutingKey:   cleanTopic,
		}
	default:
		cleanExchange := text.CleanString(topic[index+1:])
		if cleanExchange == "" {
			return nil, fmt.Errorf("invalid exchange, exchange: %s", topic[index+1:])
		}

		cleanTopic := text.CleanString(topic[:index])
		if cleanTopic == "" {
			return nil, fmt.Errorf("invalid topic, topic: %s", topic[index+1:])
		}

		// use explicit exchange
		result = &Topology{
			ExchangeKind: ExchangeKindDefault,
			ExchangeType: ExchangeTypeFanout,
			ExchangeName: cleanExchange,
			QueueName:    fmt.Sprintf("%s:%s", cleanExchange, cleanTopic),
		}
	}
	return result, nil
}

func newDelayTopology(exchangeName, topic string) (*Topology, error) {
	cleanExchangeName := text.CleanString(exchangeName)
	if cleanExchangeName == "" {
		return nil, fmt.Errorf("invalid exchange, exchange: %s", exchangeName)
	}

	cleanTopic := text.CleanString(topic)
	if cleanTopic == "" {
		return nil, fmt.Errorf("invalid topic, topic: %s", topic)
	}

	// use explicit exchange
	return &Topology{
		ExchangeKind: ExchangeKindDelay,
		ExchangeType: ExchangeTypeFanout,
		ExchangeName: fmt.Sprintf("delay:%s", cleanExchangeName),
		QueueName:    fmt.Sprintf("delay:%s:%s", cleanExchangeName, cleanTopic),
	}, nil
}

func (t *Topology) DeclareExchange(amqpChannel *amqp.Channel) error {
	var err error
	if t.ExchangeKind == ExchangeKindDelay {
		err = amqpChannel.ExchangeDeclare(
			t.ExchangeName,
			"x-delayed-message",
			true,
			false,
			false,
			false,
			amqp.Table{"x-delayed-type": string(t.ExchangeType)},
		)
	} else {
		err = amqpChannel.ExchangeDeclare(
			t.ExchangeName,
			string(t.ExchangeType),
			true,
			false,
			false,
			false,
			nil,
		)
	}
	if err != nil {
		return err
	}

	return nil
}

func (t *Topology) DeclareQueue(amqpChannel *amqp.Channel) error {
	_, err := amqpChannel.QueueDeclare(
		t.QueueName, // queue: exchangeName_key
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

func (t *Topology) BindQueue(amqpChannel *amqp.Channel) error {
	err := amqpChannel.QueueBind(
		t.QueueName, // queue: exchangeName_key
		t.BindingKey,
		t.ExchangeName,
		false,
		nil)
	if err != nil {
		return errors.Wrap(err, "cannot bind queue")
	}

	return nil
}
