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
	ProviderNameViper           ProviderName = "viper"
	ProviderNameZerolog         ProviderName = "zerolog"
	ProviderNameRedigo          ProviderName = "redigo"
	ProviderNameSqlBoilerMysql  ProviderName = "sqlboiler-mysql"
	ProviderNameSqlBoilerSqlite ProviderName = "sqlboiler-sqlite3"
	ProviderNameSqlxMysql       ProviderName = "sqlx-mysql"
	ProviderNameSquirrelMysql   ProviderName = "squirrel-mysql"
	ProviderNameRabbitMq        ProviderName = "rabbitmq"
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
