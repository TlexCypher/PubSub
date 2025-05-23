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
	port                   = os.Getenv("PORT")
	projectID              = os.Getenv("PROJECT_ID")
	helloworldSubscription = os.Getenv("HELLOWORLD_SUBSCRIPTION")
)

func main() {
	// pub/sub
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	sdkClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("failed to create pubsub client: %v", err)
	}
	defer sdkClient.Close()

	worker := kitpubsub.NewPubSubWorker(ctx)
	worker.RegisterSubscribers(
		helloworld.NewHelloWorldSubscriber(
			sdkClient,
			kitpubsub.Subscription(helloworldSubscription),
			&helloworld.HelloWorldSubscriptionHandler{},
		),
	)

	workerErr, mainServerErr := make(chan error, 1), make(chan error, 1)
	go func() {
		defer close(workerErr)
		log.Println("Starting PubSub workers....")
		workerErr <- worker.Run()
	}()

	// main server
	srv := NewApplicationServer(
		&http.Server{
			Addr:    fmt.Sprintf(":%v\n", port),
			Handler: nil,
		},
	)

	go func() {
		defer close(mainServerErr)
		log.Println("Starting Main HTTP Server...")
		if err := srv.Run(ctx, map[string]func(http.ResponseWriter, *http.Request){
			"main": func(w http.ResponseWriter, req *http.Request) {
				fmt.Fprint(w, "hello from main server")
			},
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
