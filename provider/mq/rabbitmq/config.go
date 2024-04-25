package rabbitmq

import (
	"github.com/hdget/hdsdk/v2/errdef"
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/pkg/errors"
)

type RabbitMqConfig struct {
	Connection *ConnectionConfig

	// Consumer: Whether or not to requeue when sending a negative acknowledgement in case of a failure.
	RequeueInFailure bool
	// With a prefetch count greater than zero, the server will deliver that many
	// messages to consumers before acknowledgments are received.  The server ignores
	// this option when consumers are started with noAck because no acknowledgments
	// are expected or sent.

	// Consumer: In order to defeat that we can set the prefetch count with the value of 1.
	// This tells RabbitMQ not to give more than one message to a worker at a time.
	// Or, in other words, don't dispatch a new message to a worker until it has
	// processed and acknowledged the previous one.
	// Instead, it will dispatch it to the next worker that is not still busy.
	PrefetchCount int

	// Exchange: Each exchange belongs to one of a set of exchange kinds/types implemented by
	// the server. The exchange types define the functionality of the exchange - i.e.
	// how messages are routed through it. Once an exchange is declared, its type
	// cannot be changed.  The common types are "direct", "fanout", "topic" and
	// "headers".
	ExchangeType ExchangeType

	// Exchange: whether to use explicit exchange other than default exchange
	UseExplicitExchange bool
}

type ConnectionConfig struct {
	Host     string // required, RabbitMQ host
	Port     int    // required, RabbitMQ port
	Username string // required, RabbitMQ username
	Password string // required, RabbitMQ password
	Vhost    string // required,

	// Connection: ChannelPoolSize specifies the size of a channel pool. All channels in the pool are opened when the publisher is
	// created. When a Publish operation is performed then a channel is taken from the pool to perform the operation and
	// then returned to the pool once the operation has finished. If all channels are in use then the Publish operation
	// waits until a channel is returned to the pool.
	// If this value is set to 0 (default) then channels are not pooled and a new channel is opened/closed for every
	// Publish operation.
	ChannelPoolSize int
}

const (
	configSection = "sdk.rabbitmq"
)

var (
	defaultConfig = &RabbitMqConfig{
		Connection: &ConnectionConfig{
			Port:            5672,
			Username:        "guest",
			Password:        "guest",
			Vhost:           "/",
			ChannelPoolSize: 10,
		},
		RequeueInFailure: true,
		PrefetchCount:    2,
		ExchangeType:     "fanout",
	}
)

func newConfig(configProvider intf.ConfigProvider) (*RabbitMqConfig, error) {
	if configProvider == nil {
		return nil, errdef.ErrInvalidConfig
	}

	var c *RabbitMqConfig
	err := configProvider.Unmarshal(&c, configSection)
	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, errdef.ErrEmptyConfig
	}

	err = c.validate()
	if err != nil {
		return nil, errors.Wrap(err, "validate rabbitmq provider config")
	}

	return c, nil
}

func (c *RabbitMqConfig) validate() error {
	if c.Connection == nil || c.Connection.Host == "" {
		return errdef.ErrInvalidConfig
	}

	if c.Connection.Port == 0 {
		c.Connection.Port = defaultConfig.Connection.Port
	}

	if c.Connection.Username == "" {
		c.Connection.Username = defaultConfig.Connection.Username
	}

	if c.Connection.Password == "" {
		c.Connection.Password = defaultConfig.Connection.Password
	}

	if c.Connection.Vhost == "" {
		c.Connection.Vhost = defaultConfig.Connection.Vhost
	}

	return nil
}
