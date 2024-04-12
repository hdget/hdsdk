package rabbitmq

//
//type RabbitmqProvider struct {
//	connection        rabbitMQConnectionBroker
//	channel           rabbitMQChannelBroker
//	channelMutex      sync.RWMutex
//	connectionCount   int
//	declaredExchanges map[string]bool
//
//	connectionDial func(protocol, uri, clientName string, heartBeat time.Duration, tlsCfg *tls.Config, externalSasl bool) (rabbitMQConnectionBroker, rabbitMQChannelBroker, error)
//	closeCh        chan struct{}
//	closed         atomic.Bool
//	wg             sync.WaitGroup
//}
//
//// Init	implements intf.Provider interface, used to initialize the capability
//// @author	Ryan Fan	(2021-06-09)
//// @param	baseconf.Configer	root configer interface to extract configer info
//// @return	error
//func (rp *RabbitmqProvider) Init(rootConfiger intf.Configer, logger logger.LogProvider, args ...interface{}) error {
//	// 获取日志配置信息
//	configloader, err := parseConfig(rootConfiger)
//	if err != nil {
//		return err
//	}
//
//	rp.Default, err = NewMq(intf.ProviderTypeDefault, configloader.Default, logger)
//	if err != nil {
//		logger.Error("initialize mq", "type", intf.ProviderTypeDefault, "host", configloader.Default.Host, "err", err)
//	} else {
//		logger.Debug("initialize mq", "type", intf.ProviderTypeDefault, "host", configloader.Default.Host, "err", err)
//	}
//
//	// 额外的mq
//	rp.Items = make(map[string]intf.Mq)
//	for _, otherConf := range configloader.Items {
//		instance, err := NewMq(intf.ProviderTypeOther, otherConf, logger)
//		if err != nil {
//			logger.Error("initialize mq", "type", otherConf.Name, "host", otherConf.Host, "err", err)
//			continue
//		}
//
//		rp.Items[otherConf.Name] = instance
//		logger.Debug("initialize mq", "type", otherConf.Name, "host", otherConf.Host, "err", err)
//	}
//
//	return nil
//}
