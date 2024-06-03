package sqlx_mysql

import (
	"github.com/hdget/hdsdk/v2/intf"
	"go.uber.org/fx"
)

var Capability = &intf.Capability{
	Category: intf.ProviderCategoryDbBuilder,
	Module: fx.Module(
		string(intf.ProviderNameSquirrelMysql),
		fx.Provide(New),
	),
}
