package neo4j

//
//type neo4jProviderConfig struct {
//	VirtualUri  string             `mapstructure:"virtual_uri"`
//	Username    string             `mapstructure:"username"`
//	Password    string             `mapstructure:"password"`
//	Servers     []*neo4jServerConf `mapstructure:"servers"`
//	MaxPoolSize int                `mapstructure:"max_pool_size"`
//}
//
//type neo4jServerConf struct {
//	Host string `mapstructure:"host"`
//	Port int    `mapstructure:"port"`
//}
//
//const (
//	defaultMaxPoolSize = 500
//	configSection      = "sdk.neo4j"
//)
//
//func NewConfig(configLoader intf.ConfigLoader) (*neo4jProviderConfig, error) {
//	if configLoader == nil {
//		return nil, errdef.ErrEmptyConfig
//	}
//
//	var c *neo4jProviderConfig
//	err := configLoader.Unmarshal(&c, configSection)
//	if err != nil {
//		return nil, err
//	}
//
//	if c == nil {
//		return nil, errdef.ErrEmptyConfig
//	}
//
//	err = c.validate()
//	if err != nil {
//		return nil, errors.Wrap(err, "validate neo4j config")
//	}
//
//	return c, nil
//}
//
//func (c *neo4jProviderConfig) validate() error {
//	if c.VirtualUri == "" || c.Username == "" || c.Password == "" {
//		return errdef.ErrInvalidConfig
//	}
//
//	for _, server := range c.Servers {
//		if server.Host == "" || server.Port == 0 {
//			return errdef.ErrInvalidConfig
//		}
//	}
//
//	// setup default config items
//	if c.MaxPoolSize == 0 {
//		c.MaxPoolSize = defaultMaxPoolSize
//	}
//
//	return nil
//}
