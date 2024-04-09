package neo4j

import (
	"github.com/hdget/hdsdk/v1/intf"
	"go.uber.org/fx"
)

var Capability = &intf.Capability{
	Category: intf.ProviderCategoryNeo4j,
	Name:     intf.ProviderNameNeo4j,
	Module: fx.Module(
		string(intf.ProviderNameNeo4j),
		fx.Provide(NewConfig),
		fx.Provide(New),
	),
}
