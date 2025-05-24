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
	Subscribe(ctx context.Context, subscription Subscription) error
}

type SubscriptionHandler interface {
	Handle(ctx context.Context, msg *pubsub.Message) error
}

type PubSubClient interface {
	GetClient() *pubsub.Client
	Close() error
}
