package neo4j

import (
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/hdget/hdsdk/v2/types"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type neo4jResult struct {
	neo4j.Result
}

func newResult(result neo4j.Result) intf.GraphResult {
	return &neo4jResult{
		Result: result,
	}
}

func (n neo4jResult) Record() *types.GraphRecord {
	r := n.Result.Record()
	return &types.GraphRecord{
		Values: r.Values,
		Keys:   r.Keys,
	}
}
