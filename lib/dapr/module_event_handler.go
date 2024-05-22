package dapr

import (
	"context"
	"github.com/dapr/go-sdk/service/common"
	"github.com/hdget/hdutils"
)

type EventHandler interface {
	GetTopic() string
	GetPubSub() string
	GetEventFunction() common.TopicEventHandler
}

type eventHandlerImpl struct {
	module Moduler
	pubsub string        // 消息中间件名称定义在dapr配置中
	topic  string        // 订阅主题
	fn     EventFunction // 调用函数
}

type EventFunction func(ctx context.Context, event *common.TopicEvent) (retry bool, err error)

func (h eventHandlerImpl) GetTopic() string {
	return h.topic
}

func (h eventHandlerImpl) GetPubSub() string {
	return h.pubsub
}

func (h eventHandlerImpl) GetEventFunction() common.TopicEventHandler {
	return func(ctx context.Context, event *common.TopicEvent) (retry bool, err error) {
		// 挂载defer函数
		defer func() {
			if r := recover(); r != nil {
				hdutils.RecordErrorStack(h.module.GetApp())
			}
		}()

		return false, nil
	}
}
