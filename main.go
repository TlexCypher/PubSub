package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

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
	ctx := context.TODO()
	worker := kitpubsub.NewPubSubWorker(ctx)
	oc, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("failed to create pubsub client")
	}
	worker.RegisterSubscribers(
		helloworld.NewHelloWorldSubscriber(
			oc,
			kitpubsub.Subscription(helloworldSubscription),
			&helloworld.HelloWorldSubscriptionHandler{},
		),
	)
	worker.Run()

	// main server
	srv := NewApplicationServer(
		&http.Server{
			Addr:    fmt.Sprintf(":%v\n", port),
			Handler: nil,
		},
	)
	srv.Run(map[string]func(http.ResponseWriter, *http.Request){
		"main": func(w http.ResponseWriter, req *http.Request) {
			fmt.Fprint(w, "hello from main server")
		},
	})
}
