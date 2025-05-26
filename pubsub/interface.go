package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub"
)

type (
	Topic          string
	SubscriptionID string
)

type Publisher interface {
	Publish(context.Context, Topic, []byte) (string, error)
}

type Subscriber interface {
	Subscribe(context.Context, Subscription) error
}

type Subscription interface {
	Receive(context.Context, func(context.Context, *pubsub.Message)) error
}

type PubSubClient interface {
	GetClient() *pubsub.Client
	Subscription(SubscriptionID) Subscription
	Close() error
}
