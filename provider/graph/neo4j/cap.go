package neo4j

import (
	"github.com/hdget/hdsdk/v2/types"
	"go.uber.org/fx"
)

var Capability = &types.Capability{
	Category: types.ProviderCategoryGraph,
	Name:     types.ProviderNameNeo4j,
	Module: fx.Module(
		string(types.ProviderNameNeo4j),
		fx.Provide(New),
	),
}
