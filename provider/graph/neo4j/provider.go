// Package neo4j
// @Title  logger capability of zerolog
// @Description  zerolog implementation of logger capability
// @Author  Ryan Fan 2021-06-09
// @Update  Ryan Fan 2021-06-09
package neo4j

//
//type neo4jProvider struct {
//	logger intf.LoggerProvider
//	config *neo4jProviderConfig
//	driver neo4j.Driver
//}
//
//func New(config *neo4jProviderConfig, logger intf.LoggerProvider) (intf.GraphProvider, error) {
//	provider := &neo4jProvider{
//		logger: logger,
//		config: config,
//	}
//
//	err := provider.Init()
//	if err != nil {
//		logger.Fatal("init mysql provider", "err", err)
//	}
//
//	return provider, nil
//}
//
//// Init	initialize neo4j driver
//func (p *neo4jProvider) Init(args ...any) error {
//	var err error
//	p.driver, err = p.newNeo4jDriver()
//	if err != nil {
//		return err
//	}
//	p.logger.Debug("init neo4j", "uri", p.config.VirtualUri)
//
//	return nil
//}
//
//func (p *neo4jProvider) Get(cypher string, args ...any) (any, error) {
//	session := p.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
//	defer func(session neo4j.Session) {
//		_ = session.Close()
//	}(session)
//
//	var param map[string]any
//	if len(args) > 0 {
//		param = structs.Map(args[0])
//	}
//
//	result, err := session.Run(cypher, param)
//	if err != nil {
//		return nil, err
//	}
//
//	var ret any
//	if result.Next() {
//		ret = result.Record().Values[0]
//	}
//
//	if err = result.Err(); err != nil {
//		return nil, err
//	}
//
//	return ret, nil
//}
//
//func (p *neo4jProvider) Select(cypher string, args ...any) ([]any, error) {
//	session := p.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
//	defer func(session neo4j.Session) {
//		_ = session.Close()
//	}(session)
//
//	var param map[string]any
//	if len(args) > 0 {
//		param = structs.Map(args[0])
//	}
//
//	result, err := session.Run(cypher, param)
//	if err != nil {
//		return nil, err
//	}
//
//	rets := make([]any, 0)
//	for result.Next() {
//		rets = append(rets, result.Record().Values[0])
//	}
//
//	if err = result.Err(); err != nil {
//		return nil, err
//	}
//
//	return rets, nil
//}
//
//func (p *neo4jProvider) Exec(workFuncs []neo4j.TransactionWork, bookmarks ...string) (string, error) {
//	session := p.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, Bookmarks: bookmarks})
//	defer func(session neo4j.Session) {
//		_ = session.Close()
//	}(session)
//
//	for _, fnWork := range workFuncs {
//		if _, err := session.WriteTransaction(fnWork); err != nil {
//			return "", err
//		}
//	}
//
//	return session.LastBookmark(), nil
//}
//
//func (p *neo4jProvider) Reader(bookmarks ...string) neo4j.Session {
//	return p.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead, Bookmarks: bookmarks})
//}
//
//func (p *neo4jProvider) Writer(bookmarks ...string) neo4j.Session {
//	return p.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, Bookmarks: bookmarks})
//}
//
//func (p *neo4jProvider) Read(cypher string, params map[string]any, bookmarks ...string) neo4j.Session {
//	return p.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead, Bookmarks: bookmarks})
//}
//
//func (p *neo4jProvider) newNeo4jDriver() (neo4j.Driver, error) {
//	// Address resolver is only valid for neo4j uri
//	driver, err := neo4j.NewDriver(
//		p.config.VirtualUri,
//		neo4j.BasicAuth(p.config.Username, p.config.Password, ""),
//		func(conf *neo4j.Config) {
//			conf.AddressResolver = func(address neo4j.ServerAddress) []neo4j.ServerAddress {
//				serverAddresses := make([]neo4j.ServerAddress, 0)
//				for _, server := range p.config.Servers {
//					serverAddresses = append(serverAddresses, neo4j.NewServerAddress(server.Host, cast.ToString(server.Port)))
//				}
//				return serverAddresses
//			}
//			conf.MaxConnectionPoolSize = p.config.MaxPoolSize
//		})
//	if err != nil {
//		return nil, err
//	}
//
//	// check if neo4j can be connected or not
//	_, err = driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead}).Run(
//		`CALL dbms.components() YIELD name, versions, edition RETURN name, versions, edition`,
//		nil)
//	if err != nil {
//		return nil, err
//	}
//
//	return driver, nil
//}
