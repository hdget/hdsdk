package mq

type PublisherOption struct {
	PublishDelayMessage bool
}

type SubscriberOption struct {
	SubscribeDelayMessage bool
}

var (
	DefaultPublisherOption = &PublisherOption{
		PublishDelayMessage: false,
	}

	DefaultSubscriberOption = &SubscriberOption{
		SubscribeDelayMessage: false,
	}
)
