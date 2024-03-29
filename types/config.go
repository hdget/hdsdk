package types

type Configer interface {
	GetLogConfig() interface{}      // 日志配置
	GetMysqlConfig() interface{}    // mysql数据库配置
	GetRedisConfig() interface{}    // redis缓存配置
	GetRabbitmqConfig() interface{} // GetRabbitmqConfig 获取RabbitMq队列配置
	GetKafkaConfig() interface{}    // GetKafkaConfig 获取kafka消息队列配置
	GetNosqlConfig() interface{}    // NoSQL服务配置
	GetEtcdConfig() interface{}     // kv型服务配置
	GetGraphConfig() interface{}    // 图数据库配置
}

// SdkConfigItem items under sdk config
type SdkConfigItem struct {
	Log      interface{} `mapstructure:"log"`      // 日志配置
	Mysql    interface{} `mapstructure:"mysql"`    // 数据库配置
	Redis    interface{} `mapstructure:"redis"`    // 缓存配置
	RabbitMq interface{} `mapstructure:"rabbitmq"` // rabbitmq消息队列配置
	Kafka    interface{} `mapstructure:"kafka"`    // kafka消息队列配置
	Nosql    interface{} `mapstructure:"nosql"`    // NoSQL配置
	Etcd     interface{} `mapstructure:"etcd"`     // etcd Key/Value数据库配置
	Neo4j    interface{} `mapstructure:"neo4j"`    // neo4j数据库配置
}
