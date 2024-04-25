package intf

import (
	"github.com/hdget/hdsdk/v2/provider/graph"
)

// GraphProvider 图数据库能力提供
type GraphProvider interface {
	Provider
	Exec(transFunctions []any) (string, error)
	Reader() GraphSession
	Writer() GraphSession
}

type GraphSession interface {
	// Run executes an auto-commit statement and returns a result
	Run(cypher string, params map[string]interface{}) (GraphResult, error)
	// Close closes any open resources and marks this session as unusable
	Close() error
}

type GraphResult interface {
	// Next returns true only if there is a record to be processed.
	Next() bool
	// Record returns the current record.
	Record() *graph.Record
}
