package graph

import (
	"github.com/hdget/hdsdk/v1/provider/graph/neo4j"
	"go.uber.org/fx"
)

var FxModule = fx.Module(
	"graph",
	fx.Provide(neo4j.NewConfig),
	fx.Provide(neo4j.New),
)
