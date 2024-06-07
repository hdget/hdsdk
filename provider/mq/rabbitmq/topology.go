package rabbitmq

import (
	"fmt"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"strings"
)

// direct: exchange --> routingKey --> q1,q2...
// topic: exchange --> pattern match --> q1,q2...
// fanout: exchange --> routingkey not used --> q1, q2
// publisher --- routingKey ---> exchange
// subscriber -- queue --- bindingKey --- exchange
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
)

type Topology struct {
	exchangeKind ExchangeKind
	exchangeType ExchangeType
	exchangeName string
	queueName    string
	routingKey   string
	bindingKey   string
}

// newTopology will parse topic string to decide the exchange, queue, routing key ...
// <exchange>:<routingKey>@<exchange type>
// order:close ===> exchange: order, exchangeType: fanout, queue: order_close, routingKey: "", bindingKey: "", exchangeKind: default
// cancel@delay ===> exchange: order, routingKey: cancel, exchangeKind: delay
func newTopology(topic string) (*Topology, error) {
	var result *Topology
	index := strings.Index(topic, ":")
	switch index {
	case -1:
		// use default exchange
		result = &Topology{
			exchangeKind: ExchangeKindDefault,
			exchangeType: ExchangeTypeDirect,
			queueName:    topic,
			routingKey:   topic,
		}
	default:
		// use explicit exchange
		result = &Topology{
			exchangeKind: ExchangeKindDefault,
			exchangeType: ExchangeTypeFanout,
			exchangeName: topic[index:],
			queueName:    fmt.Sprintf("%s_%s", topic[:index], topic[index:]),
		}
	}
	return result, nil
}

func newDelayTopology(topic string) (*Topology, error) {
	var result *Topology
	index := strings.Index(topic, ":")
	switch index {
	case -1:
		// use default exchange
		result = &Topology{
			exchangeKind: ExchangeKindDelay,
			exchangeType: ExchangeTypeDirect,
			queueName:    topic,
			routingKey:   topic,
		}
	default:
		// use explicit exchange
		result = &Topology{
			exchangeKind: ExchangeKindDelay,
			exchangeType: ExchangeTypeFanout,
			exchangeName: topic[index:],
			queueName:    fmt.Sprintf("%s_%s", topic[:index], topic[index:]),
		}
	}
	return result, nil
}

func (t *Topology) declareExchange(amqpChannel *amqp.Channel) error {
	var err error
	if t.exchangeKind == ExchangeKindDelay {
		err = amqpChannel.ExchangeDeclare(
			t.exchangeName,
			"x-delayed-message",
			true,
			false,
			false,
			false,
			amqp.Table{"x-delayed-type": string(t.exchangeType)},
		)
	} else {
		err = amqpChannel.ExchangeDeclare(
			t.exchangeName,
			string(t.exchangeType),
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

func (t *Topology) declareQueue(amqpChannel *amqp.Channel) error {
	_, err := amqpChannel.QueueDeclare(
		t.queueName, // queue: exchangeName_key
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
		t.queueName, // queue: exchangeName_key
		t.bindingKey,
		t.exchangeName,
		false,
		nil)
	if err != nil {
		return errors.Wrap(err, "cannot bind queue")
	}

	return nil
}
