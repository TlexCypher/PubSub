package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub"
)

type (
	Topic        string
	Subscription string
)

type Publisher interface {
	Publish(context.Context, Topic, []byte) (string, error)
}

type Subscriber interface {
	Subscribe(ctx context.Context, s Subscription, h SubscriptionHandler) error
	GetSubscription() Subscription
	GetSubscriptionHandler() SubscriptionHandler
}

type SubscriptionHandler interface {
	Handle(ctx context.Context, msg *pubsub.Message)
}

type PubSubClient interface {
	Close() error
}
