package rabbitmq

import (
	"github.com/hdget/hdsdk/types"
	"github.com/streadway/amqp"
)

// 声明exchange, 注意：如果声明出错的话会导致关闭channel
//
// 以"amp."前缀开头的exchange名字为保留名字
//
// 不同类型的exchange定义了消息的不同路由方式，当前有如下类型:
// "direct", "fanout", "topic"和"headers"
//
// Durable=true, AutoDelete=false
// 该类型的exchange在服务重启后会恢复并保持可用即使没有其他binding存在.
// 这种生命周期的exchange是最稳定的并且是缺省的exchange
//
// Durable=false, AutoDelete=true
// 如果没有其他bindings该exchange将会被删除，并且在服务重启后不会被恢复
// 这种生命周期的exchange适用于临时的并且在其他消费者完成后不需要继续保留在virtual host的场景
//
// Durable=false, AutoDelete=false
// 该类型的exchange在服务运行时始终存在即使其没有其他binding存在
// 这种生命周期的exchange适用于在binding之前需要等待较长时间
//
// Durable=true, Auto-Deleted=true
// 这种类型exchange在服务重启后会恢复，但无可用binding的时候会被删除
// 这种生命周期的exchange适用于绑定到持久化的队列但同时又要求其自动会被删除
//
// 注意: RabbitMQ会将'amq.fanout'开头的exchange声明为durable=true，
// 所以绑定到该exchange的队列也需要设置durable=true
//
// internal=true
// 当不想对外公布exchange只想做内部exchange的时候才适用,一般不用
//
// 当noWait=true,在声明的时候不会等待server端的确认，如果有错误会被NotifyClose回调函数监听到
type ExchangeOption struct {
	Durable    bool // 消息是否持久化
	AutoDelete bool // 是否会自动删除：当最后一个消费者断开后，是否将队列中的消息清除
	Internal   bool // 是否具有排他性，意思只有创建者可见，其他人不可用
	NoWait     bool // 是否阻塞
	Args       amqp.Table
}

// 所有declare的queue会获取一个queueBinding,具有以下配置：
// exchangeName="", type="direct", routingKey=queueName
//
// 通过该queueBinding, 我们可以直接通过以下配置的exchange发送消息:
// exchangeName="", routingKey=queueName,
//
// e,g:
//
//	QueueDeclare("alerts", true, false, false, false, nil)
//	Publish("", "alerts", false, false, Publishing{Body: []byte("...")})
//
//	Delivery       Exchange  Key       Queue
//	-----------------------------------------------
//	key: alerts -> ""     -> alerts -> alerts
//
// 如果queueName为空，服务器会生成一个唯一的queueName并通过该Queue结构的Name字段返回
//
// Durable=true,AutoDelete=false
// 该队列不管服务是否重启，也不管是否有消费者或者binding它就始终存在
// 持久的消息将会在服务重启后恢复，注意这些队列只能保存在durable=true的exchange上
//
// Durable=false,AutoDelete=true
// 该队列在服务重启后并不会会重新声明
// 当最后一个消费者取消或者消费者通道关闭后服务器将很快删除该队列
// 这种队列只能保存在durable=false的exchange上
//
// Durable=false,AutoDelete=false
// 该队列在服务运行时将始终处于可用状态不管有多少个消费者
// 这种生命周期的队列通常用来保存可能在不同消费者间存在很久的临时拓扑
// 这种队列只能保存在durable=false的exchange上
//
// Durable=true, AutoDelete=true
// 该队列在服务重启后会恢复，但如果没有活动的消费者该队列将会被移除
// 这种生命周期的队列一般不太会使用
//
// Exclusive=true的队列只能被声明队列的connections使用，当connection关闭的时候队列也会被删除
// 同时如果有其他channel尝试declare,bind,consume,purge或删除同样名字的队列会报错
//
// Nowait=true
// 如果该选项为true, 在声明队列时会假设总是声明成功，
// 如果服务器上已有同样的queue, 或者其他connections尝试修改该队列，channel都会抛出exception
//
// 如果返回的错误不为空，你可以认为该队列用这些配置参数不能声明成功，channel会关闭
type QueueOption struct {
	Durable    bool // 消息是否持久化
	AutoDelete bool // 是否会自动删除：当最后一个消费者断开后，是否将队列中的消息清除
	Exclusive  bool // 是否具有排他性，意思只有创建者可见，其他人不可用
	NoWait     bool // 是否阻塞等待
	Args       amqp.Table
}

