// Package neo4j
// @Title  logger capability of zerolog
// @Description  zerolog implementation of logger capability
// @Author  Ryan Fan 2021-06-09
// @Update  Ryan Fan 2021-06-09
package neo4j

//
//type Neo4jProvider struct {
//	BaseGraphProvider
//	Log    logger.LogProvider
//	driver neo4j.Driver
//}
//
//var (
//	_ intf.Provider      = (*Neo4jProvider)(nil)
//	_ intf.GraphProvider = (*Neo4jProvider)(nil)
//)
//
//const (
//	defaultMaxPoolSize = 500
//)
//
//// Init	implements intf.Provider interface, used to initialize the capability
//// @author	Ryan Fan	(2021-06-09)
//// @param	baseconf.Configer	root configer interface to extract configer info
//// @return	error
//func (np *Neo4jProvider) Init(rootConfiger intf.Configer, logger logger.LogProvider, args ...interface{}) error {
//	// 获取数据库配置信息
//	configloader, err := parseConfig(rootConfiger)
//	if err != nil {
//		return err
//	}
//
//	// 检查配置是否合法
//	err = validateConf(intf.ProviderTypeDefault, configloader)
//	if err != nil {
//		logger.Fatal("validate neo4j configer", "err", err)
//	}
//
//	// 看是否配置了多个server address
//	serverAddresses := make([]neo4j.ServerAddress, 0)
//	for _, server := range configloader.Servers {
//		serverAddresses = append(serverAddresses, neo4j.NewServerAddress(server.Host, cast.ToString(server.Port)))
//	}
//
//	maxPoolSize := defaultMaxPoolSize
//	if configloader.MaxPoolSize > 0 {
//		maxPoolSize = configloader.MaxPoolSize
//	}
//	np.driver, err = createDriverWithAddressResolver(configloader.VirtualUri, configloader.Username, configloader.Password, maxPoolSize, serverAddresses...)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func createDriverWithAddressResolver(virtualURI, username, password string, maxPoolSize int, addresses ...neo4j.ServerAddress) (neo4j.Driver, error) {
//	// Address resolver is only valid for neo4j uri
//	return neo4j.NewDriver(virtualURI, neo4j.BasicAuth(username, password, ""), func(configloader *neo4j.Config) {
//		configloader.AddressResolver = func(address neo4j.ServerAddress) []neo4j.ServerAddress {
//			return addresses
//		}
//		configloader.MaxConnectionPoolSize = maxPoolSize
//	})
//}
//
//func (np *Neo4jProvider) Get(cypher string, args ...interface{}) (interface{}, error) {
//	session := np.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
//	defer func(session neo4j.Session) {
//		_ = session.Close()
//	}(session)
//
//	var param map[string]interface{}
//	if len(args) > 0 {
//		param = structs.Map(args[0])
//	}
//
//	result, err := session.Run(cypher, param)
//	if err != nil {
//		return nil, err
//	}
//
//	var ret interface{}
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
//func (np *Neo4jProvider) Select(cypher string, args ...interface{}) ([]interface{}, error) {
//	session := np.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
//	defer func(session neo4j.Session) {
//		_ = session.Close()
//	}(session)
//
//	var param map[string]interface{}
//	if len(args) > 0 {
//		param = structs.Map(args[0])
//	}
//
//	result, err := session.Run(cypher, param)
//	if err != nil {
//		return nil, err
//	}
//
//	rets := make([]interface{}, 0)
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
//func (np *Neo4jProvider) Exec(workFuncs []neo4j.TransactionWork, bookmarks ...string) (string, error) {
//	session := np.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, Bookmarks: bookmarks})
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
//func (np *Neo4jProvider) Reader(bookmarks ...string) neo4j.Session {
//	return np.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead, Bookmarks: bookmarks})
//}
//
//func (np *Neo4jProvider) Writer(bookmarks ...string) neo4j.Session {
//	return np.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, Bookmarks: bookmarks})
//}
//
//func (np *Neo4jProvider) Read(cypher string, params map[string]interface{}, bookmarks ...string) neo4j.Session {
//	return np.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead, Bookmarks: bookmarks})
//}
