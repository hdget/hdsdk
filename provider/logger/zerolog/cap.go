package zerolog

import (
	"github.com/hdget/hdsdk/v2/intf"
	"go.uber.org/fx"
)

var Capability = &intf.Capability{
	Category: intf.ProviderCategoryLogger,
	Name:     intf.ProviderNameLoggerZerolog,
	Module: fx.Module(
		string(intf.ProviderNameLoggerZerolog),
		fx.Provide(New),
	),
}
