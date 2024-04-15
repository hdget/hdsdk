// Package neo4j
package neo4j

import (
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

type neo4jProvider struct {
	logger intf.LoggerProvider
	config *neo4jProviderConfig
	driver neo4j.Driver
}

var (
	errInvalidTransactionFunction = errors.New("invalid neo4j.TransactionWork function")
)

func New(configProvider intf.ConfigProvider, logger intf.LoggerProvider) (intf.GraphProvider, error) {
	c, err := newConfig(configProvider)
	if err != nil {
		return nil, err
	}

	provider := &neo4jProvider{
		logger: logger,
		config: c,
	}

	err = provider.Init()
	if err != nil {
		logger.Fatal("init neo4j provider", "err", err)
	}

	return provider, nil
}

// Init	initialize neo4j driver
func (p *neo4jProvider) Init(args ...any) error {
	var err error
	p.driver, err = p.newNeo4jDriver()
	if err != nil {
		return err
	}
	p.logger.Debug("init neo4j", "uri", p.config.VirtualUri)

	return nil
}

func (p *neo4jProvider) Exec(transFunctions []any) (string, error) {
	session := p.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func(session neo4j.Session) {
		_ = session.Close()
	}(session)

	for _, fn := range transFunctions {
		workFn, ok := fn.(neo4j.TransactionWork)
		if !ok {
			return "", errInvalidTransactionFunction
		}
		if _, err := session.WriteTransaction(workFn); err != nil {
			return "", err
		}
	}

	return session.LastBookmark(), nil
}

func (p *neo4jProvider) Reader() intf.GraphSession {
	return newNeo4jSession(p.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead}))
}

func (p *neo4jProvider) Writer() intf.GraphSession {
	return newNeo4jSession(p.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite}))
}

func (p *neo4jProvider) newNeo4jDriver() (neo4j.Driver, error) {
	// Address resolver is only valid for neo4j uri
	driver, err := neo4j.NewDriver(
		p.config.VirtualUri,
		neo4j.BasicAuth(p.config.Username, p.config.Password, ""),
		func(conf *neo4j.Config) {
			conf.AddressResolver = func(address neo4j.ServerAddress) []neo4j.ServerAddress {
				serverAddresses := make([]neo4j.ServerAddress, 0)
				for _, server := range p.config.Servers {
					serverAddresses = append(serverAddresses, neo4j.NewServerAddress(server.Host, cast.ToString(server.Port)))
				}
				return serverAddresses
			}
			conf.MaxConnectionPoolSize = p.config.MaxPoolSize
		})
	if err != nil {
		return nil, err
	}

	// check if neo4j can be connected or not
	_, err = driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead}).Run(
		`CALL dbms.components() YIELD name, versions, edition RETURN name, versions, edition`,
		nil)
	if err != nil {
		return nil, err
	}

	return driver, nil
}
