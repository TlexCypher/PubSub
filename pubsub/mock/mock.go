package mock

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	kitpubsub "github.com/TlexCypher/PubSub/pubsub"
)

var (
	MockServerID            = "MOCK_SERVER"
	MockServerErr           = fmt.Errorf("MOCK_SERVER_ERR")
	MockSubscribeErr        = fmt.Errorf("MOCK_SUBSCRIBE_ERR")
	MockPublishFunc         = (func(context.Context, kitpubsub.Topic, []byte) (string, error))(nil)
	MockSubscriptionHandler = (func(context.Context, *pubsub.Message))(nil)
)

/* publisher */
type MockPublisher struct {
	publishFunc func(context.Context, kitpubsub.Topic, []byte) (string, error)
	serverID    string
	serverErr   error
}

type MockPublisherBuilder struct {
	mockPublisher MockPublisher
}

func NewMockPublisherBuilder() *MockPublisherBuilder {
	return &MockPublisherBuilder{
		mockPublisher: MockPublisher{
			publishFunc: MockPublishFunc,
			serverID:    MockServerID,
			serverErr:   nil,
		},
	}
}

func (mpb *MockPublisherBuilder) PublishFunc(publishFunc func(context.Context, kitpubsub.Topic, []byte) (string, error)) *MockPublisherBuilder {
	mpb.mockPublisher.publishFunc = publishFunc
	return mpb
}

func (mpb *MockPublisherBuilder) ServerID(serverID string) *MockPublisherBuilder {
	mpb.mockPublisher.serverID = serverID
	return mpb
}

func (mpb *MockPublisherBuilder) ServerErr(serverErr error) *MockPublisherBuilder {
	mpb.mockPublisher.serverErr = serverErr
	return mpb
}

func (mpb *MockPublisherBuilder) Build() *MockPublisher {
	return &mpb.mockPublisher
}

func (mp *MockPublisher) Publish(ctx context.Context, topic kitpubsub.Topic, data []byte) (string, error) {
	if mp.publishFunc != nil {
		return mp.publishFunc(ctx, topic, data)
	}
	return mp.serverID, mp.serverErr
}

/* subscriber */
type MockSubscriberBuilder struct {
	mockSubscriber *MockSubscriber
}

func NewMockSubscriberBuilder() *MockSubscriberBuilder {
	return &MockSubscriberBuilder{
		mockSubscriber: &MockSubscriber{
			subscriptionHandler: MockSubscriptionHandler,
			subscribeErr:        nil,
		},
	}
}

func (msb *MockSubscriberBuilder) Build() *MockSubscriber {
	return msb.mockSubscriber
}

func (msb *MockSubscriberBuilder) SubscriptionHandler(subscriptionHandler func(context.Context, *pubsub.Message)) *MockSubscriberBuilder {
	msb.mockSubscriber.subscriptionHandler = subscriptionHandler
	return msb
}

func (msb *MockSubscriberBuilder) SubscribeError(subscribeErr error) *MockSubscriberBuilder {
	msb.mockSubscriber.subscribeErr = subscribeErr
	return msb
}

type MockSubscriber struct {
	subscriptionHandler func(context.Context, *pubsub.Message)
	subscribeErr        error
}

func (ms *MockSubscriber) Subscribe(ctx context.Context, subscription kitpubsub.Subscription) error {
	if ms.subscriptionHandler != nil {
		err := subscription.Receive(ctx, ms.subscriptionHandler)
		if err != nil {
			return err
		}
	}
	return ms.subscribeErr
}

/* subscription*/
type MockSubscription struct {
	subscribeErr error
}

func NewMockSubscription() kitpubsub.Subscription {
	return &MockSubscription{
		subscribeErr: nil,
	}
}

func (ms *MockSubscription) Receive(ctx context.Context, handler func(context.Context, *pubsub.Message)) error {
	if handler != nil {
		handler(ctx, &pubsub.Message{
			Data: []byte("MOCK MESSAGE"),
		})
	}
	return ms.subscribeErr
}
