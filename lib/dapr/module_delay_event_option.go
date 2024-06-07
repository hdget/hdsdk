package dapr

import (
	"github.com/cenkalti/backoff/v4"
)

type DelayEventModuleOption func(*delayEventModuleImpl)

func WithBackOff(backoff *backoff.ExponentialBackOff) DelayEventModuleOption {
	return func(m *delayEventModuleImpl) {
		m.backoffPolicy = backoff
	}
}
