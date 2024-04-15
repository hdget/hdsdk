package hdsdk

import "github.com/hdget/hdsdk/v2/intf"

func Logger() intf.LoggerProvider {
	return _instance.logger
}

func Db() intf.DbProvider {
	return _instance.db
}

func Redis() intf.RedisProvider {
	return _instance.redis
}

//func Graph() intf.GraphProvider {
//	return _instance.graph
//}

func GetName() string {
	return ""
}
