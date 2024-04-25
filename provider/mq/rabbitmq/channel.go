package rabbitmq

import (
	"github.com/hdget/hdsdk/v2/intf"
	amqp "github.com/rabbitmq/amqp091-go"
)

type channelManager interface {
	GetChannel() (channel, error)
	CloseChannel(c channel) error
	Close()
}

type channel interface {
	// AMQPChannel returns the underlying AMQP channel.
	AMQPChannel() *amqp.Channel
	// DeliveryConfirmationEnabled returns true if delivery confirmation of published messages is enabled.
	DeliveryConfirmationEnabled() bool
	// Delivered waits until confirmation of delivery has been received from the AMQP server and returns true if delivery
	// was successful, otherwise false is returned. If delivery confirmation is not enabled then true is immediately returned.
	Delivered() bool
	// Close closes the channel.
	Close() error
}

func newChannelManager(logger intf.LoggerProvider, conn *connection, poolSize int) (channelManager, error) {
	if poolSize == 0 {
		return newDefaultChannelManager(conn), nil
	}

	return newPooledChannelManager(logger, conn, poolSize)
}
