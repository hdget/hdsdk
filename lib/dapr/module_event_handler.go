package dapr

import (
	"context"
	"fmt"
	"github.com/dapr/go-sdk/service/common"
	"github.com/hdget/hdsdk/v2/intf"
	panicUtils "github.com/hdget/hdutils/panic"
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
		ctxWithTimeout, cancel := context.WithTimeout(ctx, h.module.GetAckTimeout())
		defer cancel() // 重要：释放资源

		quit := make(chan *eventHandleResult, 1)
		go func() {
			fnResult := &eventHandleResult{
				retry: false,
				err:   nil,
			}

			defer func() {
				if r := recover(); r != nil {
					fnResult.err = fmt.Errorf("panic: %v", r)
					panicUtils.RecordErrorStack(h.module.GetApp())
				}

				// 传递执行结果
				quit <- fnResult
			}()

			// 执行具体的函数
			fnResult.retry, fnResult.err = h.fn(ctxWithTimeout, event)
		}()

		select {
		case <-ctxWithTimeout.Done(): // 统一用context控制
			logger.Error("event processing timeout, discard message", "data", truncate(event.RawData))
			return false, ctxWithTimeout.Err()
		case quitResult := <-quit:
			if quitResult.err != nil {
				logger.Error("event processing", "data", truncate(event.RawData), "err", quitResult.err)
			}
			return quitResult.retry, quitResult.err
		}
	}
}
