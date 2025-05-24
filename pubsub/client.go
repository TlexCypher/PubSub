package pubsub

import "cloud.google.com/go/pubsub"

type PubSubClientImpl struct {
	c *pubsub.Client
}

func NewPubSubClientImpl(c *pubsub.Client) *PubSubClientImpl {
	return &PubSubClientImpl{
		c: c,
	}
}

func (c *PubSubClientImpl) GetClient() *pubsub.Client {
	return c.c
}

func (c *PubSubClientImpl) Close() error {
	return c.c.Close()
}
