package zerolog

import (
	"github.com/hdget/hdsdk/v2/types"
	"go.uber.org/fx"
)

var Capability = &types.Capability{
	Category: types.ProviderCategoryLogger,
	Name:     types.ProviderNameZerolog,
	Module: fx.Module(
		string(types.ProviderNameZerolog),
		fx.Provide(New),
	),
}
