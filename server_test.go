package main

import (
	"net/http/httptest"
	"testing"

	"github.com/TlexCypher/PubSub/pubsub/mock"
	"github.com/google/go-cmp/cmp"
)

func TestApplicationServer(t *testing.T) {
	t.Parallel()
	tests := []struct {
		description string
		want        string
		wantErr     error
	}{
		{
			description: "Success case",
			want:        mock.MockServerID,
			wantErr:     nil,
		},
		// {
		// 	description: "Failed case",
		// 	want:        "",
		// 	wantErr:     mock.MockServerErr,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			tt := tt
			mockPub := mock.NewMockPublisherBuilder().Build()
			handler := makePubSubHandler(mockPub)
			req := httptest.NewRequest("GET", "/main", nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			// body assertion
			if diff := cmp.Diff(tt.want, rr.Body.String()); diff != "" {
				t.Errorf("Application Server result diff (-expect +got)\n%s", diff)
			}
		})
	}
}
