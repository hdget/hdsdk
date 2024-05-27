package dapr

import (
	"context"
	"fmt"
	"github.com/dapr/go-sdk/service/common"
	"github.com/hdget/hdsdk/v2"
	"github.com/hdget/hdutils"
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
	return func(ctx context.Context, event *common.TopicEvent) (retry bool, err error) {
		defer func() {
			if r := recover(); r != nil {
				hdutils.RecordErrorStack(h.module.GetApp())
			}
			retry = false
			err = fmt.Errorf("%s panic", h.module.GetApp())
		}()

		ctx, cancel := context.WithTimeout(ctx, h.module.GetAckTimeout())
		defer cancel()

		go func() {
			hdutils.LogDebug("xxxxxx1", "instance", hdsdk.Logger())
			hdutils.LogDebug("xxxxxx2", "timeout", h.module.GetAckTimeout())
			select {
			case <-ctx.Done():
				hdutils.LogError("event processing timeout, discard message", "message", event.RawData)
				retry = false
				err = ctx.Err()
			}
		}()

		// 执行具体的函数
		retry, err = h.fn(ctx, event)
		if err != nil {
			hdutils.LogError("event processing", "message", event.RawData, "err", err)
		}
		return

	}
}

func trimData(data []byte) []rune {
	trimmed := []rune(hdutils.BytesToString(data))
	if len(trimmed) > maxRequestLength {
		trimmed = append(trimmed[:maxRequestLength], []rune("...")...)
	}
	return trimmed
}
