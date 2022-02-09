package dapr

import (
	"context"
	"encoding/json"
	"github.com/dapr/go-sdk/client"
	"github.com/hdget/hdsdk/utils"
	"github.com/pkg/errors"
)

// Publish 发布消息
func Publish(pubSubName, topic string, data interface{}) error {
	var value []byte
	switch t := data.(type) {
	case string:
		value = utils.StringToBytes(t)
	case []byte:
		value = t
	default:
		v, err := json.Marshal(data)
		if err != nil {
			return errors.Wrap(err, "marshal invoke data")
		}
		value = v
	}

	daprClient, err := client.NewClient()
	if err != nil {
		return errors.Wrap(err, "new dapr client")
	}
	if daprClient == nil {
		return errors.New("dapr client is null, name resolution service may not started, please check it")
	}

	// IMPORTANT: daprClient是全局的连接, 不能关闭
	//defer daprClient.Close()

	err = daprClient.PublishEvent(context.Background(), pubSubName, topic, value)
	if err != nil {
		return errors.Wrapf(err, "publish event, pubsub: %s, topic: %s, value: %s", pubSubName, topic, value)
	}

	return nil
}
