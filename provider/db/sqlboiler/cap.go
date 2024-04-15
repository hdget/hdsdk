package sqlboiler

import (
	"github.com/hdget/hdsdk/v2/types"
	"go.uber.org/fx"
)

var Capability = &types.Capability{
	Category: types.ProviderCategoryDb,
	Module: fx.Module(
		string(types.ProviderNameSqlBoiler),
		fx.Provide(New),
	),
}
