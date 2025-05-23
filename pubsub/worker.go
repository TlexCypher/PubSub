package pubsub

import (
	"context"
	"fmt"
	"log"

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
	eg, ctx := errgroup.WithContext(psw.ctx)
	for _, _sub := range psw.subscribers {
		sub := _sub
		eg.Go(func() error {
			log.Println("PubSubWorker: Starting subscriber...")
			err := sub.Subscribe(ctx)
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
