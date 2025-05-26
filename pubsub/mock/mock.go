package mock

import (
	"context"
	"fmt"

	kitpubsub "github.com/TlexCypher/PubSub/pubsub"
)

var (
	MockServerID    = "MOCK_SERVER"
	MockServerErr   = fmt.Errorf("MOCK_SERVER_ERR")
	MockPublishFunc = (func(context.Context, kitpubsub.Topic, []byte) (string, error))(nil)
)

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
