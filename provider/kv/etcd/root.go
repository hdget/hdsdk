// Package etcd
// @Title  log capability of zerolog
// @Description  zerolog implementation of log capability
// @Author  Ryan Fan 2021-06-09
// @Update  Ryan Fan 2021-06-09
package etcd

import (
	"github.com/hdget/hdsdk/types"
	kv "github.com/sagikazarmark/crypt/config"
)

type EtcdProvider struct {
	kv.ConfigManager
}

var (
	_ types.Provider   = (*EtcdProvider)(nil)
	_ types.KvProvider = (*EtcdProvider)(nil)
)

// Init	implements types.Provider interface, used to initialize the capability
func (ep *EtcdProvider) Init(rootConfiger types.Configer, logger types.LogProvider, args ...interface{}) error {
	// 获取数据库配置信息
	c, err := parseConfig(rootConfiger)
	if err != nil {
		return err
	}

	err = validateConf(c)
	if err != nil {
		return err
	}

	ep.ConfigManager, err = kv.NewStandardEtcdV3ConfigManager([]string{c.Url})
	if err != nil {
		return err
	}

	return nil
}
