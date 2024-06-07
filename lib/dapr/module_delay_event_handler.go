package dapr

import (
	"context"
	"github.com/cenkalti/backoff/v4"
	"github.com/hdget/hdsdk/v2/intf"
	"time"
)

type delayEventHandler interface {
	GetTopic() string
	Handle(ctx context.Context, logger intf.LoggerProvider, subscriber intf.Subscriber)
}

type delayEventHandlerImpl struct {
	module DelayEventModule
	topic  string             // 订阅主题
	fn     DelayEventFunction // 调用函数
}

type DelayEventFunction func(message []byte) (retry bool, err error)

func (h delayEventHandlerImpl) GetTopic() string {
	return h.topic
}

// Handle
// err: nil 只要错误为空，则消息成功消费, 不管retry的值为什么样
// err: not nil + retry: false 打印DROP status消息
// err: not nil + retry: true  进行重试，最后重试次数结束, 打印日志
func (h delayEventHandlerImpl) Handle(ctx context.Context, logger intf.LoggerProvider, subscriber intf.Subscriber) {
	msgChan, err := subscriber.SubscribeDelay(context.Background(), h.GetTopic())
	if err != nil {
		logger.Fatal("subscribe delay event", "topic", h.GetTopic())
	}

LOOP:
	for {
		select {
		case <-ctx.Done():
			logger.Debug("gracefully shutdown delay event handler", "topic", h.GetTopic())
			break LOOP
		case msg := <-msgChan:
			retry, err := h.fn(msg.Payload)
			if err == nil {
				msg.Ack()
			} else {
				if !retry { // err != nil && retry == false
					logger.Error("drop delay event", "err", err, "msg", trimData(msg.Payload))
					msg.Ack()
				} else { // err != nil && retry == true
					nextBackOff := h.module.GetBackOffPolicy().NextBackOff()
					if nextBackOff == backoff.Stop {
						logger.Error("drop delay event after retried many times")
						msg.Ack()
					} else {
						time.Sleep(nextBackOff)
						logger.Error("retry delay event", "err", err, "msg", trimData(msg.Payload))
						msg.Nack()
					}
				}
			}
		}
	}
}
