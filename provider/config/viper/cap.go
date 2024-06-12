package viper

import (
	"github.com/hdget/hdsdk/v2/intf"
	"go.uber.org/fx"
)

var Capability = &intf.Capability{
	Category: intf.ProviderCategoryConfig,
	Module: fx.Module(
		string(intf.ProviderNameConfigViper),
		fx.Provide(New),
	),
}
