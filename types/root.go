package types

// 底层库能力提供者接口
type Provider interface {
	Init(rootConfiger Configer, logger LogProvider, args ...interface{}) error // 初始化底层能力
}

// sdk底层能力分类
// sdk底层能力分为了几大类：
// - SdkCategoryDb 数据库底层能力
// - SdkCategoryCache 缓存底层能力
// - SdkCategoryMq 底层消息队列能力
// - SdkCategoryNosql 非典型SQL数据库底层能力

// 底层库能力类别
const SdkCategoryOffset = 10

type SdkType int

const (
	_                SdkType = SdkCategoryOffset * iota
	SdkCategoryLog           // 日志能力
	SdkCategoryDb            // 数据库能力, 例如mysql
	SdkCategoryCache         // 缓存能力，例如redis
	SdkCategoryMq            // 消息队列能力，例如rabbitmq, rocketmq, kafka
	SdkCategoryMs            // 微服务能力
	SdkCategoryNosql         // nosql数据库能力，例如es, monodb
	SdkCategoryKv            // kv型数据库能力，例如etcd
)

const (
	PROVIDER_TYPE_DEFAULT = "default"
	PROVIDER_TYPE_MASTER  = "master"
	PROVIDER_TYPE_SLAVE   = "slave"
	PROVIDER_TYPE_OTHER   = "other"
)
