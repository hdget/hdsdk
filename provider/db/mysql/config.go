package mysql

//
//type ConfigMysql struct {
//	Default *MySqlConf   `mapstructure:"default"`
//	Master  *MySqlConf   `mapstructure:"master"`
//	Slaves  []*MySqlConf `mapstructure:"slaves"`
//	Items   []*MySqlConf `mapstructure:"items"`
//}
//
//type MySqlConf struct {
//	Name     string `mapstructure:"name"`
//	User     string `mapstructure:"user"`
//	Password string `mapstructure:"password"`
//	Host     string `mapstructure:"host"`
//	Port     int    `mapstructure:"port"`
//	Database string `mapstructure:"database"`
//	Timeout  int    `mapstructure:"timeout"`
//}
//
//// /////////////////////////////////////////////////////////////////
//func parseConfig(rootConfiger intf.Configer) (*ConfigMysql, error) {
//	if rootConfiger == nil {
//		return nil, intf.ErrEmptyConfig
//	}
//
//	data := rootConfiger.GetMysqlConfig()
//	if data == nil {
//		return nil, intf.ErrEmptyConfig
//	}
//
//	values, ok := data.(map[string]interface{})
//	if !ok {
//		return nil, intf.ErrInvalidConfig
//	}
//
//	var conf ConfigMysql
//	err := mapstructure.Decode(values, &conf)
//	if err != nil {
//		return nil, errors.Wrap(err, "decode db content")
//	}
//
//	return &conf, nil
//}
//
//func validateConf(providerType string, conf *MySqlConf) error {
//	if conf == nil {
//		return intf.ErrEmptyConfig
//	}
//
//	if conf.Host == "" || conf.Database == "" || conf.User == "" {
//		return intf.ErrInvalidConfig
//	}
//
//	if providerType == intf.ProviderTypeOther && conf.Name == "" {
//		return intf.ErrInvalidConfig
//	}
//
//	return nil
//}
