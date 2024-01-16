package logger

import (
	"github.com/hdget/hdsdk/core/logger/zerolog"
	"go.uber.org/fx"
)

var FxModule = fx.Options(
	fx.Provide(zerolog.NewConfig),
	fx.Provide(zerolog.New),
)
