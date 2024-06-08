package rabbitmq

import (
	"fmt"
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync/atomic"
)

// pooledChannelManager maintains a pool of channels which are opened immediately upon creation of the provider.
// The GetChannel() function returns an existing channel from the pool. If no channel is available then the caller must
// wait until a channel is returned to the pool (with the CloseChannel function). Channels in the pool are closed when
// this manager's Close() function is called.
// This manager improves performance in high volume systems and also acts as a throttle to prevent the AMQP server from
// overloading.
type pooledChannelManager struct {
	logger     intf.LoggerProvider
	connection *connection
	channels   []*pooledChannelImpl
	closed     uint32
	chanPool   chan *pooledChannelImpl
	closedChan chan struct{}
}

func newPooledChannelManager(logger intf.LoggerProvider, connManager *connection, poolSize int) (channelManager, error) {
	logger.Info("creating pooled channel manager", "poolSize", poolSize)

	channels := make([]*pooledChannelImpl, poolSize)

	chanPool := make(chan *pooledChannelImpl, poolSize)

	// Create the channels and add them to the pool.
	for i := 0; i < poolSize; i++ {
		c, err := newPooledChannel(logger, connManager)
		if err != nil {
			return nil, err
		}
		channels[i] = c
		chanPool <- c
	}

	return &pooledChannelManager{
		logger,
		connManager,
		channels,
		0,
		chanPool,
		make(chan struct{}),
	}, nil
}

func (m *pooledChannelManager) GetChannel() (channel, error) {
	if m.isClosed() {
		return nil, errors.New("channel pool is closed")
	}

	select {
	case c := <-m.chanPool:
		// Ensure that the existing AMQP channel is still open.
		if err := c.validate(); err != nil {
			return nil, err
		}

		return c, nil
	case <-m.closedChan:
		return nil, errors.New("pooled channel manager is closed")
	}
}

func (m *pooledChannelManager) CloseChannel(c channel) error {
	if m.isClosed() {
		return nil
	}

	pc, ok := c.(*pooledChannelImpl)
	if !ok {
		return errors.New("channel must be of type pooledChannelImpl")
	}

	m.chanPool <- pc

	return nil
}

func (m *pooledChannelManager) Close() {
	if !atomic.CompareAndSwapUint32(&m.closed, 0, 1) {
		// Already closed.
		return
	}

	close(m.closedChan)

	m.logger.Info("closing all channels in the pool", "poolSize", len(m.channels))

	for _, c := range m.channels {
		if err := c.Close(); err != nil {
			m.logger.Error("closing channel: %s", "err", err)
		}
	}
}

func (m *pooledChannelManager) isClosed() bool {
	return atomic.LoadUint32(&m.closed) != 0
}

type pooledChannelImpl struct {
	logger      intf.LoggerProvider
	connection  *connection
	amqpChan    *amqp.Channel
	closedChan  chan *amqp.Error
	confirmChan chan amqp.Confirmation
}

func newPooledChannel(logger intf.LoggerProvider, conn *connection) (*pooledChannelImpl, error) {
	c := &pooledChannelImpl{
		logger,
		conn,
		nil,
		nil,
		nil,
	}

	if err := c.openAMQPChannel(); err != nil {
		return nil, fmt.Errorf("open AMQP channel: %w", err)
	}

	return c, nil
}

func (c *pooledChannelImpl) AMQPChannel() *amqp.Channel {
	return c.amqpChan
}

func (c *pooledChannelImpl) Delivered() bool {
	if c.confirmChan == nil {
		// Delivery confirmation is not enabled. Simply return true.
		return true
	}

	confirmed := <-c.confirmChan

	return confirmed.Ack
}

// DeliveryConfirmationEnabled returns true if delivery confirmation of published messages is enabled.
func (c *pooledChannelImpl) DeliveryConfirmationEnabled() bool {
	return c.confirmChan != nil
}

func (c *pooledChannelImpl) Close() error {
	return c.amqpChan.Close()
}

func (c *pooledChannelImpl) validate() error {
	select {
	case e := <-c.closedChan:
		c.logger.Info("AMQP channel was closed. opening new channel.", "close-error", e.Error())
		return c.openAMQPChannel()
	default:
		return nil
	}
}

func (c *pooledChannelImpl) openAMQPChannel() error {
	var err error

	c.amqpChan, err = c.connection.AmqpConnection().Channel()
	if err != nil {
		return fmt.Errorf("create AMQP channel: %w", err)
	}

	c.closedChan = make(chan *amqp.Error, 1)

	c.amqpChan.NotifyClose(c.closedChan)

	err = c.amqpChan.Confirm(false)
	if err != nil {
		return fmt.Errorf("set AMQP channel to confirmed mode: %w", err)
	}

	c.confirmChan = c.amqpChan.NotifyPublish(make(chan amqp.Confirmation, 1))

	return nil
}
