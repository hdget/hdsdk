package neo4j

import (
	"github.com/hdget/hdsdk/v2/intf"
	"go.uber.org/fx"
)

var Capability = &intf.Capability{
	Category: intf.ProviderCategoryGraph,
	Name:     intf.ProviderNameNeo4j,
	Module: fx.Module(
		string(intf.ProviderNameNeo4j),
		fx.Provide(New),
	),
}