// mandatory=true
// 当没有队列匹配routingKey, 发布的消息也可能处于不能递交状态
// immediate=true
// 如果在匹配的队列上没有消费者准备好，发布的消息也可能处于不能递交状态
type PublishOption struct {
	ExchangeName string
	ExchangeType string
	Mandatory    bool
	Immediate    bool

	ContentType  string // MIME content type
	DeliveryMode uint8  // Transient (0 or 1) or Persistent (2)
	//Priority        uint8     // 0 to 9
	//CorrelationId   string    // correlation identifier
	//ReplyTo         string    // address to to reply to (ex: RPC)
	//Expiration      string    // message expiration spec
	//MessageId       string    // message identifier
	//Timestamp       time.Time // message timestamp
	//Type            string    // message type name
	//UserId          string    // creating user id - ex: "guest"
	//AppId           string    // creating application id
}

// mandatory=true
// 当没有队列匹配routingKey, 发布的消息也可能处于不能递交状态
// immediate=true
// 如果在匹配的队列上没有消费者准备好，发布的消息也可能处于不能递交状态
type ConsumeOption struct {
	ExchangeName string
	ExchangeType string
	QueueName    string
	RoutingKeys  []string
	ConsumerTag  string // consumer_tag 消费者标签
	NoLocal      bool   // 这个功能属于AMQP的标准,但是rabbitMQ并没有做实现.
	NoAck        bool   // 收到消息后,是否不需要回复确认即被认为被消费
	Exclusive    bool   // 排他消费者,即这个队列只能由一个消费者消费.适用于任务不允许进行并发处理的情况下.比如系统对接
	NoWait       bool   // 不返回执行结果,但是如果exclusive开启的话,则必须需要等待结果的,如果exclusive和nowait都为true就会报错
	Arguments    amqp.Table
}

type QosOption struct {
	PrefetchCount int
	PrefetchSize  int
	Global        bool
}

var (
	defaultQueueOption = &QueueOption{
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Args:       nil,
	}
	defaultExchangeOption = &ExchangeOption{
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
		Args:       nil,
	}
	defaultPublishOption = &PublishOption{
		Mandatory:    false,
		Immediate:    false,
		ContentType:  "application/json",
		DeliveryMode: 2, // 不设置persistent，broker重启后消息会丢失
	}
	defaultConsumeOption = &ConsumeOption{
		ConsumerTag: "",
		NoLocal:     false,
		NoAck:       false,
		Exclusive:   false,
		NoWait:      false,
		Arguments:   nil,
	}
	// 从测试来看prefetchCount=20会大大提高吞吐量
	// http://www.rabbitmq.com/blog/2012/04/25/rabbitmq-performance-measurements-part-2/
	defaultQosOption = &QosOption{
		PrefetchCount: 20,
		PrefetchSize:  0,
		Global:        false,
	}
)

func (q QueueOption) GetType() types.MqOptionType {
	return types.MqOptionQueue
}

func (q ExchangeOption) GetType() types.MqOptionType {
	return types.MqOptionExchange
}

func (q PublishOption) GetType() types.MqOptionType {
	return types.MqOptionPublish
}

func (q ConsumeOption) GetType() types.MqOptionType {
	return types.MqOptionConsume
}

func (q QosOption) GetType() types.MqOptionType {
	return types.MqOptionQos
}

func GetQueueOption(options map[types.MqOptionType]types.MqOptioner) *QueueOption {
	if len(options) == 0 {
		return defaultQueueOption
	}

	v := options[types.MqOptionQueue]
	if v == nil {
		return defaultQueueOption
	}

	option, ok := v.(*QueueOption)
	if !ok {
		return defaultQueueOption
	}

	return option
}

func GetPublishOption(options map[types.MqOptionType]types.MqOptioner) *PublishOption {
	if len(options) == 0 {
		return defaultPublishOption
	}

	v := options[types.MqOptionPublish]
	if v == nil {
		return defaultPublishOption
	}

	option, ok := v.(*PublishOption)
	if !ok {
		return defaultPublishOption
	}
	return option
}

func GetConsumeOption(options map[types.MqOptionType]types.MqOptioner) *ConsumeOption {
	if len(options) == 0 {
		return defaultConsumeOption
	}

	v := options[types.MqOptionPublish]
	if v == nil {
		return defaultConsumeOption
	}

	option, ok := v.(*ConsumeOption)
	if !ok {
		return defaultConsumeOption
	}
	return option
}

func GetExchangeOption(options map[types.MqOptionType]types.MqOptioner) *ExchangeOption {
	if len(options) == 0 {
		return defaultExchangeOption
	}

	v := options[types.MqOptionExchange]
	if v == nil {
		return defaultExchangeOption
	}

	option, ok := v.(*ExchangeOption)
	if !ok {
		return defaultExchangeOption
	}
	return option
}

func GetQosOption(options map[types.MqOptionType]types.MqOptioner) *QosOption {
	if len(options) == 0 {
		return defaultQosOption
	}

	v := options[types.MqOptionQos]
	if v == nil {
		return defaultQosOption
	}

	option, ok := v.(*QosOption)
	if !ok {
		return defaultQosOption
	}
	return option
}
