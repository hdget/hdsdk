package logger

import (
	"github.com/hdget/hdsdk/v1/core/logger/zerolog"
	"go.uber.org/fx"
)

var FxModule = fx.Module(
	"logger",
	fx.Provide(zerolog.NewConfig),
	fx.Provide(zerolog.New),
)
