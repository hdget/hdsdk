package sqlx

import (
	"github.com/hdget/hdsdk/v2/intf"
	"go.uber.org/fx"
)

var Capability = &intf.Capability{
	Category: intf.ProviderCategoryDb,
	Module: fx.Module(
		string(intf.ProviderNameSqlx),
		fx.Provide(New),
	),
}
