package helloworld

import (
	"context"

	"cloud.google.com/go/pubsub"
	kitpubsub "github.com/TlexCypher/PubSub/pubsub"
)

type HelloWorldPublisher struct {
	c kitpubsub.PubSubClient
}

func NewHelloWorldPublisher(c kitpubsub.PubSubClient) kitpubsub.Publisher {
	return &HelloWorldPublisher{
		c: c,
	}
}

func (hp *HelloWorldPublisher) Publish(ctx context.Context, topic kitpubsub.Topic, data []byte) (string, error) {
	return "", nil
}

type HelloWorldSubscriber struct {
	c            kitpubsub.PubSubClient
	subscription kitpubsub.Subscription
	handler      kitpubsub.SubscriptionHandler
}

func NewHelloWorldSubscriber(c kitpubsub.PubSubClient, subscription kitpubsub.Subscription, handler kitpubsub.SubscriptionHandler) kitpubsub.Subscriber {
	return &HelloWorldSubscriber{
		c:            c,
		subscription: subscription,
		handler:      handler,
	}
}

type HelloWorldSubscriptionHandler struct{}

func (h *HelloWorldSubscriptionHandler) Handle(ctx context.Context, msg *pubsub.Message) error {
	return nil
}

func (hs *HelloWorldSubscriber) Subscribe(ctx context.Context) error {
	return nil
}
