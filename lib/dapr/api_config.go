package dapr

import (
	"context"
	"github.com/dapr/go-sdk/client"
	"github.com/pkg/errors"
)

// GetConfigurationItems 获取配置项
func (a apiImpl) GetConfigurationItems(configStore string, keys []string) (map[string]*client.ConfigurationItem, error) {
	daprClient, err := client.NewClient()
	if err != nil {
		return nil, errors.Wrap(err, "new dapr client")
	}
	if daprClient == nil {
		return nil, errors.New("dapr client is null, name resolution service may not started, please check it")
	}

	items, err := daprClient.GetConfigurationItems(a.ctx, a.normalize(configStore), keys)
	if err != nil {
		return nil, errors.Wrap(err, "get configuration items")
	}

	return items, nil
}

// SubscribeConfigurationItems 订阅配置项更改
func (a apiImpl) SubscribeConfigurationItems(ctx context.Context, configStore string, keys []string, handler client.ConfigurationHandleFunction) (string, error) {
	daprClient, err := client.NewClient()
	if err != nil {
		return "", errors.Wrap(err, "new dapr client")
	}
	if daprClient == nil {
		return "", errors.New("dapr client is null, name resolution service may not started, please check it")
	}

	subscriberId, err := daprClient.SubscribeConfigurationItems(ctx, a.normalize(configStore), keys, handler)
	if err != nil {
		return "", errors.Wrap(err, "subscribe configuration items update")
	}
	return subscriberId, nil
}
