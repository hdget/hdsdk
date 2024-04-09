package intf

//
//type MqProvider interface {
//	Provider
//	My() Mq
//	By(string) Mq // 获取某个Mq能力提供者
//}
//
//// MqMsgProcessFunc MQ的消息处理函数
//// 入参为处理的数据
//// 第一个返回值是否要ack
//// 第二个返回值是否要把同一个channel上的上一次ack之前的所有
//type MqMsgProcessFunc func([]byte) MqMsgAction
//
//// MqMsgAction 消息处理后的动作
//type MqMsgAction int
//
//const (
//	_          MqMsgAction = iota
//	Ack                    // 消息标志处理成功
//	Retry                  // 消息重传进行重新处理，当条消息会被重传
//	Next                   // 取下一条消息
//	BatchAck               // 批量消息标志处理成功
//	BatchRetry             // 批量消息进行重传并重新处理，自上次ack到现在的消息都会被重传
//)
//
//// MqOptioner 选项接口
//type MqOptioner interface {
//	GetType() MqOptionType // 获取配置项类型，现在有几个配置项: exchange配置项, queue配置项, publish配置项
//}
//
//type Mq interface {
//	GetDefaultOptions() map[MqOptionType]MqOptioner
//	CreateProducer(parameters map[string]interface{}, args ...MqOptioner) (MqProducer, error)
//	CreateConsumer(processFunc MqMsgProcessFunc, parameters map[string]interface{}, args ...MqOptioner) (MqConsumer, error)
//}
//
//// 消息发布者，负责生产并发送消息至Topic
//type MqProducer interface {
//	Publish(data []byte, args ...interface{}) error                 // MQ发送消息
//	PublishDelay(data []byte, ttl int64, args ...interface{}) error // MQ发送延迟消息
//	GetLastConfirmedId() uint64                                     // 获取上一次确认发送成功的消息Tag
//	Close()
//}
//
//// 消息订阅者，负责从Topic接收并消费消息。
//type MqConsumer interface {
//	Consume() // 消费消息
//	Close()
//}
//
//// option intf
//type MqOptionType int
//
//const (
//	_                MqOptionType = iota
//	MqOptionQueue                 // 队列选项
//	MqOptionExchange              // exchange选项
//	MqOptionPublish               // 发送选项
//	MqOptionConsume               // 消费选项
//	MqOptionQos                   // Qos选项
//)
