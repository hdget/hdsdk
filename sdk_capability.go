package hdsdk

import "github.com/hdget/hdsdk/v2/intf"

func Logger() intf.LoggerProvider {
	return _instance.logger
}

func Db() intf.DbProvider {
	return _instance.db
}

func DbBuilder(sqlizer intf.Sqlizer) intf.DbBuilderProvider {
	_instance.dbBuilder.Set(sqlizer)
	return _instance.dbBuilder
}

func Sqlx() intf.SqlxDbProvider {
	return _instance.sqlxDb
}

func Redis() intf.RedisProvider {
	return _instance.redis
}

func Mq() intf.MqProvider {
	return _instance.mq
}

func Config() intf.ConfigProvider {
	return _instance.configProvider
}
