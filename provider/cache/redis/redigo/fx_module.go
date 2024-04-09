package redigo

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"cache_redigo",
	fx.Provide(NewConfig),
	fx.Provide(New),
)
