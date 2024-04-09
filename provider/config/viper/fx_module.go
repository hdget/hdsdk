package viper

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"config_viper",
	fx.Provide(New),
)
