package rabbitmq

import (
	"errors"
	"github.com/hdget/hdsdk/types"
	"github.com/hdget/hdsdk/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/streadway/amqp"
	"time"
)

type RabbitMqProducer struct {
	Logger             types.LogProvider
	Option             *PublishOption
	ExchangeName       string
	Client             *ProducerClient
	CurrentDeliveryTag uint64     // 当前确认的deliveryTag
	LastDeliverTag     uint64     // 上一次确认的deliveryTag
	chanDelivery       chan error // 发送递达确认通道
}

const PUBLISH_CONFIRM_RETRY_TIMEOUT = 1 * time.Second

var (
	ErrPublishNotConfirmed = errors.New("publish is not acknowledged")
	ErrPublishAckLost      = errors.New("publish ack lost")
)

var _ types.MqProducer = (*RabbitMqProducer)(nil)

type ProducerParameter struct {
	ExchangeName string `mapstructure:"exchangeName"`
	ExchangeType string `mapstructure:"exchangeType"`
}

// CreateProducer producer的名字和route中的名字对应
func (rmq *RabbitMq) CreateProducer(parameters map[string]interface{}, options ...types.MqOptioner) (types.MqProducer, error) {
	// 初始化rabbitmq client
	client, err := rmq.newProducerClient(options...)
	if err != nil {
		return nil, err
	}

	// 客户端连接
	err = client.connect()
	if err != nil {
		return nil, err
	}

	// exchange不管producer还是consumer都必须要设置好
	producerParam, err := parseProducerParameter(parameters)
	if err != nil {
		return nil, err
	}

	err = client.setupExchange(producerParam.ExchangeName, producerParam.ExchangeType)
	if err != nil {
		return nil, err
	}

	p := &RabbitMqProducer{
		Logger:       rmq.Logger,
		Client:       client,
		Option:       GetPublishOption(client.Options),
		ExchangeName: producerParam.ExchangeName,
		chanDelivery: make(chan error),
	}

	// 一定要确保这些eventListener必须在exchange和queue设置之后
	// 因为在setupExchange和setupQueue的时候我们已经可能导致channel close了
	// 我们在那里已经强制reconnect了, 如果提前设置了closeNotify，会导致混乱

	// 确保publish confirm
	p.addEventListener()

	// 确保断线重连
	client.addEventListener()

	return p, nil
}

// nolint: errcheck
// 监听是否需要关闭producer, 处理消息是否投递成功事件
func (rmp *RabbitMqProducer) addEventListener() {
	// 设置channel进入消息投递确认模式
	//The capacity of the chan Confirmation must be at least as large as the
	//number of outstanding publishings.  Not having enough buffered chans will
	//create a deadlock if you attempt to perform other operations on the Connection
	//or Channel while confirms are in-flight.
	chanPublishConfirm := rmp.Client.Channel.NotifyPublish(make(chan amqp.Confirmation, 1))
	rmp.Client.Channel.Confirm(false)

	// Publish Confirm事件监听
	go func(confirms chan amqp.Confirmation) {
		for {
			select {
			// 如果channel或者connection重连了
			case <-rmp.Client.chanReconnect:
				rmp.Logger.Error("publish confirm: reestablished as channel has reconnected")

				// 重连可能发生多次，如果重连，需要累加上次出错的DeliveryTag
				rmp.LastDeliverTag += rmp.CurrentDeliveryTag

				// 放置在这里因为每次channel关闭并且重连需要重新初始化channel并使新channel进入confirm模式
				// 注意这里需要放在NotifyClose后面,因为这里执行出错会导致channel关闭，后面的select就会检测到并重连
				confirms = rmp.Client.Channel.NotifyPublish(make(chan amqp.Confirmation, 1))
				rmp.Client.Channel.Confirm(false)
				// 在重新setup后需要返回给客户
				rmp.chanDelivery <- ErrPublishAckLost

			// 消息投递成功是否成功, 如果成功返回nil, 不成功给出提示信息
			// 因为网络等其他原因，channel可能被关闭，chanPublishConfirm的isOpen可能会返回false
			case c, isOpen := <-confirms:
				if isOpen {
					// 保存上一次成功的deliveryTag
					rmp.CurrentDeliveryTag = c.DeliveryTag

					// 尝试返回confirm确认给客户端
					var errAck error
					if !c.Ack {
						errAck = ErrPublishNotConfirmed
					}
					rmp.chanDelivery <- errAck
				} else {
					rmp.Logger.Debug("publish confirm: channel is closed")
					time.Sleep(PUBLISH_CONFIRM_RETRY_TIMEOUT)
				}

			}
		}
	}(chanPublishConfirm)
}

