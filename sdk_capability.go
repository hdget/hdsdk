package hdsdk

import (
	"github.com/hdget/hdsdk/v1/instance"
	"github.com/hdget/hdsdk/v1/intf"
)

func Logger() intf.LoggerProvider {
	return instance.Logger
}

func Mysql() intf.DbProvider {
	return instance.Mysql
}

func Redis() intf.RedisProvider {
	return instance.Redis
}

func Neo4j() intf.GraphProvider {
	return instance.Neo4j
}
