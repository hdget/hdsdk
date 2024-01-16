// Package hdsdk
// 提供各类底层能力的直接访问方式，SDK包在使用前必须要初始化
//
// 首先必须创建一个继承自sdk.BaseConfig的配置struct
// e,g:
//
//	   import hdget
//
//			type XXXConfig struct {
//				*sdk.Config `mapstructure:",squash"`
//	     }
package hdsdk

//
//import (
//	"github.com/hdget/hdsdk/intf"
//	"github.com/hdget/hdsdk/provider/db/mysql"
//	"github.com/hdget/hdsdk/provider/etcd"
//	"github.com/hdget/hdsdk/provider/logger"
//	"github.com/hdget/hdsdk/provider/mq/kafka"
//	"github.com/hdget/hdsdk/provider/mq/rabbitmq_deprecated"
//	"github.com/hdget/hdsdk/provider/neo4j"
//	"github.com/hdget/hdsdk/provider/redis"
//	"github.com/pkg/errors"
//)
//
//// SdkProvider 底层能力实例
//type SdkProvider struct {
//	Kind     intf.SdkType  // 底层能力类型
//	Name     string        // 底层能力名字
//	Instance intf.Provider // 底层能力实现实例
//}
//
//var (
//	Logger   intf.LogProvider   // 日志能力
//	Mysql    intf.DbProvider    // mysql数据库能力
//	Redis    intf.CacheProvider // redis缓存能力
//	Rabbitmq intf.MqProvider    // rabbitmq能力
//	Kafka    intf.MqProvider    // kafka能力
//	Neo4j    intf.GraphProvider // 图数据库能力
//	Etcd     intf.KvProvider    // kv能力
//)
//
//var (
//	LogProvider = &SdkProvider{
//		Kind:     intf.SdkCategoryLog,
//		Name:     "logger",
//		Instance: &logger.LoggerImpl{},
//	}
//
//	// 除去日志外其他能力能力提供者实例
//	providers = []*SdkProvider{
//		{
//			Kind:     intf.SdkTypeDbMysql,
//			Name:     "mysql",
//			Instance: &mysql.MysqlProvider{},
//		},
//		{
//			Kind:     intf.SdkTypeCacheRedis,
//			Name:     "redis",
//			Instance: &redis.RedisProvider{},
//		},
//		{
//			Kind:     intf.SdkTypeMqRabbitmq,
//			Name:     "rabbitmq",
//			Instance: &rabbitmq_deprecated.RabbitmqProvider{},
//		},
//		{
//			Kind:     intf.SdkTypeMqKafka,
//			Name:     "kafka",
//			Instance: &kafka.KafkaProvider{},
//		},
//		{
//			Kind:     intf.SdkTypeGraphNeo4j,
//			Name:     "neo4j",
//			Instance: &neo4j.Neo4jProvider{},
//		},
//		{
//			Kind:     intf.SdkTypeKvEtcd,
//			Name:     "etcd",
//			Instance: &etcd.EtcdProvider{},
//		},
//	}
//)
//
//// Initialize 初始化SDK, 指定的配置文件里面有什么配置就配置什么能力
//func Initialize(configer intf.Configer) error {
//	var err error
//	Logger, err = newLogger(configer)
//	if err != nil {
//		return err
//	}
//
//	for _, p := range providers {
//		err = p.Instance.Init(configer, Logger)
//		// 如果没有对应能力的配置，忽略该底层能力的后续初始化动作
//		if errors.Is(err, intf.ErrEmptyConfig) {
//			continue
//		}
//
//		// 打印提示日志
//		if err != nil {
//			Logger.Error("initialize provider", "name", p.Name, "err", err)
//		} else {
//			Logger.Info("initialize provider", "name", p.Name)
//		}
//
//		setGlobalVars(p)
//	}
//
//	return nil
//}
//
//func setGlobalVars(p *SdkProvider) {
//	// 根据不同的能力类型，将provider Instance转换成具体的provider
//	switch p.Kind {
//	case intf.SdkTypeDbMysql:
//		Mysql = p.Instance.(*mysql.MysqlProvider)
//	case intf.SdkTypeCacheRedis:
//		Redis = p.Instance.(*redis.RedisProvider)
//	case intf.SdkTypeMqRabbitmq:
//		Rabbitmq = p.Instance.(*rabbitmq_deprecated.RabbitmqProvider)
//	case intf.SdkTypeMqKafka:
//		Kafka = p.Instance.(*kafka.KafkaProvider)
//	case intf.SdkTypeGraphNeo4j:
//		Neo4j = p.Instance.(*neo4j.Neo4jProvider)
//	case intf.SdkTypeKvEtcd:
//		Etcd = p.Instance.(*etcd.EtcdProvider)
//	}
//}
//
//// 初始化日志服务
//func newLogger(configer intf.Configer) (intf.LogProvider, error) {
//	err := LogProvider.Instance.Init(configer, nil)
//	if err != nil {
//		return nil, err
//	}
//
//	logger, ok := LogProvider.Instance.(*logger.LoggerImpl)
//	if !ok {
//		return nil, errors.New("error convert to LoggerImpl")
//	}
//
//	return logger, nil
//}
