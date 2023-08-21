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

import (
	"github.com/hdget/hdsdk/provider/cache/redis"
	"github.com/hdget/hdsdk/provider/db/mysql"
	"github.com/hdget/hdsdk/provider/graph/neo4j"
	"github.com/hdget/hdsdk/provider/log"
	"github.com/hdget/hdsdk/provider/mq/kafka"
	"github.com/hdget/hdsdk/provider/mq/rabbitmq"
	"github.com/hdget/hdsdk/types"
	"github.com/pkg/errors"
)

// SdkProvider 底层能力实例
type SdkProvider struct {
	Kind     types.SdkType  // 底层能力类型
	Name     string         // 底层能力名字
	Instance types.Provider // 底层能力实现实例
}

var (
	Logger   types.LogProvider   // 日志能力
	Mysql    types.DbProvider    // mysql数据库能力
	Redis    types.CacheProvider // redis缓存能力
	Rabbitmq types.MqProvider    // rabbitmq能力
	Kafka    types.MqProvider    // kafka能力
	Neo4j    types.GraphProvider // 图数据库能力
)

var (
	LogProvider = &SdkProvider{
		Kind:     types.SdkCategoryLog,
		Name:     "log",
		Instance: &log.LoggerImpl{},
	}

	// 除去日志外其他能力能力提供者实例
	providers = []*SdkProvider{
		{
			Kind:     types.SdkTypeDbMysql,
			Name:     "mysql",
			Instance: &mysql.MysqlProvider{},
		},
		{
			Kind:     types.SdkTypeCacheRedis,
			Name:     "redis",
			Instance: &redis.RedisProvider{},
		},
		{
			Kind:     types.SdkTypeMqRabbitmq,
			Name:     "rabbitmq",
			Instance: &rabbitmq.RabbitmqProvider{},
		},
		{
			Kind:     types.SdkTypeMqKafka,
			Name:     "kafka",
			Instance: &kafka.KafkaProvider{},
		},
		{
			Kind:     types.SdkTypeGraphNeo4j,
			Name:     "neo4j",
			Instance: &neo4j.Neo4jProvider{},
		},
	}
)

// Initialize 初始化SDK, 指定的配置文件里面有什么配置就配置什么能力
func Initialize(configer types.Configer) error {
	var err error
	Logger, err = newLogger(configer)
	if err != nil {
		return err
	}

	for _, p := range providers {
		err = p.Instance.Init(configer, Logger)
		// 如果没有对应能力的配置，忽略该底层能力的后续初始化动作
		if errors.Is(err, types.ErrEmptyConfig) {
			continue
		}

		// 打印提示日志
		if err != nil {
			Logger.Error("initialize provider", "name", p.Name, "err", err)
		} else {
			Logger.Info("initialize provider", "name", p.Name)
		}

		setGlobalVars(p)
	}

	return nil
}

func setGlobalVars(p *SdkProvider) {
	// 根据不同的能力类型，将provider Instance转换成具体的provider
	switch p.Kind {
	case types.SdkTypeDbMysql:
		Mysql = p.Instance.(*mysql.MysqlProvider)
	case types.SdkTypeCacheRedis:
		Redis = p.Instance.(*redis.RedisProvider)
	case types.SdkTypeMqRabbitmq:
		Rabbitmq = p.Instance.(*rabbitmq.RabbitmqProvider)
	case types.SdkTypeMqKafka:
		Kafka = p.Instance.(*kafka.KafkaProvider)
	case types.SdkTypeGraphNeo4j:
		Neo4j = p.Instance.(*neo4j.Neo4jProvider)
	}
}

// 初始化日志服务
func newLogger(configer types.Configer) (types.LogProvider, error) {
	err := LogProvider.Instance.Init(configer, nil)
	if err != nil {
		return nil, err
	}

	logger, ok := LogProvider.Instance.(*log.LoggerImpl)
	if !ok {
		return nil, errors.New("error convert to LoggerImpl")
	}

	return logger, nil
}
