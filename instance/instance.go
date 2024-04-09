package instance

import "github.com/hdget/hdsdk/v1/intf"

var (
	Logger intf.LoggerProvider
	Mysql  intf.DbProvider
	Redis  intf.RedisProvider
	Neo4j  intf.GraphProvider
)

func Register(p intf.Provider) {
	switch v := p.(type) {
	case intf.LoggerProvider:
		Logger = v
	case intf.DbProvider:
		Mysql = v
	case intf.RedisProvider:
		Redis = v
	case intf.GraphProvider:
		Neo4j = v
	}
}
