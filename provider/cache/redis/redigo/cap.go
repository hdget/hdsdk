package redigo

import (
	"github.com/hdget/hdsdk/v1/intf"
	"go.uber.org/fx"
)

var Capability = &intf.Capability{
	Category: intf.ProviderCategoryCache,
	Name:     intf.ProviderNameRedigo,
	Module: fx.Module(
		string(intf.ProviderNameRedigo),
		fx.Provide(NewConfig),
		fx.Provide(New),
	),
}
