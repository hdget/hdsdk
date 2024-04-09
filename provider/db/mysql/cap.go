package mysql

import (
	"github.com/hdget/hdsdk/v1/intf"
	"go.uber.org/fx"
)

var Capability = &intf.Capability{
	Category: intf.ProviderCategoryDb,
	Module: fx.Module(
		string(intf.ProviderNameMysql),
		fx.Provide(NewConfig),
		fx.Provide(New),
	),
}
