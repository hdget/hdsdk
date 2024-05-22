package dapr

type EventModuleOption func(*eventModuleImpl)

//func WithConsumerTimeout(duration time.Duration) EventModuleOption {
//	return func(m *eventModuleImpl) {
//		m.consumerTimeout = duration
//	}
//}
