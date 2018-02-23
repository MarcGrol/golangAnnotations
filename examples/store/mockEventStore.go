package store

import (
	"golang.org/x/net/context"

	"github.com/MarcGrol/golangAnnotations/examples/store/request"
)

type EventStore interface {
	Put(c context.Context, rc request.Context, envelope interface{}) error
}

var New = func() EventStore {
	return &MockEventStore{}
}

// MockEventStore simulates a real persistent store
type MockEventStore struct {
}

// Put stores an events that is wrapped in an envelope
func (s *MockEventStore) Put(c context.Context, rc request.Context, envelope interface{}) error {
	return nil
}
