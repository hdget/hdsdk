package intf

import "go.uber.org/fx"

type ProviderCategory int

const (
	ProviderCategoryUnknown ProviderCategory = iota
	ProviderCategoryConfig
	ProviderCategoryLogger
	ProviderCategoryDb
	ProviderCategoryRedis
	ProviderCategoryMq
	ProviderCategoryDbSqlx
	ProviderCategoryDbBuilder
)

type ProviderName string

const (
	ProviderNameConfigViper       ProviderName = "config-viper"
	ProviderNameLoggerZerolog     ProviderName = "logger-zerolog"
	ProviderNameRedisRedigo       ProviderName = "redis-redigo"
	ProviderNameDbSqlBoilerMysql  ProviderName = "db-sqlboiler-mysql"
	ProviderNameDbSqlBoilerSqlite ProviderName = "db-sqlboiler-sqlite3"
	ProviderNameDbSqlxMysql       ProviderName = "db-sqlx-mysql"
	ProviderNameDbSquirrelMysql   ProviderName = "db-squirrel-mysql"
	ProviderNameMqRabbitMq        ProviderName = "mq-rabbitmq"
)

// Capability 能力提供者
type Capability struct {
	Category ProviderCategory
	Name     ProviderName
	Module   fx.Option
}

// Provider 底层库能力提供者接口
type Provider interface {
	Init(args ...any) error // 初始化
}
