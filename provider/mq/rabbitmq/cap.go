package rabbitmq

import (
	"github.com/hdget/hdsdk/v2/intf"
	"go.uber.org/fx"
)

var Capability = &intf.Capability{
	Category: intf.ProviderCategoryMq,
	Name:     intf.ProviderNameRabbitMq,
	Module: fx.Module(
		string(intf.ProviderNameRabbitMq),
		fx.Provide(New),
	),
}
