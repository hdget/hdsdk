package dapr

import (
	"context"
	"fmt"
	"github.com/dapr/go-sdk/service/common"
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/hdget/hdutils/convert"
	panicUtils "github.com/hdget/hdutils/panic"
	"time"
)

type eventHandler interface {
	GetTopic() string
	GetEventFunction(logger intf.LoggerProvider) common.TopicEventHandler
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

func (h eventHandlerImpl) GetEventFunction(logger intf.LoggerProvider) common.TopicEventHandler {
	return func(ctx context.Context, event *common.TopicEvent) (retry bool, err error) {
		quit := make(chan bool, 1)
		defer func() {
			if r := recover(); r != nil {
				panicUtils.RecordErrorStack(h.module.GetApp())
			}
			quit <- true
			retry = false
			err = fmt.Errorf("%s panic", h.module.GetApp())
		}()

		go func() {
			select {
			case <-time.After(h.module.GetAckTimeout()):
				logger.Error("event processing timeout, discard message", "message", trimData(event.RawData))
				retry = false
				err = ctx.Err()
			case <-quit:
				break
			}
		}()

		// 执行具体的函数
		retry, err = h.fn(ctx, event)
		if err != nil {
			logger.Error("event processing", "message", trimData(event.RawData), "err", err)
		}
		return
	}
}

func trimData(data []byte) []rune {
	trimmed := []rune(convert.BytesToString(data))
	if len(trimmed) > maxRequestLength {
		trimmed = append(trimmed[:maxRequestLength], []rune("...")...)
	}
	return trimmed
}
