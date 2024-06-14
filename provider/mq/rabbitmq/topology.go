package rabbitmq

import (
	"fmt"
	"github.com/hdget/hdutils/text"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

// direct: exchange --> RoutingKey --> q1,q2...
// topic: exchange --> pattern match --> q1,q2...
// fanout: exchange --> routingkey not used --> q1, q2
// publisher --- RoutingKey ---> exchange
// subscriber -- queue --- BindingKey --- exchange
// RoutingKey: publisher send message to which exchange depends on RoutingKey
// BindingKey: exchange bind to which queue depends on BindingKey

type TopologyKind int

const (
	TopologyKindDefault TopologyKind = iota // default topology
	TopologyKindDelay                       // delay topology kind
)

type ExchangeKind string

const (
	ExchangeKindFanout ExchangeKind = "fanout"
)

type Topology struct {
	Kind         TopologyKind
	ExchangeKind ExchangeKind
	ExchangeName string
	QueueName    string
	RoutingKey   string
	BindingKey   string
}

// 相同name的多个订阅者如果订阅同一个topic,则只有一个订阅者会收到消息
// 不同name的多个订阅者果订阅同一个topic,则所有订阅者都会收到消息
func newTopology(name, topic string, useDelayTopology bool) (*Topology, error) {
	cleanName := text.CleanString(name)
	if cleanName == "" {
		return nil, fmt.Errorf("invalid name, name: %s", name)
	}

	cleanTopic := text.CleanString(topic)
	if cleanTopic == "" {
		return nil, fmt.Errorf("invalid topic, topic: %s", topic)
	}

	// use explicit exchange
	kind := TopologyKindDefault
	if useDelayTopology {
		kind = TopologyKindDelay
	}
	return &Topology{
		Kind:         kind,
		ExchangeKind: ExchangeKindFanout,
		ExchangeName: cleanTopic,
		QueueName:    fmt.Sprintf("%s@%s", cleanTopic, cleanName),
	}, nil
}

func (t *Topology) DeclareExchange(amqpChannel *amqp.Channel) error {
	var err error
	if t.Kind == TopologyKindDelay {
		err = amqpChannel.ExchangeDeclare(
			t.ExchangeName,
			"x-delayed-message",
			true,
			false,
			false,
			false,
			amqp.Table{"x-delayed-type": string(t.ExchangeKind)},
		)
	} else {
		err = amqpChannel.ExchangeDeclare(
			t.ExchangeName,
			string(t.ExchangeKind),
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
		t.QueueName,
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
