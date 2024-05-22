package dapr

import (
	"context"
	"github.com/dapr/go-sdk/service/common"
	"github.com/hdget/hdsdk/v2"
	"github.com/hdget/hdutils"
	"sync"
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

var (
	// 记录正在处理的消息
	processingMessages = sync.Map{}
)

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

		// 冥等处理，如果消息正正在处理中，直接ACK
		if _, exists := processingMessages.Load(event.ID); exists {
			hdsdk.Logger().Warn("same message received or message is processingMessages", "module", h.module.GetMeta().StructName, "topic", h.GetTopic(), "handler", hdutils.Reflect().GetFuncName(h.fn), "message", trimData(event.RawData))
			return false, nil
		} else {
			processingMessages.Store(event.ID, struct{}{})
		}
		defer processingMessages.Delete(event.ID)

		// 执行具体的函数
		retry, err := h.fn(ctx, event)
		if err != nil {
			hdsdk.Logger().Error("event processingMessages", "module", h.module.GetMeta().StructName, "topic", h.GetTopic(), "handler", hdutils.Reflect().GetFuncName(h.fn), "message", trimData(event.RawData), "err", err)
		}

		return retry, err
	}
}

func trimData(data []byte) []rune {
	trimmed := []rune(hdutils.BytesToString(data))
	if len(trimmed) > maxRequestLength {
		trimmed = append(trimmed[:maxRequestLength], []rune("...")...)
	}
	return trimmed
}
