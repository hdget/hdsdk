package neo4j

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"graph_neo4j",
	fx.Provide(NewConfig),
	fx.Provide(New),
)
