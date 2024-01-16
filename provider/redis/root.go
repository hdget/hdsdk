// Package redis
// @Title  logger capability of zerolog
// @Description  zerolog implementation of logger capability
// @Author  Ryan Fan 2021-06-09
// @Update  Ryan Fan 2021-06-09
package redis

import (
	_ "github.com/go-sql-driver/mysql"
)

//
//type RedisProvider struct {
//	BaseCacheProvider
//	Log logger.LogProvider
//}
//
//var (
//	_ intf.Provider      = (*RedisProvider)(nil)
//	_ intf.CacheProvider = (*RedisProvider)(nil)
//)
//
//// Init implements intf.Provider interface, used to initialize the capability
//// @author	Ryan Fan	(2021-06-09)
//// @param	baseconf.Configer	root configer interface to extract configer info
//// @return	error
//func (rp *RedisProvider) Init(rootConfiger intf.Configer, logger logger.LogProvider, _ ...interface{}) error {
//	// 获取日志配置信息
//	configloader, err := parseConfig(rootConfiger)
//	if err != nil {
//		return err
//	}
//
//	// 缺省redis必须要配置合法
//	err = validateConf(intf.ProviderTypeDefault, configloader.Default)
//	if err != nil {
//		logger.Fatal("validate redis configer", "type", intf.ProviderTypeDefault, "err", err)
//	}
//
//	rp.Default, err = rp.connect(configloader.Default)
//	if err != nil {
//		logger.Fatal("connect redis", "type", intf.ProviderTypeDefault, "host", configloader.Default.Host, "err", err)
//	}
//	logger.Debug("connect redis", "type", intf.ProviderTypeDefault, "host", configloader.Default.Host)
//
//	// 额外的redis
//	rp.Items = make(map[string]intf.CacheClient)
//	for _, otherConf := range configloader.Items {
//		if err := validateConf(intf.ProviderTypeOther, otherConf); err == nil {
//			instance, err := rp.connect(otherConf)
//			if instance != nil {
//				rp.Items[otherConf.Name] = instance
//			}
//			logger.Debug("connect redis", "type", otherConf.Name, "host", otherConf.Host, "err", err)
//		}
//	}
//
//	return nil
//}
//
//func (rp *RedisProvider) connect(conf *RedisConf) (intf.CacheClient, error) {
//	client := NewRedisClient(conf)
//	err := client.Ping()
//	if err != nil {
//		return nil, err
//	}
//	return client, nil
//}
