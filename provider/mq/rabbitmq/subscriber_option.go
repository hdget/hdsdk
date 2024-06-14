package rabbitmq

type subscriberOption func(impl *rmpSubscriberImpl)

func withSubscriberDelayTopology() subscriberOption {
	return func(impl *rmpSubscriberImpl) {
		impl.useDelayTopology = true
	}
}
