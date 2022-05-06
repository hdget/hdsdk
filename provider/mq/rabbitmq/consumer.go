package rabbitmq

import (
	"github.com/hdget/hdsdk/types"
	"github.com/streadway/amqp"
	"sync/atomic"
)

type RabbitMqConsumer struct {
	Logger        types.LogProvider
	ConsumeOption *ConsumeOption
	QosOption     *QosOption
	Client        *ConsumerClient
	QueueName     string
	Process       types.MqMsgProcessFunc
	closed        int32
}

var _ types.MqConsumer = (*RabbitMqConsumer)(nil)

// CreateConsumer producer的名字和route中的名字对应
func (rmq *RabbitMq) CreateConsumer(name string, processFunc types.MqMsgProcessFunc, args ...map[types.MqOptionType]types.MqOptioner) (types.MqConsumer, error) {
	options := rmq.GetDefaultOptions()
	if len(args) > 0 {
		options = args[0]
	}

	// 初始化rabbitmq client
	client, err := rmq.NewConsumerClient(name, options)
	if err != nil {
		return nil, err
	}

	// 连接
	err = client.connect()
	if err != nil {
		return nil, err
	}

	// exchange不管producer还是consumer都必须要设置好
	err = client.setupExchange(client.Config.ExchangeName, client.Config.ExchangeType)
	if err != nil {
		return nil, err
	}

	qname, err := client.setupQueue()
	if err != nil {
		return nil, err
	}

	// 添加事件监听, 一定要确保这个eventListener必须在exchange和queue设置之后
	client.addEventListener()

	c := &RabbitMqConsumer{
		Logger:        rmq.Logger,
		Client:        client,
		ConsumeOption: getConsumeOption(options),
		QosOption:     getQosOption(options),
		QueueName:     qname,
		Process:       processFunc,
	}

	return c, nil
}

// Consume 消费消息
// consume -> spawn a routing -> handle(deliveries)
func (mc *RabbitMqConsumer) Consume() {
	countRetry := 0
	for !mc.isClosed() {
		// 防止有大量消息堆积或者producer生产数据的能力比消费者强很多的时候导致消费端出现问题
		// 需要进行消费者流控，这个时候消息的消费需要配合Qos
		err := mc.Client.Channel.Qos(
			// 每次队列只消费一个消息 这个消息处理不完服务器不会发送第二个消息过来
			// 当前消费者一次能接受的最大消息数量
			mc.QosOption.PrefetchCount,
			//服务器传递的最大容量
			mc.QosOption.PrefetchSize,
			//如果为true 对channel可用 false则只对当前队列可用
			mc.QosOption.Global,
		)
		// 如果发现consume错误，等待一定时间让client重连成功
		if err != nil {
			mc.Logger.Error("aliyun setup qos", "name", mc.Client.Name, "retry", countRetry, "err", err)
		}

		messages, err := mc.Client.Channel.Consume(
			mc.QueueName,               // queue
			mc.Client.Name,             // 消费者标志用client.Name来标识
			mc.ConsumeOption.NoAck,     // auto-ack
			mc.ConsumeOption.Exclusive, // exclusive
			mc.ConsumeOption.NoLocal,   // no-local
			mc.ConsumeOption.NoWait,    // no-wait
			mc.ConsumeOption.Arguments, // args
		)
		// 如果发现consume错误，等待一定时间让client重连成功
		if err != nil {
			countRetry += 1
			mc.Logger.Error("aliyun consume", "name", mc.Client.Name, "retry", countRetry, "err", err)
			mc.Client.wait(countRetry)
			continue
		}

		// 在另外一个routine中处理消息
		chanDone := make(chan interface{})
		// 处理消息
		go mc.handle(messages, chanDone)

		// chanDone在下面两种情况下会收到消息:
		// 1. channel或者connection的close会关闭deliveries
		// 2. shutdown consumer也会关闭deliveries
		// 只有当handle处理消息的时候退出的时候才会继续下一次for循环
		<-chanDone
		mc.Logger.Debug("aliyun consume quit", "name", mc.Client.Name)
	}
}

// 实际处理消息
func (mc *RabbitMqConsumer) handle(deliveries <-chan amqp.Delivery, done chan interface{}) {
	for msg := range deliveries {
		var err error
		action := mc.Process(msg.Body)
		switch action {
		case types.Ack:
			err = msg.Ack(false)
		case types.Retry:
			err = msg.Nack(false, true)
		case types.BatchAck:
			err = msg.Ack(true)
		case types.BatchRetry:
			err = msg.Nack(true, true)
		default:
			mc.Logger.Error("handle: unsupported action", "action", action)
			// if it is non-supported action type, need retry
			err = msg.Nack(false, true)
		}
		if err != nil {
			mc.Logger.Error("handle: process deliver", "action", action, "err", err)
		}
	}
	mc.Logger.Debug("handle: deliveries channel closed")

	close(done)
}

func (mc *RabbitMqConsumer) isClosed() bool {
	return atomic.LoadInt32(&mc.closed) == 1
}

func (mc *RabbitMqConsumer) setClosed() {
	atomic.StoreInt32(&mc.closed, 1)
}

func (mc *RabbitMqConsumer) shutdown() {
	// 先标识consumer处于关闭状态
	mc.setClosed()

	// 关闭consumer delivery, 这里用client.Name作为consumer的customerFlag
	// consumer.consume()也会使用同样的参数，这样保证其cancel掉的是同一个对象
	if err := mc.Client.Channel.Cancel(mc.Client.Name, false); err != nil {
		mc.Logger.Error("cancel consumer delivery", "name", mc.Client.Name, "err", err)
	}

	mc.Client.close()
}

func (mc *RabbitMqConsumer) Close() {
	mc.shutdown()
}
