package dapr

import (
	"context"
	"encoding/json"
	"github.com/dapr/go-sdk/client"
	"github.com/hdget/hdsdk/utils"
	"github.com/pkg/errors"
)

// SaveState 保存状态
func SaveState(storeName, key string, data interface{}) error {
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

	err = daprClient.SaveState(context.Background(), storeName, key, value)
	if err != nil {
		return errors.Wrapf(err, "save state, store: %s, key: %s, value: %s", storeName, key, value)
	}

	return nil
}

// GetState 保存状态
func GetState(storeName, key string) ([]byte, error) {
	daprClient, err := client.NewClient()
	if err != nil {
		return nil, errors.Wrap(err, "new dapr client")
	}
	if daprClient == nil {
		return nil, errors.New("dapr client is null, name resolution service may not started, please check it")
	}

	// IMPORTANT: daprClient是全局的连接, 不能关闭
	//defer daprClient.Close()
	item, err := daprClient.GetState(context.Background(), storeName, key)
	if err != nil {
		return nil, errors.Wrapf(err, "get state, store: %s, key: %s", storeName, key)
	}

	return item.Value, nil
}

// DeleteState 删除状态
func DeleteState(storeName, key string) error {
	daprClient, err := client.NewClient()
	if err != nil {
		return errors.Wrap(err, "new dapr client")
	}
	if daprClient == nil {
		return errors.New("dapr client is null, name resolution service may not started, please check it")
	}

	// IMPORTANT: daprClient是全局的连接, 不能关闭
	//defer daprClient.Close()
	err = daprClient.DeleteState(context.Background(), storeName, key)
	if err != nil {
		return errors.Wrapf(err, "delete state, store: %s, key: %s", storeName, key)
	}

	return nil
}
