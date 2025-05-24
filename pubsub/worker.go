package pubsub

import (
	"context"
	"fmt"
	"log"
	"maps"

	"golang.org/x/sync/errgroup"
)

type Worker interface {
	Run(context.Context) error
}

type PubSubWorker struct {
	subscriptions map[Subscriber]Subscription
	ctx           context.Context
}

func NewPubSubWorker(ctx context.Context) *PubSubWorker {
	return &PubSubWorker{
		subscriptions: make(map[Subscriber]Subscription, 0),
		ctx:           ctx,
	}
}

func (psw *PubSubWorker) RegisterSubscribers(subscriptionMap map[Subscriber]Subscription) {
	maps.Copy(psw.subscriptions, subscriptionMap)
}

func (psw *PubSubWorker) Run() error {
	eg, ctx := errgroup.WithContext(psw.ctx)
	for subscriber, subscription := range psw.subscriptions {
		sub := subscriber
		eg.Go(func() error {
			log.Println("PubSubWorker: Starting subscriber...")
			err := sub.Subscribe(ctx, subscription)
			if err != nil {
				log.Printf("PubSubWorker: Subscriber exited with error: %v\n", err)
			} else {
				log.Println("PubSubWorker: Subscriber exited gracefully.")
			}
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return fmt.Errorf("pubsub worker run failed: %w", err)
	}
	log.Println("PubSubWorker: All subscribers finished successfully.")
	return nil
}
