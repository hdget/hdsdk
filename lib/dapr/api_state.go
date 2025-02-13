package dapr

import (
	"fmt"
	"github.com/dapr/go-sdk/client"
	"github.com/hdget/hdutils/convert"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// SaveState 保存状态
func (a apiImpl) SaveState(storeName, key string, value interface{}) error {
	data, err := convert.ToBytes(value)
	if err != nil {
		return err
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

	err = daprClient.SaveState(a.ctx, a.standardize(storeName), key, data, nil)
	if err != nil {
		return errors.Wrapf(err, "save state, store: %s, key: %s, value: %s", storeName, key, value)
	}

	return nil
}

// GetState 获取状态
func (a apiImpl) GetState(storeName, key string) ([]byte, error) {
	daprClient, err := client.NewClient()
	if err != nil {
		return nil, errors.Wrap(err, "new dapr client")
	}
	if daprClient == nil {
		return nil, errors.New("dapr client is null, name resolution service may not started, please check it")
	}

	// IMPORTANT: daprClient是全局的连接, 不能关闭
	item, err := daprClient.GetState(a.ctx, a.standardize(storeName), key, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "get state, store: %s, key: %s", storeName, key)
	}

	return item.Value, nil
}

// GetBulkState 批量获取状态
func (a apiImpl) GetBulkState(storeName string, keys any) (map[string][]byte, error) {
	strKeys, err := cast.ToStringSliceE(keys)
	if err != nil {
		return nil, fmt.Errorf("invalid keys, keys: %v", keys)
	}

	daprClient, err := client.NewClient()
	if err != nil {
		return nil, errors.Wrap(err, "new dapr client")
	}
	if daprClient == nil {
		return nil, errors.New("dapr client is null, name resolution service may not started, please check it")
	}

	// IMPORTANT: daprClient是全局的连接, 不能关闭
	items, err := daprClient.GetBulkState(a.ctx, a.standardize(storeName), strKeys, nil, 100)
	if err != nil {
		return nil, errors.Wrapf(err, "get bulk state, store: %s, keys: %s", storeName, keys)
	}

	results := make(map[string][]byte, len(items))
	for _, item := range items {
		if item.Error == "" {
			results[item.Key] = item.Value
		}
	}
	return results, nil
}

// DeleteState 删除状态
func (a apiImpl) DeleteState(storeName, key string) error {
	daprClient, err := client.NewClient()
	if err != nil {
		return errors.Wrap(err, "new dapr client")
	}
	if daprClient == nil {
		return errors.New("dapr client is null, name resolution service may not started, please check it")
	}

	// IMPORTANT: daprClient是全局的连接, 不能关闭
	//defer daprClient.Close()
	err = daprClient.DeleteState(a.ctx, a.standardize(storeName), key, nil)
	if err != nil {
		return errors.Wrapf(err, "delete state, store: %s, key: %s", storeName, key)
	}

	return nil
}
