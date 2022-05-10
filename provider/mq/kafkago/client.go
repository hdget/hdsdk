package kafkago

//
//import (
//	"context"
//	"github.com/hdget/hdsdk/types"
//)
//
//// 消息队列客户端维护connection和channel
//type BaseClient struct {
//	Logger types.LogProvider
//
//	Name    string // 客户端名字
//	Options map[types.MqOptionType]types.MqOptioner
//	ctx     context.Context
//}
//
//func (k *Kafka) newBaseClient(name string, options ...types.MqOptioner) *BaseClient {
//	// 构造所有option
//	allOptions := k.GetDefaultOptions()
//	for _, opt := range options {
//		allOptions[opt.GetType()] = opt
//	}
//
//	// 连接URL
//	return &BaseClient{
//		Logger:  k.Logger,
//		Name:    name,
//		Options: allOptions,
//		ctx:     context.Background(),
//	}
//}
