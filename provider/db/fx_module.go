package db

import (
	"github.com/hdget/hdsdk/provider/db/mysql"
	"go.uber.org/fx"
)

var FxModule = fx.Module(
	"db",
	fx.Provide(mysql.NewConfig),
	fx.Provide(mysql.New),
)
