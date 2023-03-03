// Package mysql
// @Title  log capability of zerolog
// @Description  zerolog implementation of log capability
// @Author  Ryan Fan 2021-06-09
// @Update  Ryan Fan 2021-06-09
package neo4j

import (
	"github.com/fatih/structs"
	"github.com/hdget/hdsdk/provider/graph"
	"github.com/hdget/hdsdk/types"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/spf13/cast"
)

type Neo4jProvider struct {
	graph.BaseGraphProvider
	Log    types.LogProvider
	driver neo4j.Driver
}

var (
	_ types.Provider      = (*Neo4jProvider)(nil)
	_ types.GraphProvider = (*Neo4jProvider)(nil)
)

const (
	defaultMaxPoolSize = 500
)

// Init	implements types.Provider interface, used to initialize the capability
// @author	Ryan Fan	(2021-06-09)
// @param	baseconf.Configer	root config interface to extract config info
// @return	error
func (np *Neo4jProvider) Init(rootConfiger types.Configer, logger types.LogProvider, args ...interface{}) error {
	// 获取数据库配置信息
	config, err := parseConfig(rootConfiger)
	if err != nil {
		return err
	}

	// 检查配置是否合法
	err = validateConf(types.PROVIDER_TYPE_DEFAULT, config)
	if err != nil {
		logger.Fatal("validate neo4j config", "err", err)
	}

	// 看是否配置了多个server address
	serverAddresses := make([]neo4j.ServerAddress, 0)
	for _, server := range config.Servers {
		serverAddresses = append(serverAddresses, neo4j.NewServerAddress(server.Host, cast.ToString(server.Port)))
	}

	maxPoolSize := defaultMaxPoolSize
	if config.MaxPoolSize > 0 {
		maxPoolSize = config.MaxPoolSize
	}
	np.driver, err = createDriverWithAddressResolver(config.VirtualUri, config.Username, config.Password, maxPoolSize, serverAddresses...)
	if err != nil {
		return err
	}

	return nil
}

func createDriverWithAddressResolver(virtualURI, username, password string, maxPoolSize int, addresses ...neo4j.ServerAddress) (neo4j.Driver, error) {
	// Address resolver is only valid for neo4j uri
	return neo4j.NewDriver(virtualURI, neo4j.BasicAuth(username, password, ""), func(config *neo4j.Config) {
		config.AddressResolver = func(address neo4j.ServerAddress) []neo4j.ServerAddress {
			return addresses
		}
		config.MaxConnectionPoolSize = maxPoolSize
	})
}

func (np *Neo4jProvider) Get(cypher string, args ...interface{}) (interface{}, error) {
	session := np.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	var param map[string]interface{}
	if len(args) > 0 {
		param = structs.Map(args[0])
	}

	result, err := session.Run(cypher, param)
	if err != nil {
		return nil, err
	}

	var ret interface{}
	if result.Next() {
		ret = result.Record().Values[0]
	}

	if err = result.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}

func (np *Neo4jProvider) Select(cypher string, args ...interface{}) ([]interface{}, error) {
	session := np.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	var param map[string]interface{}
	if len(args) > 0 {
		param = structs.Map(args[0])
	}

	result, err := session.Run(cypher, param)
	if err != nil {
		return nil, err
	}

	rets := make([]interface{}, 0)
	for result.Next() {
		rets = append(rets, result.Record().Values[0])
	}

	if err = result.Err(); err != nil {
		return nil, err
	}

	return rets, nil
}

func (np *Neo4jProvider) Exec(workFuncs []neo4j.TransactionWork, bookmarks ...string) (string, error) {
	session := np.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, Bookmarks: bookmarks})
	defer session.Close()

	for _, fnWork := range workFuncs {
		if _, err := session.WriteTransaction(fnWork); err != nil {
			return "", err
		}
	}

	return session.LastBookmark(), nil
}

func (np *Neo4jProvider) Reader(bookmarks ...string) neo4j.Session {
	return np.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead, Bookmarks: bookmarks})
}

func (np *Neo4jProvider) Writer(bookmarks ...string) neo4j.Session {
	return np.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, Bookmarks: bookmarks})
}

func (np *Neo4jProvider) Read(cypher string, params map[string]interface{}, bookmarks ...string) neo4j.Session {
	return np.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead, Bookmarks: bookmarks})
}
