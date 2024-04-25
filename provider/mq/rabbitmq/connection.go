package rabbitmq

import (
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
	"sync/atomic"
)

// connection manages an AMQP connection.
type connection struct {
	config             *ConnectionConfig
	logger             intf.LoggerProvider
	amqpConnection     *amqp.Connection
	amqpConnectionLock sync.Mutex
	connected          chan struct{}

	closing chan struct{}
	closed  uint32

	connectionWaitGroup sync.WaitGroup
}

const (
	amqpURITemplate = "amqp://%s:%s@%s:%d/%s"
)

var (
	defaultBackoffConfig = backoff.NewExponentialBackOff()
)

// newConnection returns a new connection manager.
func newConnection(logger intf.LoggerProvider, connectionConfig *ConnectionConfig) (*connection, error) {
	cm := &connection{
		config:    connectionConfig,
		logger:    logger,
		closing:   make(chan struct{}),
		connected: make(chan struct{}),
	}
	if err := cm.connect(); err != nil {
		return nil, err
	}

	// fork a go routine to monitor the connection close event
	go cm.handleConnectionClose()

	return cm, nil
}

func (c *connection) connect() error {
	c.amqpConnectionLock.Lock()
	defer c.amqpConnectionLock.Unlock()

	var err error

	c.amqpConnection, err = amqp.Dial(c.getURI())
	if err != nil {
		return fmt.Errorf("cannot connect to AMQP, uri: %s", c.getSecuredURI())
	}

	close(c.connected)

	c.logger.Info("connected to AMQP", "uri", c.getSecuredURI())

	return nil
}

func (c *connection) handleConnectionClose() {
	for {
		c.logger.Debug("waiting for AMQP connection to be created")
		<-c.connected
		c.logger.Debug("monitor AMQP connection close event", nil)

		notifyCloseConnection := c.amqpConnection.NotifyClose(make(chan *amqp.Error))

		select {
		case <-c.closing:
			c.logger.Debug("stopping connection")
			c.connected = make(chan struct{})
			return
		case err := <-notifyCloseConnection:
			c.connected = make(chan struct{})
			c.logger.Error("received close notification from AMQP, reconnecting", "err", err)
			c.reconnect()
		}
	}
}

func (c *connection) reconnect() {
	err := backoff.Retry(func() error {
		err := c.connect()
		if err == nil {
			return nil
		}

		c.logger.Error("cannot reconnect to AMQP, retrying", "err", err)

		if c.IsClosed() {
			return backoff.Permanent(errors.Wrap(err, "closing AMQP connection"))
		}

		return err
	}, defaultBackoffConfig)
	if err != nil {
		// should only exit, if closing Pub/Sub
		c.logger.Error("AMQP reconnect failed", "err", err)
	}
}

func (c *connection) Close() error {
	if !atomic.CompareAndSwapUint32(&c.closed, 0, 1) {
		// Already closed.
		return nil
	}

	close(c.closing)

	c.logger.Info("closing AMQP connection")
	defer c.logger.Info("closed AMQP connection")

	c.connectionWaitGroup.Wait()

	if err := c.amqpConnection.Close(); err != nil {
		c.logger.Error("connection close error", "err", err)
	}

	return nil
}

func (c *connection) IsClosed() bool {
	return atomic.LoadUint32(&c.closed) == 1
}

func (c *connection) IsConnected() bool {
	select {
	case <-c.connected:
		return true
	default:
		return false
	}
}

func (c *connection) getURI() string {
	return fmt.Sprintf(amqpURITemplate, c.config.Username, c.config.Password, c.config.Host, c.config.Port, c.config.Vhost)
}

func (c *connection) getSecuredURI() string {
	return fmt.Sprintf(amqpURITemplate, c.config.Username, "***", c.config.Host, c.config.Port, c.config.Vhost)
}
