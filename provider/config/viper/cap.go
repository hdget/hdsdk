package viper

import (
	"github.com/hdget/hdsdk/v2/types"
	"go.uber.org/fx"
)

var Capability = &types.Capability{
	Category: types.ProviderCategoryConfig,
	Module: fx.Module(
		string(types.ProviderNameViper),
		fx.Provide(New),
	),
}
