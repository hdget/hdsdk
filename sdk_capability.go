package hdsdk

import "github.com/hdget/hdsdk/v1/intf"

func Logger() intf.LoggerProvider {
	return _instance.logger
}

func Db() intf.DbProvider {
	return _instance.db
}

func Cache() intf.CacheProvider {
	return _instance.cache
}

func Graph() intf.GraphProvider {
	return _instance.graph
}
