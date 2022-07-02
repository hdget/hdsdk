package graph

import "github.com/neo4j/neo4j-go-driver/v4/neo4j"

type BaseGraphProvider struct {
}

func (c *BaseGraphProvider) Get(cypher string, args ...interface{}) (interface{}, error) {
	return nil, nil
}

func (c *BaseGraphProvider) Select(cypher string, args ...interface{}) ([]interface{}, error) {
	return nil, nil
}

func (c *BaseGraphProvider) Exec(workFuncs []neo4j.TransactionWork, bookmarks ...string) (string, error) {
	return "", nil
}

func (c *BaseGraphProvider) Reader(bookmarks ...string) neo4j.Session {
	return nil
}

func (c *BaseGraphProvider) Writer(bookmarks ...string) neo4j.Session {
	return nil
}
