package types

import "go.uber.org/fx"

type ProviderCategory int

const (
	ProviderCategoryUnknown ProviderCategory = iota
	ProviderCategoryConfig
	ProviderCategoryLogger
	ProviderCategoryDb
	ProviderCategoryRedis
	ProviderCategoryGraph
)

type ProviderName string

const (
	ProviderNameViper     ProviderName = "viper"
	ProviderNameZerolog   ProviderName = "zerolog"
	ProviderNameRedigo    ProviderName = "redigo"
	ProviderNameSqlBoiler ProviderName = "sqlboiler"
	ProviderNameSqlx      ProviderName = "sqlx"
	ProviderNameNeo4j     ProviderName = "neo4j"
)

// Capability 能力提供者
type Capability struct {
	Category ProviderCategory
	Name     ProviderName
	Module   fx.Option
}
