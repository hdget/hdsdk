// Package mysql
// @Title  log capability of zerolog
// @Description  zerolog implementation of log capability
// @Author  Ryan Fan 2021-06-09
// @Update  Ryan Fan 2021-06-09
package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"hdsdk/provider/db"
	"hdsdk/types"
	"time"
)

type MysqlProvider struct {
	db.BaseDbProvider
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
func (mp *MysqlProvider) Init(rootConfiger types.Configer, logger types.LogProvider, args ...interface{}) error {
	// 获取数据库配置信息
	config, err := parseConfig(rootConfiger)
	if err != nil {
		return err
	}

	// 缺省数据库必须要配置合法
	err = validateConf(types.PROVIDER_TYPE_DEFAULT, config.Default)
	if err != nil {
		logger.Fatal("validate mysql config", "type", types.PROVIDER_TYPE_DEFAULT, "err", err)
	}

	// 缺省数据库必须确保能够连接成功，否则fatal
	mp.Default, err = mp.connect(config.Default)
	if err != nil {
		logger.Fatal("connect db", "type", types.PROVIDER_TYPE_DEFAULT, "host", config.Default.Host, "dbname", config.Default.Database, "err", err)
	}
	logger.Debug("connect db", "type", types.PROVIDER_TYPE_DEFAULT, "host", config.Default.Host, "dbname", config.Default.Database)

	// 主库
	if err := validateConf(types.PROVIDER_TYPE_MASTER, config.Master); err == nil {
		mp.Main, err = mp.connect(config.Master)
		logger.Debug("connect db", "type", types.PROVIDER_TYPE_MASTER, "host", config.Master.Host, "dbname", config.Master.Database, "err", err)
	}

	// 从库
	mp.Slaves = make([]*sqlx.DB, 0)
	for i, slaveConf := range config.Slaves {
		if err := validateConf(types.PROVIDER_TYPE_SLAVE, slaveConf); err == nil {
			instance, err := mp.connect(slaveConf)
			if instance != nil {
				mp.Slaves = append(mp.Slaves, instance)
			}
			logger.Debug("connect db", "type", fmt.Sprintf("slave_%d", i), "host", slaveConf.Host, "dbname", slaveConf.Database, "err", err)
		}
	}

	// 外部库
	mp.Items = make(map[string]*sqlx.DB)
	for _, otherConf := range config.Items {
		if err := validateConf(types.PROVIDER_TYPE_OTHER, otherConf); err == nil {
			instance, err := mp.connect(otherConf)
			if instance != nil {
				mp.Items[otherConf.Name] = instance
			}
			logger.Debug("connect db", "type", otherConf.Name, "host", otherConf.Host, "dbname", otherConf.Database, "err", err)
		}
	}

	return nil
}

func (mp *MysqlProvider) connect(conf *MySqlConf) (*sqlx.DB, error) {
	// DSN (Data Type NickName): username:password@protocol(address)/dbname?param=value
	t := "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4,utf8"

	// 构造连接参数
	//params := make([]string, 0)
	//timeout := fmt.Sprintf("timeout=%ds", dbConf.Timeout)
	//params = append(params, timeout)
	//strParams := strings.Join(params, "&")
	connStr := fmt.Sprintf(t, conf.User, conf.Password, conf.Host, conf.Port, conf.Database)
	instance, err := sqlx.Connect("mysql", connStr)

	// https://www.alexedwards.net/blog/configuring-sqldb
	// https://making.pusher.com/production-ready-connection-pooling-in-go
	// Avoid issue:
	// packets.go:123: closing bad idle connection: EOF
	// connection.go:173: driver: bad connection
	instance.SetConnMaxLifetime(3 * time.Minute)
	if err != nil {
		return nil, err
	}

	return instance, nil
}
