package pubsub

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Worker interface {
	Run(context.Context) error
}

type PubSubWorker struct {
	subscribers []Subscriber
	ctx         context.Context
}

func NewPubSubWorker(ctx context.Context) *PubSubWorker {
	return &PubSubWorker{
		subscribers: make([]Subscriber, 0),
		ctx:         ctx,
	}
}

func (psw *PubSubWorker) RegisterSubscribers(subscribers ...Subscriber) {
	for _, sub := range subscribers {
		psw.subscribers = append(psw.subscribers, sub)
	}
}

func (psw *PubSubWorker) Run() error {
	eg := errgroup.Group{}
	eg.Go(func() error {
		for _, sub := range psw.subscribers {
			if err := sub.Subscribe(
				psw.ctx,
				sub.GetSubscription(),
				sub.GetSubscriptionHandler(),
			); err != nil {
				return err
			}
		}
		return nil
	})
	return nil
}
