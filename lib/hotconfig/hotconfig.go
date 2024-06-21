package hotconfig

import (
	"encoding/json"
	"github.com/friendsofgo/errors"
	jsonUtils "github.com/hdget/hdutils/json"
)

type HotConfig interface {
	GetName() string               // 获取名字
	GetValue() (any, error)        // 获取配置
	UpdateValue(data []byte) error // 更新值
}

type Manager interface {
	SaveConfig(name string, data []byte) error
	LoadConfig(configName string) ([]byte, error)
	GetInstance(configName string) HotConfig
	Register(configName string, defaultConfigValue any)
}

type hotConfigObject struct {
	manager  Manager
	name     string
	isLoaded bool
	value    any
}

func (h *hotConfigObject) GetName() string {
	return h.name
}

func (h *hotConfigObject) GetValue() (any, error) {
	if h.isLoaded {
		return h.value, nil
	}

	data, err := h.manager.LoadConfig(h.name)
	if err != nil {
		return nil, errors.Wrap(err, "load hot config")
	}

	err = h.UpdateValue(data)
	if err != nil {
		return nil, errors.Wrap(err, "update hot config value")
	}

	h.isLoaded = true
	return h.value, nil
}

func (h *hotConfigObject) UpdateValue(data []byte) error {
	if jsonUtils.IsEmptyJsonObject(data) {
		return errors.New("empty hot config data")
	}
	return json.Unmarshal(data, &h.value)
}
