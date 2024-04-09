package zerolog

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"logger_zerolog",
	fx.Provide(NewConfig),
	fx.Provide(New),
)
