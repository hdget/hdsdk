package redigo

import (
	"github.com/hdget/hdsdk/v2/types"
	"go.uber.org/fx"
)

var Capability = &types.Capability{
	Category: types.ProviderCategoryRedis,
	Name:     types.ProviderNameRedigo,
	Module: fx.Module(
		string(types.ProviderNameRedigo),
		fx.Provide(New),
	),
}
