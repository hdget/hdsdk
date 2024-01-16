package config

import (
	"go.uber.org/fx"
)

var FxModule = fx.Options(
	fx.Provide(NewConfigLoader),
	fx.Provide(NewConfiger),
)
