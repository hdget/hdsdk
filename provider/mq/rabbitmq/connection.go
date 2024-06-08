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
	config              *RabbitMqConfig
	logger              intf.LoggerProvider
	amqpConnection      *amqp.Connection
	connectionWaitGroup sync.WaitGroup
	amqpConnectionLock  sync.Mutex
	chanConnected       chan struct{}
	chanClosing         chan struct{}
	chanClosed          uint32
}

const (
	amqpURITemplate = "amqp://%s:%s@%s:%d/%s"
)

var (
	defaultBackoffConfig = backoff.NewExponentialBackOff()
)

// newConnection returns a new connection manager.
func newConnection(logger intf.LoggerProvider, c *RabbitMqConfig) (*connection, error) {
	cm := &connection{
		config:        c,
		logger:        logger,
		chanClosing:   make(chan struct{}),
		chanConnected: make(chan struct{}),
	}
	if err := cm.connect(); err != nil {
		return nil, err
	}

	// fork a go routine to monitor the connection close event
	go cm.handleConnectionClose()

	return cm, nil
}

// Begin 连接开始执行逻辑
func (c *connection) Begin() {
	c.connectionWaitGroup.Add(1)
}

// End 连接结束执行完逻辑
func (c *connection) End() {
	c.connectionWaitGroup.Done()
}

// Connected 是否已经连接
func (c *connection) Connected() chan struct{} {
	return c.chanConnected
}

// Closing 正在关闭
func (c *connection) Closing() chan struct{} {
	return c.chanClosing
}

// AmqpConnection 正在关闭
func (c *connection) AmqpConnection() *amqp.Connection {
	return c.amqpConnection
}

func (c *connection) connect() error {
	c.amqpConnectionLock.Lock()
	defer c.amqpConnectionLock.Unlock()

	var err error

	c.amqpConnection, err = amqp.Dial(c.getURI())
	if err != nil {
		return fmt.Errorf("cannot connect to AMQP, uri: %s", c.getSecuredURI())
	}

	close(c.chanConnected)

	c.logger.Info("Connected to AMQP", "uri", c.getSecuredURI())

	return nil
}

func (c *connection) handleConnectionClose() {
	for {
		c.logger.Debug("waiting for AMQP connection to be created")
		<-c.chanConnected
		c.logger.Debug("monitor AMQP connection close event", nil)

		notifyCloseConnection := c.amqpConnection.NotifyClose(make(chan *amqp.Error))

		select {
		case <-c.chanClosing:
			c.logger.Debug("stopping connection")
			c.chanConnected = make(chan struct{})
			return
		case err := <-notifyCloseConnection:
			c.chanConnected = make(chan struct{})
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
			return backoff.Permanent(errors.Wrap(err, "chanClosing AMQP connection"))
		}

		return err
	}, defaultBackoffConfig)
	if err != nil {
		// should only exit, if Closing Pub/Sub
		c.logger.Error("AMQP reconnect failed", "err", err)
	}
}

func (c *connection) Close() error {
	if !atomic.CompareAndSwapUint32(&c.chanClosed, 0, 1) {
		// Already Closed.
		return nil
	}

	close(c.chanClosing)

	c.logger.Info("chanClosing AMQP connection")
	defer c.logger.Info("chanClosed AMQP connection")

	c.connectionWaitGroup.Wait()

	if err := c.amqpConnection.Close(); err != nil {
		c.logger.Error("connection close error", "err", err)
	}

	return nil
}

func (c *connection) IsClosed() bool {
	return atomic.LoadUint32(&c.chanClosed) == 1
}

func (c *connection) IsConnected() bool {
	select {
	case <-c.chanConnected:
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
