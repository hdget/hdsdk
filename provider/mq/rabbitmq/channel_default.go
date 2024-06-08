package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

// defaultChannelProvider simply opens a new channel when channel() is called and closes the channel
// when CloseChannel is called.
type defaultChannelManager struct {
	connection *connection
}

func newDefaultChannelManager(conn *connection) channelManager {
	return &defaultChannelManager{connection: conn}
}

func (m *defaultChannelManager) GetChannel() (channel, error) {
	amqpChan, err := m.connection.AmqpConnection().Channel()
	if err != nil {
		return nil, fmt.Errorf("create AMQP channel: %w", err)
	}

	err = amqpChan.Confirm(false)
	if err != nil {
		return nil, fmt.Errorf("set AMQP channel to confirmed mode: %w", err)
	}

	confirmChan := amqpChan.NotifyPublish(make(chan amqp.Confirmation, 1))

	return &defaultChannelImpl{amqpChan, confirmChan}, nil
}

func (m *defaultChannelManager) CloseChannel(c channel) error {
	return c.Close()
}

func (m *defaultChannelManager) Close() {
	// do nothing
}

type defaultChannelImpl struct {
	*amqp.Channel
	confirmChan chan amqp.Confirmation
}

func (c *defaultChannelImpl) AMQPChannel() *amqp.Channel {
	return c.Channel
}

func (c *defaultChannelImpl) DeliveryConfirmationEnabled() bool {
	return c.confirmChan != nil
}

func (c *defaultChannelImpl) Delivered() bool {
	if c.confirmChan == nil {
		// Delivery confirmation is not enabled. Simply return true.
		return true
	}

	confirmed := <-c.confirmChan

	return confirmed.Ack
}
