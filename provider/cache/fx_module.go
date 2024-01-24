package cache

import (
	"github.com/hdget/hdsdk/v1/provider/cache/redis/redigo"
	"go.uber.org/fx"
)

var FxModule = fx.Module(
	"cache",
	fx.Provide(redigo.NewConfig),
	fx.Provide(redigo.New),
)
