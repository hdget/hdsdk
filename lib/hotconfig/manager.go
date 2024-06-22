package hotconfig

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dapr/go-sdk/client"
	"github.com/elliotchance/pie/v2"
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/hdget/hdsdk/v2/lib/dapr"
	"github.com/hdget/hdutils/convert"
	"github.com/hdget/hdutils/logger"
	"github.com/pkg/errors"
)

type Transactor interface {
	Commit() error
	Rollback() error
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type hotConfigManager struct {
	app             string
	registry        map[string]HotConfig
	subscribed      bool // 是否订阅了动态配置的变化
	saveFunction    SaveFunction
	loadFunction    LoadFunction
	daprConfigStore string
	redisClient     intf.RedisClient
}

var (
	_managerInstance Manager
)

func NewManager(app string, options ...Option) Manager {
	if _managerInstance == nil {
		v := &hotConfigManager{
			app:        app,
			registry:   make(map[string]HotConfig),
			subscribed: false,
		}

		for _, option := range options {
			option(v)
		}
		_managerInstance = v
	}
	return _managerInstance
}

func (impl *hotConfigManager) LoadConfig(configName string) ([]byte, error) {
	if impl.loadFunction == nil {
		return nil, errors.Errorf("load function not specified")
	}

	// 订阅配置变化
	err := impl.subscribeConfigChanges()
	if err != nil {
		return nil, errors.Wrap(err, "subscribe config changes")
	}

	return impl.loadFunction(configName)

}

func (impl *hotConfigManager) SaveConfig(configName string, data []byte) error {
	if !pie.Contains(pie.Keys(impl.registry), configName) {
		return fmt.Errorf("invalid config name, name: %s", configName)
	}

	if impl.saveFunction == nil {
		return errors.New("save function not specified")
	}

	if impl.redisClient == nil {
		return errors.New("redis client not initialized")
	}

	if impl.daprConfigStore == "" {
		return errors.New("config store not found")
	}

	// 先执行persistent函数，比如存入数据库
	tx, err := impl.saveFunction(configName, data)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
		_ = tx.Commit()
	}()

	// 保存到数据库的同时，写入到缓存中
	err = impl.redisClient.Set(impl.getConfigKey(configName), data)
	if err != nil {
		return err
	}

	return nil
}

func (impl *hotConfigManager) GetInstance(name string) HotConfig {
	return impl.registry[name]
}

func (impl *hotConfigManager) Register(configName string, defaultConfigValue any) {
	impl.registry[configName] = &hotConfigObject{
		manager: impl,
		name:    configName,
		value:   defaultConfigValue,
	}
}

func (impl *hotConfigManager) subscribeConfigChanges() error {
	if impl.subscribed {
		return nil
	}

	configKey2name := make(map[string]string)
	for k := range impl.registry {
		configKey2name[impl.getConfigKey(k)] = k
	}

	if len(configKey2name) == 0 {
		return nil
	}

	subscriberId, err := dapr.Api().SubscribeConfigurationItems(context.Background(), impl.daprConfigStore, pie.Keys(configKey2name), func(id string, items map[string]*client.ConfigurationItem) {
		for configKey, configItem := range items {
			instance := impl.GetInstance(configKey2name[configKey])
			if instance != nil {
				err := instance.UpdateValue(convert.StringToBytes(configItem.Value))
				if err != nil {
					logger.Error("update value", "configKey", configKey, "value", configItem.Value, "err", err)
				}
			}
		}
	})
	if err != nil {
		return errors.Wrap(err, "subscribe hot config changes")
	}

	logger.Debug("subscribe hot config changes", "subscriberId:", subscriberId)
	impl.subscribed = true
	return nil
}

func (impl *hotConfigManager) getConfigKey(configName string) string {
	return fmt.Sprintf("hotconfig:%s:%s", impl.app, configName)
}
