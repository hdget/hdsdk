package sqlboiler_sqlite3

import (
	"github.com/hdget/hdsdk/v2/intf"
	"go.uber.org/fx"
)

var Capability = &intf.Capability{
	Category: intf.ProviderCategoryDb,
	Module: fx.Module(
		string(intf.ProviderNameDbSqlBoilerSqlite),
		fx.Provide(New),
	),
}
