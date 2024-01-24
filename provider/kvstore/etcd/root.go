// Package etcd
// @Title  logger capability of zerolog
// @Description  zerolog implementation of logger capability
// @Author  Ryan Fan 2021-06-09
// @Update  Ryan Fan 2021-06-09
package etcd

//
//type EtcdProvider struct {
//	kv.ConfigManager
//}
//
//var (
//	_ intf.Provider   = (*EtcdProvider)(nil)
//	_ intf.KvProvider = (*EtcdProvider)(nil)
//)
//
//// Init	implements intf.Provider interface, used to initialize the capability
//func (ep *EtcdProvider) Init(rootConfiger intf.Configer, logger logger.LogProvider, args ...interface{}) error {
//	// 获取数据库配置信息
//	c, err := parseConfig(rootConfiger)
//	if err != nil {
//		return err
//	}
//
//	err = validateConf(c)
//	if err != nil {
//		return err
//	}
//
//	ep.ConfigManager, err = kv.NewStandardEtcdV3ConfigManager([]string{c.Url})
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
