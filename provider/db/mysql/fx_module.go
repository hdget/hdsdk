package mysql

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"db_mysql",
	fx.Provide(NewConfig),
	fx.Provide(New),
)
