package intf

import "github.com/neo4j/neo4j-go-driver/v4/neo4j"

// GraphProvider 图数据库能力提供
type GraphProvider interface {
	Provider
	Get(cypher string, args ...any) (any, error)
	Select(cypher string, args ...any) ([]any, error)
	Exec(workFuncs []neo4j.TransactionWork, bookmarks ...string) (string, error)
	Reader(bookmarks ...string) neo4j.Session
	Writer(bookmarks ...string) neo4j.Session
}
