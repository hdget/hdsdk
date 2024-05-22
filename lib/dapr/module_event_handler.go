package dapr

import (
	"context"
	"github.com/dapr/go-sdk/service/common"
	"github.com/hdget/hdsdk/v2"
	"github.com/hdget/hdutils"
)

type EventHandler interface {
	GetTopic() string
	GetEventFunction() common.TopicEventHandler
}

type eventHandlerImpl struct {
	module EventModule
	topic  string        // 订阅主题
	fn     EventFunction // 调用函数
}

type EventFunction func(ctx context.Context, event *common.TopicEvent) (retry bool, err error)

func (h eventHandlerImpl) GetTopic() string {
	return h.topic
}

func (h eventHandlerImpl) GetEventFunction() common.TopicEventHandler {
	return func(ctx context.Context, event *common.TopicEvent) (bool, error) {
		// 挂载defer函数
		defer func() {
			if r := recover(); r != nil {
				hdutils.RecordErrorStack(h.module.GetApp())
			}
		}()

		// 强制设置超时时间
		ctx, cancel := context.WithTimeout(ctx, h.module.GetConsumerTimeout())
		defer cancel()

		// 执行具体的函数
		retry, err := h.fn(ctx, event)
		if err != nil {
			req := []rune(hdutils.BytesToString(event.RawData))
			if len(req) > maxRequestLength {
				req = append(req[:maxRequestLength], []rune("...")...)
			}
			hdsdk.Logger().Error("event invoke", "module", h.module.GetMeta().StructName, "topic", h.GetTopic(), "handler", hdutils.Reflect().GetFuncName(h.fn), "req", req, "err", err)
		}

		return retry, err
	}
}
