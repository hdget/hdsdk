package config

import (
	"go.uber.org/fx"
)

var FxModule = fx.Module(
	"config",
	fx.Provide(NewConfigLoader),
	fx.Provide(NewSdkConfiger),
)
