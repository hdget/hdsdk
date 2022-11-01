package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"hdsdk/types"
	"hdsdk/utils"
	"strings"
	"sync/atomic"
	"time"
)

// 消息队列客户端维护connection和channel
type BaseClient struct {
	Logger types.LogProvider

	Name    string // 客户端名字
	Url     string // 连接url
	Options map[types.MqOptionType]types.MqOptioner

	Connection *amqp.Connection
	Channel    *amqp.Channel
	closed     int32

	chanReconnect chan interface{} // 发送确认通道
}

// 重连等待的最大等待时间
const (
	MAX_CLIENT_RECONNECT_WAIT_TIME = 10 * time.Second
)

// 支持的exchangeTypes
var (
	ExchangeTypeDirect      = "direct"
	ExchangeTypeFanout      = "fanout"
	ExchangeTypeTopic       = "topic"
	ExchangeTypeDelayDirect = "delay:direct"
	ExchangeTypeDelayFanout = "delay:fanout"
	ExchangeTypeDelayTopic  = "delay:topic"

	SupportedExchangeTypes = []string{
		ExchangeTypeDirect,
		ExchangeTypeFanout,
		ExchangeTypeTopic,
		ExchangeTypeDelayDirect,
		ExchangeTypeDelayFanout,
		ExchangeTypeDelayTopic,
	}
)

func (rmq *RabbitMq) newBaseClient(options ...types.MqOptioner) *BaseClient {
	allOptions := rmq.GetDefaultOptions()
	for _, option := range options {
		allOptions[option.GetType()] = option
	}

	// 连接URL
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", rmq.Config.Username, rmq.Config.Password, rmq.Config.Host, rmq.Config.Port, rmq.Config.Vhost)
	return &BaseClient{
		Logger:        rmq.Logger,
		Url:           url,
		Options:       allOptions,
		chanReconnect: make(chan interface{}),
	}
}

// 重连
func (c *BaseClient) connect() error {
	// 建立连接
	var err error
	if c.Connection == nil || c.Connection.IsClosed() {
		c.Connection, err = amqp.Dial(c.Url)
		if err != nil {
			return err
		}
	}

	// 创建channel
	c.Channel, err = c.Connection.Channel()
	if err != nil {
		e := c.Connection.Close()
		if e != nil {
			c.Logger.Error("close connection if error create channel", "name", c.Name, "err", e)
		}
		return err
	}
	return nil
}

// nolint:errcheck
// 添加事件监听, 处理connection和channel的notify close消息
// 同时如果有正常的quit事件，退出事件监听
func (c *BaseClient) addEventListener() {
	chanCloseNotify := c.Channel.NotifyClose(make(chan *amqp.Error))

	go func(closeNotifies chan *amqp.Error) {
		for {
			// 如果客户端关闭了，不需要再监听closeNotify了
			if c.isClosed() {
				break
			}

			// 遍历打印所有close event，直到chanCloseNotify被关闭
			// IMPORTANT: 必须清空Notify，否则死连接不会释放
			// 这里用for循环来读取通道中的内容防止其阻塞Close动作
			// 因为同时amqp在尝试关闭这些通道，在某个时刻chanCloseNotify一定会收到nil
			for closeError := range closeNotifies {
				if closeError != nil {
					c.Logger.Error("notify error", "name", c.Name, "err", closeError)
				} else {
					c.Logger.Debug("receive close event", "name", c.Name)
				}
			}

			// 如果client没被关闭，需要尝试重连
			if !c.isClosed() {
				// 重连
				c.reconnect()

				// 重连成功通知
				c.chanReconnect <- nil

				// IMPORTANT: 这里注意每次for循环需要重新notifyClose注册listener
				// 检查是否收到notify close消息，打印相关错误信息
				// 如果是quit, 退出事件监听, 这里只在channel上监听了close事件
				// 因为connection的close消息同样会广播到channel上来
				closeNotifies = c.Channel.NotifyClose(make(chan *amqp.Error))
			}
		}
	}(chanCloseNotify)
}

// 声明exchange
// @return exchangeName
// @return exchangeType
// @return error
func (c *BaseClient) setupExchange(exchangeName, exchangeType string) error {
	if !utils.Contains(SupportedExchangeTypes, strings.ToLower(exchangeType)) {
		return fmt.Errorf("unsupported exchange type: %s", exchangeType)
	}

	// 如果未指定exchangeName, 缺省使用default exchange, 不需要declare
	if exchangeName == "" {
		return nil
	}

	// 如果指定了exchangeName, 尝试检测exchange是否声明了, 如果已经声明了的话就会无错误
	option := GetExchangeOption(c.Options)
	if strings.HasPrefix(exchangeType, "delay:") {
		routeType := exchangeType[len("delay:"):]
		if !utils.Contains(SupportedExchangeTypes, strings.ToLower(routeType)) {
			return fmt.Errorf("unsupported route type: %s", routeType)
		}

		if option.Args == nil {
			option.Args = amqp.Table{"x-delayed-type": routeType}
		} else {
			option.Args["x-delayed-type"] = routeType
		}
		exchangeType = "x-delayed-message"
	}

	err := c.Channel.ExchangeDeclarePassive(
		exchangeName,      // exchange name
		exchangeType,      // exchange type
		option.Durable,    // durable
		option.AutoDelete, // auto-deleted
		option.Internal,   // internal
		option.NoWait,     // no-wait
		option.Args,       // arguments
	)
	// 如果被动声明没成功，尝试显示声明，这里注意上述隐式声明出错，channel可能已经被关闭了
	if err != nil {
		// 因为之前出错会关闭channel, 这里尝试尝试重连
		err := c.connect()
		if err != nil {
			return err
		}

		if err := c.Channel.ExchangeDeclare(
			exchangeName,      // exchange name
			exchangeType,      // exchange type
			option.Durable,    // durable
			option.AutoDelete, // auto-deleted
			option.Internal,   // internal
			option.NoWait,     // no-wait
			option.Args,       // arguments
		); err != nil {
			return err
		}
	}

	return nil
}

// 等待重连
func (c BaseClient) wait(countRetry int) {
	// 计算每次重连的等待时间
	waitTime := time.Duration(countRetry) * time.Second
	if waitTime > MAX_CLIENT_RECONNECT_WAIT_TIME {
		waitTime = MAX_CLIENT_RECONNECT_WAIT_TIME
	}
	time.Sleep(waitTime)
}

// 等待重连
func (c *BaseClient) reconnect() {
	// 尝试重连
	countRetry := 0
	for {
		// 每次需要重连和重新初始化，需要计数+1
		countRetry += 1
		// 尝试重连
		err := c.connect()
		// 如果重连失败，尝试等待一段时间后继续重连
		if err != nil {
			c.Logger.Error("reconnect", "name", c.Name, "retry", countRetry, "err", err)
			c.wait(countRetry)
			continue
		}

		// 如果到这里，说明重连成功，需要跳出for循环
		break
	}

	c.Logger.Debug("successfully reconnected")
}

func (c *BaseClient) close() {
	// 标志client为关闭状态
	c.setClosed()

	// 触发connection close事件,
	// 后续会被closeNotify监听收到, 但因为设置为关闭状态了，不会重试重连
	if err := c.Connection.Close(); err != nil {
		c.Logger.Error("client close: close connection", "err", err)
	}
}

func (c *BaseClient) isClosed() bool {
	return atomic.LoadInt32(&c.closed) == 1
}

func (c *BaseClient) setClosed() {
	atomic.StoreInt32(&c.closed, 1)
}
