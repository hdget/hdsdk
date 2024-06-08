package rabbitmq

import (
	"github.com/hdget/hdsdk/v2/errdef"
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
)

type RabbitMqConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Vhost    string `mapstructure:"vhost"`

	// Consumer: Whether or not to requeue when sending a negative acknowledgement in case of a failure.
	RequeueInFailure bool `mapstructure:"requeue_in_failure"`
	// With a prefetch count greater than zero, the server will deliver that many
	// messages to consumers before acknowledgments are received.  The server ignores
	// this option when consumers are started with noAck because no acknowledgments
	// are expected or sent.

	// Consumer: In order to defeat that we can set the prefetch count with the value of 1.
	// This tells RabbitMQ not to give more than one message to a worker at a time.
	// Or, in other words, don't dispatch a new message to a worker until it has
	// processed and acknowledged the previous one.
	// Instead, it will dispatch it to the next worker that is not still busy.
	PrefetchCount int `mapstructure:"prefetch_count"`

	// connection: ChannelPoolSize specifies the size of a channel pool. All channels in the pool are opened when the publisher is
	// created. When a Publish operation is performed then a channel is taken from the pool to perform the operation and
	// then returned to the pool once the operation has finished. If all channels are in use then the Publish operation
	// waits until a channel is returned to the pool.
	// If this value is set to 0 (default) then channels are not pooled and a new channel is opened/closed for every
	// Publish operation.
	ChannelPoolSize int `mapstructure:"channel_pool_size"`
}

type ConnectionConfig struct {
	Host     string // required, RabbitMQ host
	Port     int    // required, RabbitMQ port
	Username string // required, RabbitMQ username
	Password string // required, RabbitMQ password
	Vhost    string // required,

}

const (
	configSection = "sdk.rabbitmq"
)

var (
	defaultConfig = RabbitMqConfig{
		Port:             5672,
		Username:         "guest",
		Password:         "guest",
		Vhost:            "/",
		ChannelPoolSize:  10,
		RequeueInFailure: true,
		PrefetchCount:    2,
	}
)

func newConfig(configProvider intf.ConfigProvider) (*RabbitMqConfig, error) {
	if configProvider == nil {
		return nil, errdef.ErrInvalidConfig
	}

	// copy default config
	var c RabbitMqConfig
	_ = copier.Copy(&c, &defaultConfig)

	// unmarshal config
	err := configProvider.Unmarshal(&c, configSection)
	if err != nil {
		return nil, err
	}

	err = c.validate()
	if err != nil {
		return nil, errors.Wrap(err, "validate rabbitmq provider config")
	}

	return &c, nil
}

func (c *RabbitMqConfig) validate() error {
	if c.Host == "" {
		return errdef.ErrInvalidConfig
	}

	if c.Port == 0 {
		c.Port = defaultConfig.Port
	}

	if c.Username == "" {
		c.Username = defaultConfig.Username
	}

	if c.Password == "" {
		c.Password = defaultConfig.Password
	}

	if c.Vhost == "" {
		c.Vhost = defaultConfig.Vhost
	}

	return nil
}
