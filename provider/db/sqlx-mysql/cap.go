package sqlx_mysql

import (
	"github.com/hdget/hdsdk/v2/intf"
	"go.uber.org/fx"
)

var Capability = &intf.Capability{
	Category: intf.ProviderCategoryDbSqlx,
	Module: fx.Module(
		string(intf.ProviderNameDbSqlxMysql),
		fx.Provide(New),
	),
}
