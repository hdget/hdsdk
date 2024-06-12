package redigo

import (
	"github.com/hdget/hdsdk/v2/intf"
	"go.uber.org/fx"
)

var Capability = &intf.Capability{
	Category: intf.ProviderCategoryRedis,
	Name:     intf.ProviderNameRedisRedigo,
	Module: fx.Module(
		string(intf.ProviderNameRedisRedigo),
		fx.Provide(New),
	),
}
