package neo4j

//
//type ConfigNeo4j struct {
//	VirtualUri  string             `mapstructure:"virtual_uri"`
//	Username    string             `mapstructure:"username"`
//	Password    string             `mapstructure:"password"`
//	Servers     []*Neo4jServerConf `mapstructure:"servers"`
//	MaxPoolSize int                `mapstructure:"max_pool_size"`
//}
//
//type Neo4jServerConf struct {
//	Host string `mapstructure:"host"`
//	Port int    `mapstructure:"port"`
//}
//
//// /////////////////////////////////////////////////////////////////
//func parseConfig(rootConfiger intf.Configer) (*ConfigNeo4j, error) {
//	if rootConfiger == nil {
//		return nil, intf.ErrEmptyConfig
//	}
//
//	data := rootConfiger.GetNeo4jConfig()
//	if data == nil {
//		return nil, intf.ErrEmptyConfig
//	}
//
//	values, ok := data.(map[string]interface{})
//	if !ok {
//		return nil, intf.ErrInvalidConfig
//	}
//
//	var conf ConfigNeo4j
//	err := mapstructure.Decode(values, &conf)
//	if err != nil {
//		return nil, errors.Wrap(err, "decode neo4j configer")
//	}
//
//	return &conf, nil
//}
//
//func validateConf(providerType string, conf *ConfigNeo4j) error {
//	if conf == nil {
//		return intf.ErrEmptyConfig
//	}
//
//	if conf.VirtualUri == "" || conf.Username == "" || conf.Password == "" {
//		return intf.ErrInvalidConfig
//	}
//
//	for _, server := range conf.Servers {
//		if server.Host == "" || server.Port == 0 {
//			return intf.ErrInvalidConfig
//		}
//	}
//
//	return nil
//}
