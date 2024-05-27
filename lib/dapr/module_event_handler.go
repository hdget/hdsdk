package dapr

import (
	"context"
	"github.com/dapr/go-sdk/service/common"
	"github.com/hdget/hdsdk/v2"
	"github.com/hdget/hdutils"
	"github.com/pkg/errors"
	"time"
)

type eventHandler interface {
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
		defer func() {
			if r := recover(); r != nil {
				hdutils.RecordErrorStack(h.module.GetApp())
			}
		}()

		ctx, cancel := context.WithTimeout(ctx, h.module.GetAckTimeout()-1*time.Minute)
		defer cancel()

		// 执行具体的函数
		retry, err := h.fn(ctx, event)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				hdsdk.Logger().Error("event processing timeout", "module", h.module.GetMeta().StructName, "topic", h.GetTopic(), "handler", hdutils.Reflect().GetFuncName(h.fn), "message", trimData(event.RawData))
				return false, nil
			}
			hdsdk.Logger().Error("event processing", "module", h.module.GetMeta().StructName, "topic", h.GetTopic(), "handler", hdutils.Reflect().GetFuncName(h.fn), "message", trimData(event.RawData))
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
