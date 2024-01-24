package etcd

//
//type ConfigEtcd struct {
//	Url string `mapstructure:"url"`
//}
//
//// /////////////////////////////////////////////////////////////////
//func parseConfig(rootConfiger intf.Configer) (*ConfigEtcd, error) {
//	if rootConfiger == nil {
//		return nil, intf.ErrEmptyConfig
//	}
//
//	data := rootConfiger.GetEtcdConfig()
//	if data == nil {
//		return nil, intf.ErrEmptyConfig
//	}
//
//	values, ok := data.(map[string]interface{})
//	if !ok {
//		return nil, intf.ErrInvalidConfig
//	}
//
//	var conf ConfigEtcd
//	err := mapstructure.Decode(values, &conf)
//	if err != nil {
//		return nil, errors.Wrap(err, "decode db content")
//	}
//
//	return &conf, nil
//}
//
//func validateConf(conf *ConfigEtcd) error {
//	if conf == nil {
//		return intf.ErrEmptyConfig
//	}
//
//	if conf.Url == "" {
//		return intf.ErrInvalidConfig
//	}
//
//	return nil
//}
