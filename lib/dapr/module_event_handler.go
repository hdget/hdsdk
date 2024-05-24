package dapr

import (
	"context"
	"github.com/dapr/go-sdk/service/common"
	"github.com/hdget/hdsdk/v2"
	"github.com/hdget/hdutils"
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

type handleResult struct {
	retry bool
	err   error
}

type EventFunction func(ctx context.Context, event *common.TopicEvent) (retry bool, err error)

func (h eventHandlerImpl) GetTopic() string {
	return h.topic
}

func (h eventHandlerImpl) GetEventFunction() common.TopicEventHandler {
	return func(ctx context.Context, event *common.TopicEvent) (bool, error) {
		// 在go routine中执行具体的函数
		quit := make(chan *handleResult, 1)
		go func(quit chan *handleResult) {
			var retry bool
			var err error
			// 挂载defer函数
			defer func() {
				if r := recover(); r != nil {
					hdutils.RecordErrorStack(h.module.GetApp())
				}
				quit <- &handleResult{
					retry: retry,
					err:   err,
				}
			}()

			// 执行具体的函数
			retry, err = h.fn(ctx, event)
		}(quit)

		select {
		case <-time.After(h.module.GetAckTimeout() - 1*time.Minute):
			hdsdk.Logger().Error("event processing timeout", "module", h.module.GetMeta().StructName, "topic", h.GetTopic(), "handler", hdutils.Reflect().GetFuncName(h.fn), "message", trimData(event.RawData))
			return false, nil // 丢弃消息不重试
		case ret := <-quit:
			if ret.err != nil {
				hdsdk.Logger().Error("event processing done", "module", h.module.GetMeta().StructName, "topic", h.GetTopic(), "handler", hdutils.Reflect().GetFuncName(h.fn), "message", trimData(event.RawData), "err", ret.err)
			}
			return ret.retry, ret.err
		}
	}
}

func trimData(data []byte) []rune {
	trimmed := []rune(hdutils.BytesToString(data))
	if len(trimmed) > maxRequestLength {
		trimmed = append(trimmed[:maxRequestLength], []rune("...")...)
	}
	return trimmed
}
