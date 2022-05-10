package kafkago

//
//import (
//	"fmt"
//	"github.com/hdget/hdsdk/types"
//	"github.com/segmentio/kafka-go"
//	"github.com/segmentio/kafka-go/sasl/plain"
//	"time"
//)
//
//type ConsumerClient struct {
//	*BaseClient
//	Option *ConsumeOption
//	Config *ConsumerConfig
//	Reader *kafka.Reader
//}
//
//func (k *Kafka) NewConsumerClient(name string, args ...types.MqOptioner) (*ConsumerClient, error) {
//	// 获取匹配的路由配置
//	config := k.getConsumerConfig(name)
//	if config == nil {
//		return nil, fmt.Errorf("no matched consumer config for: %s", name)
//	}
//
//	client := k.newBaseClient(name, args...)
//	return &ConsumerClient{
//		BaseClient: client,
//		Config:     config,
//		Option:     GetConsumeOption(client.Options),
//	}, nil
//}
//
//// connect 连接kafka, 如果args不为空，则args[0]为groupId
//func (cc *ConsumerClient) connect(brokers []string) error {
//	// 检查是否需要创建安全SASL认证
//	var dialer *kafka.Dialer
//	if cc.Config.User != "" && cc.Config.Password != "" {
//		// in aliyun alidts, the username consists of user-groupid
//		m := plain.Mechanism{
//			Username: fmt.Sprintf("%s-%s", cc.Config.User, cc.Config.GroupId),
//			Password: cc.Config.Password,
//		}
//		dialer = &kafka.Dialer{
//			Timeout:       10 * time.Second,
//			DualStack:     true,
//			SASLMechanism: m,
//		}
//	}
//
//	// make a new reader that consumes from topic
//	cc.Reader = kafka.NewReader(
//		kafka.ReaderConfig{
//			Brokers:        brokers,
//			Topic:          cc.Config.Topic,
//			MinBytes:       cc.Option.MinBytes,
//			MaxBytes:       cc.Option.MaxBytes,
//			CommitInterval: time.Duration(cc.Option.CommitInterval) * time.Second,
//			// 下面的是可选的
//			Dialer:  dialer,
//			GroupID: cc.Config.GroupId,
//		},
//	)
//
//	return nil
//}
