// Package mysql
// @Title  logger capability of zerolog
// @Description  zerolog implementation of logger capability
// @Author  Ryan Fan 2021-06-09
// @Update  Ryan Fan 2021-06-09
package mysql

import (
	_ "github.com/go-sql-driver/mysql"
)

//
//type MysqlProvider struct {
//	db.BaseDbProvider
//	Log logger.LogProvider
//}
//
//var (
//	_ intf.Provider   = (*MysqlProvider)(nil)
//	_ intf.DbProvider = (*MysqlProvider)(nil)
//)
//
//// Init	implements intf.Provider interface, used to initialize the capability
//// @author	Ryan Fan	(2021-06-09)
//// @param	baseconf.Configer	root configer interface to extract configer info
//// @return	error
//func (p *MysqlProvider) Init(rootConfiger intf.Configer, logger logger.LogProvider, _ ...interface{}) error {
//	// 获取数据库配置信息
//	configloader, err := parseConfig(rootConfiger)
//	if err != nil {
//		return err
//	}
//
//	// 缺省数据库必须要配置合法
//	err = validateConf(intf.ProviderTypeDefault, configloader.Default)
//	if err != nil {
//		logger.Fatal("validate mysql configer", "type", intf.ProviderTypeDefault, "err", err)
//	}
//
//	// 缺省数据库必须确保能够连接成功，否则fatal
//	p.Default, err = p.connect(configloader.Default)
//	if err != nil {
//		logger.Fatal("connect db", "type", intf.ProviderTypeDefault, "host", configloader.Default.Host, "dbname", configloader.Default.Database, "err", err)
//	}
//	// 为了使用sqlboiler, 这里缺省数据库必须设置进去
//	boil.SetDB(p.Default.Db())
//	logger.Debug("connect db", "type", intf.ProviderTypeDefault, "host", configloader.Default.Host, "dbname", configloader.Default.Database)
//
//	// 主库
//	if err := validateConf(intf.ProviderTypeMaster, configloader.Master); err == nil {
//		p.Main, err = p.connect(configloader.Master)
//		logger.Debug("connect db", "type", intf.ProviderTypeMaster, "host", configloader.Master.Host, "dbname", configloader.Master.Database, "err", err)
//	}
//
//	// 从库
//	p.Slaves = make([]intf.DbClient, 0)
//	for i, slaveConf := range configloader.Slaves {
//		if err := validateConf(intf.ProviderTypeSlave, slaveConf); err == nil {
//			instance, err := p.connect(slaveConf)
//			if instance != nil {
//				p.Slaves = append(p.Slaves, instance)
//			}
//			logger.Debug("connect db", "type", fmt.Sprintf("slave_%d", i), "host", slaveConf.Host, "dbname", slaveConf.Database, "err", err)
//		}
//	}
//
//	// 外部库
//	p.Items = make(map[string]intf.DbClient)
//	for _, otherConf := range configloader.Items {
//		if err := validateConf(intf.ProviderTypeOther, otherConf); err == nil {
//			instance, err := p.connect(otherConf)
//			if instance != nil {
//				p.Items[otherConf.Name] = instance
//			}
//			logger.Debug("connect db", "type", otherConf.Name, "host", otherConf.Host, "dbname", otherConf.Database, "err", err)
//		}
//	}
//
//	return nil
//}
//
//func (p *MysqlProvider) connect(conf *MySqlConf) (*db.BaseDbClient, error) {
//	// 这里设置解析时间类型https://github.com/go-sql-driver/mysql#timetime-support
//	// DSN (Data Type NickName): username:password@protocol(address)/dbname?param=value
//	t := "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local"
//	// 构造连接参数
//	connStr := fmt.Sprintf(t, conf.User, conf.Password, conf.Host, conf.Port, conf.Database)
//	instance, err := sqlx.Connect("mysql", connStr)
//	if err != nil {
//		return nil, err
//	}
//
//	// https://www.alexedwards.net/blog/configuring-sqldb
//	// https://making.pusher.com/production-ready-connection-pooling-in-go
//	// Avoid issue:
//	// packets.go:123: closing bad idle connection: EOF
//	// connection.go:173: driver: bad connection
//	instance.SetConnMaxLifetime(3 * time.Minute)
//
//	return &db.BaseDbClient{DB: instance}, nil
//}
