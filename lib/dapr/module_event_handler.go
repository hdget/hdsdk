package dapr

import (
	"context"
	"github.com/dapr/go-sdk/service/common"
	"github.com/hdget/hdsdk/v2/intf"
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

type eventHandleResult struct {
	retry bool
	err   error
}

type EventFunction func(ctx context.Context, event *common.TopicEvent) (retry bool, err error)

func (h eventHandlerImpl) GetTopic() string {
	return h.topic
}

// GetEventFunction
// err: nil 只要错误为空，则消息成功消费, 不管retry的值为什么样
// err: not nil + retry: false DAPR打印DROP status消息
// err: not nil + retry: true  根据DAPR resilience策略进行重试，最后重试次数结束, DAPR打印日志
func (h eventHandlerImpl) GetEventFunction(logger intf.LoggerProvider) common.TopicEventHandler {
	return func(ctx context.Context, event *common.TopicEvent) (bool, error) {
		quit := make(chan *eventHandleResult, 1)
		go func(chanResult chan *eventHandleResult) {
			fnResult := &eventHandleResult{}
			defer func() {
				if r := recover(); r != nil {
					panicUtils.RecordErrorStack(h.module.GetApp())
				}

				// 传递执行结果
				chanResult <- fnResult
			}()

			// 执行具体的函数
			fnResult.retry, fnResult.err = h.fn(ctx, event)
		}(quit)

		var result *eventHandleResult
		select {
		case <-time.After(h.module.GetAckTimeout()): // 超时则丢弃消息
			logger.Error("event processing timeout, discard message", "data", truncate(event.RawData))
			result = &eventHandleResult{
				retry: false,
				err:   nil,
			}
		case result = <-quit: // 如果gorouting中的函数在没超时之前退出,获取执行结果
			if result.err != nil {
				logger.Error("event processing", "data", truncate(event.RawData), "err", result.err)
			}
		}
		return result.retry, result.err
	}
}
