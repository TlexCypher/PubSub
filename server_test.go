package main

import (
	"context"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"cloud.google.com/go/pubsub"
	kitpubsub "github.com/TlexCypher/PubSub/pubsub"
	"github.com/TlexCypher/PubSub/pubsub/mock"
	"github.com/google/go-cmp/cmp"
)

/*
NOTE: Entire description for unit tests

We have two main processes, Application Server and PubSub-Worker.

Application server has only one endpoint, named '/main'.
When that endpoint is hit, message is published to PubSubWorker.

So, we need to implement two test cases, TestApplicationServer, TestPubSubWorker.
*/

func TestApplicationServer(t *testing.T) {
	t.Parallel()
	tests := []struct {
		description     string
		wantBody        string
		wantPublishFunc func(context.Context, kitpubsub.Topic, []byte) (string, error)
		wantServerID    string
		wantErr         error
	}{
		{
			description:     "Success case",
			wantBody:        mock.MockServerID,
			wantPublishFunc: nil,
			wantServerID:    mock.MockServerID,
			wantErr:         nil,
		},
		{
			description: "Failed case",
			wantBody:    "",
			wantPublishFunc: func(context.Context, kitpubsub.Topic, []byte) (string, error) {
				return "", fmt.Errorf("failed to publish")
			},
			wantServerID: "",
			wantErr:      fmt.Errorf("failed to publish"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			tt := tt
			mockPub := mock.NewMockPublisherBuilder().PublishFunc(tt.wantPublishFunc).ServerID(tt.wantServerID).ServerErr(tt.wantErr).Build()
			handler := makePubSubHandler(mockPub)
			req := httptest.NewRequest("GET", "/main", nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			// body assertion
			if diff := cmp.Diff(tt.wantBody, rr.Body.String()); diff != "" {
				t.Errorf("Application Server result diff (-expect +got)\n%s", diff)
			}
			// err assertion
		})
	}
}

func TestPubSubWorker(t *testing.T) {
	t.Parallel()
	tests := []struct {
		description         string
		wantErr             error
		subscriptionHandler func(context.Context, *pubsub.Message)
		subscribeErr        error
	}{
		{
			description:         "Success case",
			wantErr:             nil,
			subscriptionHandler: nil,
			subscribeErr:        nil,
		},
		{
			description:         "Failed case",
			wantErr:             mock.MockServerErr,
			subscriptionHandler: nil,
			subscribeErr:        mock.MockServerErr,
		},
	}
	ctx := context.Background()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			mockedPubsubWorker := kitpubsub.NewPubSubWorker(ctx)
			builder := mock.NewMockSubscriberBuilder().SubscriptionHandler(tt.subscriptionHandler).SubscribeError(tt.subscribeErr)
			mockedPubsubWorker.RegisterSubscribers(
				map[kitpubsub.Subscriber]kitpubsub.Subscription{
					builder.Build(): mock.NewMockSubscription(),
				},
			)
			if err := mockedPubsubWorker.Run(); !errors.Is(err, tt.wantErr) {
				t.Errorf("failed to run pubsub workers: (-expected, +got)\n-%v\n+%v\n", tt.wantErr, err)
			}
		})
	}
}
