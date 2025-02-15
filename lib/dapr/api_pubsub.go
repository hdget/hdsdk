package dapr

import (
	"github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	"github.com/pkg/errors"
)

type daprEvent struct {
	subscription *common.Subscription
	handler      common.TopicEventHandler
}

// Publish 发布消息
// isRawPayLoad 发送原始的消息，非cloudevent message
func (a apiImpl) Publish(pubSubName, topic string, data interface{}, args ...bool) error {
	daprClient, err := client.NewClient()
	if err != nil {
		return errors.Wrap(err, "new dapr client")
	}
	if daprClient == nil {
		return errors.New("dapr client is null, name resolution service may not started, please check it")
	}

	// IMPORTANT: daprClient是全局的连接, 不能关闭
	//defer daprClient.Close()

	var opt client.PublishEventOption
	metaOptions := getPublishMetaOptions(args...)
	if metaOptions != nil {
		opt = client.PublishEventWithMetadata(metaOptions)
		err = daprClient.PublishEvent(a.ctx, a.normalize(pubSubName), topic, data, opt)
	} else {
		err = daprClient.PublishEvent(a.ctx, a.normalize(pubSubName), topic, data)
	}

	if err != nil {
		return err
	}

	return nil
}

func getDaprEvent(pubsubName, topic string, handler common.TopicEventHandler, args ...bool) daprEvent {
	metaOptions := getPublishMetaOptions(args...)
	return daprEvent{
		subscription: &common.Subscription{
			PubsubName: pubsubName,
			Topic:      topic,
			Metadata:   metaOptions,
		},
		handler: handler,
	}
}

func getPublishMetaOptions(args ...bool) map[string]string {
	isRawPayLoad := false
	if len(args) > 0 {
		isRawPayLoad = args[0]
	}

	if isRawPayLoad {
		return map[string]string{"rawPayload": "true"}
	}
	return nil
}
