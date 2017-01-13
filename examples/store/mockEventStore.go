package store

import "golang.org/x/net/context"

// EventStore represents an event store
type EventStore interface {
	// Put stores an event related to a specific aggregate-root persistenently in the store
	Put(c context.Context, envelope interface{}) error
}

var (
	// New provides an environment specific implementation of Store
	New func() EventStore
)

func init() {
	New = NewMockEventStore
}

// MockEventStore simulates a real persistent store
type MockEventStore struct {
}

// NewMockEventStore is a factory function that returns our mock store
func NewMockEventStore() EventStore {
	return &MockEventStore{}
}

// Put stores an events that is wrapped in an envelope
func (s *MockEventStore) Put(c context.Context, envelope interface{}) error {
	return nil
}
