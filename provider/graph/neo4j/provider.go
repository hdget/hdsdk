// Package neo4j
// @Title  logger capability of zerolog
// @Description  zerolog implementation of logger capability
// @Author  Ryan Fan 2021-06-09
// @Update  Ryan Fan 2021-06-09
package neo4j

import (
	"github.com/fatih/structs"
	"github.com/hdget/hdsdk/v1/intf"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

type neo4jProvider struct {
	driver neo4j.Driver
}

func New(providerConfig *neo4jProviderConfig, logger intf.LoggerProvider) (intf.GraphProvider, error) {
	provider := &neo4jProvider{}
	err := provider.Init(logger, providerConfig)
	if err != nil {
		logger.Fatal("init mysql provider", "err", err)
	}

	return provider, nil
}

// Init	implements intf.Provider interface, used to initialize the capability
// @author	Ryan Fan	(2021-06-09)
// @param	baseconf.Configer	root configer interface to extract configer info
// @return	error
func (p *neo4jProvider) Init(logger intf.LoggerProvider, args ...any) error {
	if len(args) == 0 {
		return errors.New("need neo4j provider config")
	}

	providerConfig, ok := args[0].(*neo4jProviderConfig)
	if !ok {
		return errors.New("invalid neo4j config")
	}

	var err error
	p.driver, err = newNeo4jDriver(providerConfig)
	if err != nil {
		return err
	}

	return nil
}

func newNeo4jDriver(providerConfig *neo4jProviderConfig) (neo4j.Driver, error) {
	// Address resolver is only valid for neo4j uri
	return neo4j.NewDriver(
		providerConfig.VirtualUri,
		neo4j.BasicAuth(providerConfig.Username, providerConfig.Password, ""),
		func(conf *neo4j.Config) {
			conf.AddressResolver = func(address neo4j.ServerAddress) []neo4j.ServerAddress {
				serverAddresses := make([]neo4j.ServerAddress, 0)
				for _, server := range providerConfig.Servers {
					serverAddresses = append(serverAddresses, neo4j.NewServerAddress(server.Host, cast.ToString(server.Port)))
				}
				return serverAddresses
			}
			conf.MaxConnectionPoolSize = providerConfig.MaxPoolSize
		})
}

func (p *neo4jProvider) Get(cypher string, args ...interface{}) (interface{}, error) {
	session := p.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer func(session neo4j.Session) {
		_ = session.Close()
	}(session)

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

func (p *neo4jProvider) Select(cypher string, args ...interface{}) ([]interface{}, error) {
	session := p.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer func(session neo4j.Session) {
		_ = session.Close()
	}(session)

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

func (p *neo4jProvider) Exec(workFuncs []neo4j.TransactionWork, bookmarks ...string) (string, error) {
	session := p.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, Bookmarks: bookmarks})
	defer func(session neo4j.Session) {
		_ = session.Close()
	}(session)

	for _, fnWork := range workFuncs {
		if _, err := session.WriteTransaction(fnWork); err != nil {
			return "", err
		}
	}

	return session.LastBookmark(), nil
}

func (p *neo4jProvider) Reader(bookmarks ...string) neo4j.Session {
	return p.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead, Bookmarks: bookmarks})
}

func (p *neo4jProvider) Writer(bookmarks ...string) neo4j.Session {
	return p.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, Bookmarks: bookmarks})
}

func (p *neo4jProvider) Read(cypher string, params map[string]interface{}, bookmarks ...string) neo4j.Session {
	return p.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead, Bookmarks: bookmarks})
}
