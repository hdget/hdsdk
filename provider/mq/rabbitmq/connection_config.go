package rabbitmq

//
//type connectionConfig struct {
//	host      string // RabbitMQ host
//	port      int    // RabbitMQ port
//	username  string // RabbitMQ username
//	password  string // RabbitMQ password
//	vhost     string
//	reconnect *reconnectConfig // reconnect config
//}
//
//// increases the back off period for each retry attempt using a randomization function that grows exponentially
//type reconnectConfig struct {
//	backoffInitialInterval     time.Duration
//	backoffRandomizationFactor float64
//	backoffMultiplier          float64
//	backoffMaxInterval         time.Duration
//}
//
//var (
//	defaultReconnectConfig = &reconnectConfig{
//		backoffInitialInterval:     500 * time.Millisecond,
//		backoffRandomizationFactor: 0.5,
//		backoffMultiplier:          1.5,
//		backoffMaxInterval:         60 * time.Second,
//	}
//)
//
//func (c *connectionConfig) validate() error {
//	if c.host == "" || c.username == "" || c.password == "" {
//		return errdef.ErrInvalidConfig
//	}
//
//	if c.reconnect == nil {
//		c.reconnect = defaultReconnectConfig
//	}
//
//	if c.port == 0 {
//		c.port = 5672
//	}
//
//	return nil
//}
//
//func (c *connectionConfig) getExponentialBackOff() *backoff.ExponentialBackOff {
//	return &backoff.ExponentialBackOff{
//		InitialInterval:     c.reconnect.backoffInitialInterval,
//		RandomizationFactor: c.reconnect.backoffRandomizationFactor,
//		Multiplier:          c.reconnect.backoffMultiplier,
//		MaxInterval:         c.reconnect.backoffMaxInterval,
//		MaxElapsedTime:      0, // no support for disabling reconnect, only close of Pub/Sub can stop reconnecting
//		Clock:               backoff.SystemClock,
//	}
//}
