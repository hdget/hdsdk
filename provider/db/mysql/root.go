// Package mysql
// @Title  log capability of zerolog
// @Description  zerolog implementation of log capability
// @Author  Ryan Fan 2021-06-09
// @Update  Ryan Fan 2021-06-09
package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hdget/hdsdk/provider/db"
	"github.com/hdget/hdsdk/types"
	"github.com/jmoiron/sqlx"
	"time"
)

type MysqlProvider struct {
	*db.BaseDbProvider
	Log types.LogProvider
}

var (
	_ types.Provider   = (*MysqlProvider)(nil)
	_ types.DbProvider = (*MysqlProvider)(nil)
)

// Init	implements types.Provider interface, used to initialize the capability
// @author	Ryan Fan	(2021-06-09)
// @param	baseconf.Configer	root config interface to extract config info
// @return	error
func (mp *MysqlProvider) Init(rootConfiger types.Configer, logger types.LogProvider, _ ...interface{}) error {
	// 获取数据库配置信息
	config, err := parseConfig(rootConfiger)
	if err != nil {
		return err
	}

	// 缺省数据库必须要配置合法
	err = validateConf(types.ProviderTypeDefault, config.Default)
	if err != nil {
		logger.Fatal("validate mysql config", "type", types.ProviderTypeDefault, "err", err)
	}

	// 缺省数据库必须确保能够连接成功，否则fatal
	mp.Default, err = mp.connect(config.Default)
	if err != nil {
		logger.Fatal("connect db", "type", types.ProviderTypeDefault, "host", config.Default.Host, "dbname", config.Default.Database, "err", err)
	}
	logger.Debug("connect db", "type", types.ProviderTypeDefault, "host", config.Default.Host, "dbname", config.Default.Database)

	// 主库
	if err := validateConf(types.ProviderTypeMaster, config.Master); err == nil {
		mp.Main, err = mp.connect(config.Master)
		logger.Debug("connect db", "type", types.ProviderTypeMaster, "host", config.Master.Host, "dbname", config.Master.Database, "err", err)
	}

	// 从库
	mp.Slaves = make([]*db.BaseDbClient, 0)
	for i, slaveConf := range config.Slaves {
		if err := validateConf(types.ProviderTypeSlave, slaveConf); err == nil {
			instance, err := mp.connect(slaveConf)
			if instance != nil {
				mp.Slaves = append(mp.Slaves, instance)
			}
			logger.Debug("connect db", "type", fmt.Sprintf("slave_%d", i), "host", slaveConf.Host, "dbname", slaveConf.Database, "err", err)
		}
	}

	// 外部库
	mp.Items = make(map[string]*db.BaseDbClient)
	for _, otherConf := range config.Items {
		if err := validateConf(types.ProviderTypeOther, otherConf); err == nil {
			instance, err := mp.connect(otherConf)
			if instance != nil {
				mp.Items[otherConf.Name] = instance
			}
			logger.Debug("connect db", "type", otherConf.Name, "host", otherConf.Host, "dbname", otherConf.Database, "err", err)
		}
	}

	return nil
}

func (mp *MysqlProvider) connect(conf *MySqlConf) (*db.BaseDbClient, error) {
	// DSN (Data Type NickName): username:password@protocol(address)/dbname?param=value
	t := "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4,utf8"

	// 构造连接参数
	connStr := fmt.Sprintf(t, conf.User, conf.Password, conf.Host, conf.Port, conf.Database)
	instance, err := sqlx.Connect("mysql", connStr)
	if err != nil {
		return nil, err
	}

	// https://www.alexedwards.net/blog/configuring-sqldb
	// https://making.pusher.com/production-ready-connection-pooling-in-go
	// Avoid issue:
	// packets.go:123: closing bad idle connection: EOF
	// connection.go:173: driver: bad connection
	instance.SetConnMaxLifetime(3 * time.Minute)

	return &db.BaseDbClient{DB: instance}, nil
}
