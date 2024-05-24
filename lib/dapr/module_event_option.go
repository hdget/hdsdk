package dapr

import "time"

type EventModuleOption func(*eventModuleImpl)

func WithConsumerTimeout(duration time.Duration) EventModuleOption {
	return func(m *eventModuleImpl) {
		m.ackTimeout = duration
	}
}
