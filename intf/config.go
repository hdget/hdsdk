package intf

type Configer interface {
	GetLogConfig() map[string]any      // 日志配置
	GetMysqlConfig() map[string]any    // mysql数据库配置
	GetRedisConfig() map[string]any    // redis缓存配置
	GetRabbitmqConfig() map[string]any // GetRabbitmqConfig 获取RabbitMq队列配置
	GetEtcdConfig() map[string]any     // kv型服务配置
	GetNeo4jConfig() map[string]any    // 图数据库配置
}

type ConfigLoader interface {
	Load(configVar any) error
	LoadKey(key string, configVar any) error
}
