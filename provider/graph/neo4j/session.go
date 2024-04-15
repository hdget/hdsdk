package neo4j

import (
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type neo4jSession struct {
	session neo4j.Session
}

var (
	_ intf.GraphSession = (*neo4jSession)(nil)
)

func newNeo4jSession(session neo4j.Session) *neo4jSession {
	return &neo4jSession{
		session: session,
	}
}

func (n neo4jSession) Run(cypher string, params map[string]interface{}) (intf.GraphResult, error) {
	ret, err := n.session.Run(cypher, params)
	if err != nil {
		return nil, err
	}

	return newResult(ret), nil
}

func (n neo4jSession) Close() error {
	return n.session.Close()
}
