package kafkago

//
//import (
//	"errors"
//	"github.com/hdget/hdsdk/types"
//	"github.com/segmentio/kafka-go"
//)
//
//type Producer struct {
//	Logger types.LogProvider
//	Option *PublishOption
//	Client *ProducerClient
//}
//
//var (
//	ErrInvalidBalancer = errors.New("invalid balancer")
//)
//
//var _ types.MqProducer = (*Producer)(nil)
//
//// CreateProducer 创造一个生产者
//func (k *Kafka) CreateProducer(name string, args ...types.MqOptioner) (types.MqProducer, error) {
//	// 初始化rabbitmq client
//	client, err := k.newProducerClient(name, args...)
//	if err != nil {
//		return nil, err
//	}
//
//	// 客户端连接
//	err = client.connect(k.Config.Brokers)
//	if err != nil {
//		return nil, err
//	}
//
//	p := &Producer{
//		Logger: k.Logger,
//		Client: client,
//	}
//
//	return p, nil
//}
//
//func (p *Producer) GetLastConfirmedId() uint64 {
//	return 0
//}
//
//func (p *Producer) Close() {
//	p.Client.Writer.Close()
//}
//
//func (p Producer) Publish(data []byte, args ...interface{}) error {
//	var key []byte
//	if len(args) > 0 {
//		v, ok := args[0].([]byte)
//		if ok {
//			key = v
//		}
//	}
//
//	msgs := make([]kafka.Message, 0)
//	for _, topic := range p.Client.Config.Topics {
//		m := kafka.Message{
//			Topic: topic,
//			Key:   key,
//			Value: data,
//		}
//		msgs = append(msgs, m)
//	}
//
//	return p.Client.Writer.WriteMessages(p.Client.ctx, msgs...)
//}
//
//func (p Producer) PublishDelay(data []byte, ttl int64, args ...interface{}) error {
//	return nil
//}
