package hddapr

import (
	"context"
	"github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	"github.com/pkg/errors"
)

type Event struct {
	Subscription *common.Subscription
	Handler      common.TopicEventHandler
}

// Publish 发布消息
// isRawPayLoad 发送原始的消息，非cloudevent message
func Publish(pubSubName, topic string, data interface{}, args ...bool) error {
	daprClient, err := client.NewClient()
	if err != nil {
		return errors.Wrap(err, "new hddapr client")
	}
	if daprClient == nil {
		return errors.New("hddapr client is null, handlerName resolution service may not started, please check it")
	}

	// IMPORTANT: daprClient是全局的连接, 不能关闭
	//defer daprClient.Close()

	var opt client.PublishEventOption
	metaOptions := getMeta(args...)
	if metaOptions != nil {
		opt = client.PublishEventWithMetadata(metaOptions)
		err = daprClient.PublishEvent(context.Background(), pubSubName, topic, data, opt)
	} else {
		err = daprClient.PublishEvent(context.Background(), pubSubName, topic, data)
	}

	if err != nil {
		return err
	}

	return nil
}

func GetEvent(pubsubName, topic string, handler common.TopicEventHandler, args ...bool) Event {
	metaOptions := getMeta(args...)
	return Event{
		Subscription: &common.Subscription{
			PubsubName: pubsubName,
			Topic:      topic,
			Metadata:   metaOptions,
		},
		Handler: handler,
	}
}

func getMeta(args ...bool) map[string]string {
	isRawPayLoad := false
	if len(args) > 0 {
		isRawPayLoad = args[0]
	}

	if isRawPayLoad {
		return map[string]string{"rawPayload": "true"}
	}
	return nil
}
