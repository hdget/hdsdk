package intf

import "go.uber.org/fx"

// Provider 底层库能力提供者接口
type Provider interface {
	Init(args ...any) error // 初始化
}

type ProviderCategory int

const (
	ProviderCategoryUnknown ProviderCategory = iota
	ProviderCategoryLogger
	ProviderCategoryDb
	ProviderCategoryCache
	ProviderCategoryNeo4j
)

type ProviderName string

const (
	ProviderNameZerolog ProviderName = "zerolog"
	ProviderNameRedigo  ProviderName = "redigo"
	ProviderNameMysql   ProviderName = "mysql"
	ProviderNameNeo4j   ProviderName = "neo4j"
)

// Capability 能力提供者
type Capability struct {
	Category ProviderCategory
	Name     ProviderName
	Module   fx.Option
}
