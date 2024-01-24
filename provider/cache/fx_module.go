package cache

import (
	redigo2 "github.com/hdget/hdsdk/v1/provider/cache/redis/redigo"
	"go.uber.org/fx"
)

var FxModule = fx.Module(
	"cache",
	fx.Provide(redigo2.NewConfig),
	fx.Provide(redigo2.New),
)
