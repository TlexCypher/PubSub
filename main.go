package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/pubsub"

	kitpubsub "github.com/TlexCypher/PubSub/pubsub"
	"github.com/TlexCypher/PubSub/pubsub/helloworld"
)

var (
	projectID              = os.Getenv("PROJECT_ID")
	topicID                = os.Getenv("TOPIC_ID")
	helloworldSubscription = os.Getenv("HELLOWORLD_SUBSCRIPTION")
	msg                    = "HELLO WORLD FROM PUB/SUB"
)

func makePubSubHandler(publisher kitpubsub.Publisher) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		serverID, err := publisher.Publish(req.Context(), kitpubsub.Topic(topicID), []byte(msg))
		if err != nil {
			log.Printf("failed to publish message: %v", err)
			return
		}
		log.Printf("message was successfully published. serverID:%v\n", serverID)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", serverID)
	}
}

func main() {
	// pub/sub
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	sdkClient, err := pubsub.NewClient(ctx, projectID)
	pubsubClient := kitpubsub.NewPubSubClientImpl(sdkClient)
	if err != nil {
		log.Fatalf("failed to create pubsub client: %v", err)
	}
	defer sdkClient.Close()

	worker := kitpubsub.NewPubSubWorker(ctx)
	worker.RegisterSubscribers(
		map[kitpubsub.Subscriber]kitpubsub.Subscription{
			helloworld.NewHelloWorldSubscriber(
				pubsubClient,
				kitpubsub.SubscriptionID(helloworldSubscription),
			): pubsubClient.Subscription(kitpubsub.SubscriptionID(helloworldSubscription)),
		},
	)
	workerErr, mainServerErr := make(chan error, 1), make(chan error, 1)
	go func() {
		defer close(workerErr)
		log.Println("Starting PubSub workers....")
		workerErr <- worker.Run()
	}()

	// main server
	srv := NewApplicationServer()
	go func() {
		defer close(mainServerErr)
		log.Println("Starting Main HTTP Server...")
		if err := srv.Run(ctx, map[string]func(http.ResponseWriter, *http.Request){
			"/main": makePubSubHandler(helloworld.NewHelloWorldPublisher(pubsubClient)),
		}); err != nil {
			mainServerErr <- fmt.Errorf("HTTP server Run failed: %w", err)
		}
	}()

	select {
	case err := <-workerErr:
		if err != nil {
			log.Printf("PubSub worker exited with error: %v\n", err)
		} else {
			log.Println("PubSub workers were shutdowned gracefully ")
		}
		stop()
	case err := <-mainServerErr:
		if err != nil {
			log.Printf("Main HTTP server exited with error: %v\n", err)
		} else {
			log.Println("Main HTTP Server were shutdowned gracefully ")
		}
		stop()
	case <-ctx.Done():
		log.Println("Shutdown signal received, initiating graceful shutdown...")
	}

	log.Println("Waiting for services to complete shutdown...")
	<-workerErr
	<-mainServerErr
	log.Println("Application shutdowned.")
}
