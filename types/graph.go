package types

import "github.com/neo4j/neo4j-go-driver/v4/neo4j"

// GraphProvider 图数据库能力提供
type GraphProvider interface {
	Provider
	Get(cypher string, args ...interface{}) (interface{}, error)
	Select(cypher string, args ...interface{}) ([]interface{}, error)
	Exec(workFuncs []neo4j.TransactionWork, bookmarks ...string) (string, error)
	Reader(bookmarks ...string) neo4j.Session
	Writer(bookmarks ...string) neo4j.Session
}

// database capability
const (
	_ = SdkCategoryGraph + iota
	SdkTypeGraphNeo4j
)
