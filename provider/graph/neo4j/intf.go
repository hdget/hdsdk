package neo4j

import "github.com/neo4j/neo4j-go-driver/v4/neo4j"

type Provider interface {
	Get(cypher string, args ...interface{}) (interface{}, error)
	Select(cypher string, args ...interface{}) ([]interface{}, error)
	Exec(workFuncs []neo4j.TransactionWork, bookmarks ...string) (string, error)
	Reader(bookmarks ...string) neo4j.Session
	Writer(bookmarks ...string) neo4j.Session
}