func (rmp *RabbitMqProducer) GetLastConfirmedId() uint64 {
	return rmp.CurrentDeliveryTag
}

func (rmp *RabbitMqProducer) Close() {
	rmp.Client.close()
	close(rmp.chanDelivery)
}

// Publish 发布内容到exchange
// 当你想发送一条消息到单个队列的时候，你可以用queueName当做routingKey来发送到缺省的exchange
// 因为每个申明的queue都有一个隐式的route到缺省exchange
//
// mandatory=true
// 当没有队列匹配routingKey, 发布的消息也可能处于不能递交状态
// ---> immediate=true
// 如果在匹配的队列上没有消费者准备好，发布的消息也可能处于不能递交状态
//
// 当connection、channel被关闭的时候都会返回错误，所以我们不能通过判断没有错误来认为服务器已经接收到发布的内容
// 同时因为发布动作是异步的，不能递交的消息会被服务器返回, 我们需要实现Channel.NotifyReturn接口来监听并处理不能递交的消息
//
// 当底层socket被关闭的时候如果没有将等待发送的消息报从内核缓存中进行保存，有可能导致发布内容不能到达broker
// 最简单的防止消息对视就是在终止发布应用的时候需要调用Connection.Close来保证消息不丢失
// 另外为了确保消息到达服务器需要添加一个Channel.NotifyPublish的监听，并且让Channel处于Confirm模式
// 发布递交的标签和对应的确认从1开始, 当所有发布确认后会退出
//
// 当发布没有返回错误并且channel在confirm模式， DeliveryTags的内部计数器首先确认从1开始
func (rmp RabbitMqProducer) Publish(data []byte, args ...interface{}) error {
	key := ""
	if len(args) > 0 {
		v, ok := args[0].(string)
		if ok {
			key = v
		}
	}

	// 尝试发送
	errPublish := rmp.Client.Channel.Publish(
		rmp.ExchangeName,
		key,
		rmp.Option.Mandatory,
		rmp.Option.Immediate,
		amqp.Publishing{
			ContentType:  rmp.Option.ContentType,
			DeliveryMode: rmp.Option.DeliveryMode,
			Body:         data,
		})

	// 等待发送确认
	errAck := <-rmp.chanDelivery
	if errAck != nil {
		return errAck
	}
	return errPublish
}

func (rmp RabbitMqProducer) PublishDelay(data []byte, ttl int64, args ...interface{}) error {
	key := ""
	if len(args) > 0 {
		v, ok := args[0].(string)
		if ok {
			key = v
		}
	}

	// 尝试发送
	errPublish := rmp.Client.Channel.Publish(
		rmp.ExchangeName,
		key,
		rmp.Option.Mandatory,
		rmp.Option.Immediate,
		amqp.Publishing{
			ContentType:  rmp.Option.ContentType,
			DeliveryMode: rmp.Option.DeliveryMode,
			Body:         data,
			Headers: map[string]interface{}{
				"x-delay": ttl, // 消息从交换机过期时间,毫秒（x-dead-message插件提供）
			},
		})

	// 等待发送确认
	errAck := <-rmp.chanDelivery
	if errAck != nil {
		return errAck
	}
	return errPublish
}

func parseProducerParameter(params map[string]interface{}) (*ProducerParameter, error) {
	var producerParams ProducerParameter
	err := mapstructure.Decode(params, &producerParams)
	if err != nil {
		return nil, err
	}

	if producerParams.ExchangeName == "" || !utils.Contains(SupportedExchangeTypes, producerParams.ExchangeType) {
		return nil, ErrInvalidProducerParam
	}

	return &producerParams, nil
}
