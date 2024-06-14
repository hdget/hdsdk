package rabbitmq

type publisherOption func(impl *rmqPublisherImpl)

func withPublisherDelayTopology() publisherOption {
	return func(impl *rmqPublisherImpl) {
		impl.useDelayTopology = true
	}
}
