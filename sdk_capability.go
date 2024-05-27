package hdsdk

import "github.com/hdget/hdsdk/v2/intf"

var (
	_configProvider intf.ConfigProvider
	_logger         intf.LoggerProvider
	_db             intf.DbProvider
	_graph          intf.GraphProvider
	_redis          intf.RedisProvider
	_mq             intf.MqProvider
)

func Logger() intf.LoggerProvider {
	return _logger
}

func Db() intf.DbProvider {
	return _db
}

func Redis() intf.RedisProvider {
	return _redis
}

func Graph() intf.GraphProvider {
	return _graph
}

func Mq() intf.MqProvider {
	return _mq
}
