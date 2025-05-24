package helloworld

import (
	"context"
	"fmt"
	"log"

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
	pubsubClient := hp.c.GetClient()
	t := pubsubClient.Topic(string(topic))
	result := t.Publish(ctx, &pubsub.Message{
		Data: data,
	})
	return result.Get(ctx)
}

type HelloWorldSubscriber struct {
	c            kitpubsub.PubSubClient
	subscription kitpubsub.Subscription
}

func NewHelloWorldSubscriber(c kitpubsub.PubSubClient, subscription kitpubsub.Subscription) kitpubsub.Subscriber {
	return &HelloWorldSubscriber{
		c:            c,
		subscription: subscription,
	}
}

func (hs *HelloWorldSubscriber) Subscribe(ctx context.Context, subscription kitpubsub.Subscription) error {
	pubsubClient := hs.c.GetClient()
	sub := pubsubClient.Subscription(string(subscription))

	err := sub.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		log.Printf("Got message: %q\n", string(msg.Data))
		msg.Ack()
	})
	if err != nil {
		return fmt.Errorf("sub.Receive: %w", err)
	}
	return nil
}
