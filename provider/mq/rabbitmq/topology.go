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

type ExchangeKind string

const (
	ExchangeKindDefault ExchangeKind = ""               // use default exchange
	ExchangeKindDelay   ExchangeKind = "x-delayed-type" // delay exchange
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
	exchangeKind := ExchangeKindDefault
	index := strings.Index(strings.ToUpper(topic), "@DELAY")
	if index > -1 {
		exchangeKind = ExchangeKindDelay
	}

	result := &Topology{
		exchangeKind: exchangeKind,
		exchangeType: ExchangeTypeFanout,
	}

	remains := topic
	if index > -1 {
		remains = topic[:index]
	}

	tokens := strings.Split(remains, ":")
	switch len(tokens) {
	case 1: // use default exchange
		if exchangeKind == ExchangeKindDelay {
			return nil, errors.New("default exchange doesn't support delay feature")
		}
		result.exchangeType = ExchangeTypeDirect
		result.queueName = tokens[0]
		result.routingKey = tokens[0]
	case 2: // use explicit exchange
		result.exchangeName = tokens[0]
		result.queueName = fmt.Sprintf("%s_%s", tokens[0], tokens[1])
	}
	return result, nil
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
