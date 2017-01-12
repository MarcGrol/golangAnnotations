package store

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/Duxxie/platform/backend/lib/bus"
	"github.com/Duxxie/platform/backend/lib/events"
	"github.com/Duxxie/platform/backend/lib/logging"
)

// StoredItemHandlerFunc is called when iterating over all events
type StoredItemHandlerFunc func(envelope events.Envelope)

// EventStore represents an event store
type EventStore interface {

	// Exists checks if an aggregate exists
	Exists(c context.Context, aggregateName string, aggregateUID string) bool

	GetAllAggregateUIDs(c context.Context, aggregateName string) ([]string, error)

	// Put stores an event related to a specific aggregate-root persistenently in the store
	Put(c context.Context, envelope *events.Envelope) error

	// Iterate fetches all persistent events from the store in order of arrival
	Iterate(c context.Context, callback StoredItemHandlerFunc) error

	// Search fetches all events related to a sinle aggregate-root
	Search(c context.Context, aggregateName string, aggregateUID string) ([]events.Envelope, error)

	// Get fetches a single persisted event from the store
	//Get(c context.Context, UUID string) (*Envelope, error)
}

type eventStoreFactory func() EventStore

var (
	// New provides an environment specific implementation of Store
	New eventStoreFactory
)

var debug = false

func init() {
	New = NewMockEventStore
}

// StoredEnvelopes allow tests to verify events stored
var StoredEnvelopes = []events.Envelope{}
var StoredEnvelopeAggregates = []events.EnvelopeAggregate{}

// MockEventStore simulates a real persistent store
type MockEventStore struct {
	logger logging.Logger
}

// NewMockEventStore is a factory function that returns our mock store
func NewMockEventStore() EventStore {
	return &MockEventStore{
		logger: logging.New(),
	}
}

// Exists checks if an aggregate exists
func (s *MockEventStore) Exists(c context.Context, aggregateName string, aggregateUID string) bool {
	envlps, err := s.Search(c, aggregateName, aggregateUID)
	if err != nil {
		return false
	}
	return len(envlps) > 0
}

func (s *MockEventStore) GetAllAggregateUIDs(c context.Context, aggregateName string) ([]string, error) {
	aggregateUIDs := make([]string, 0)
	for _, aggregate := range StoredEnvelopeAggregates {
		if aggregate.AggregateName == aggregateName {
			aggregateUIDs = append(aggregateUIDs, aggregate.AggregateUID)
		}
	}
	return aggregateUIDs, nil
}

// Put stores an events that is wrapped in an envelope
func (s *MockEventStore) Put(c context.Context, envelope *events.Envelope) error {
	found := false
	for _, env := range StoredEnvelopes {
		if env.UUID == envelope.UUID {
			env = *envelope
			found = true
		}
	}
	if found == false {
		if debug {
			s.logger.Debug(c, "Stored %s event %s.%s: %+v",
				envelope.EventTypeName,
				envelope.AggregateName, envelope.AggregateUID,
				*envelope)
		}
		StoredEnvelopes = append(StoredEnvelopes, *envelope)

		StoredEnvelopeAggregates = append(StoredEnvelopeAggregates, events.EnvelopeAggregate{
			AggregateName: envelope.AggregateName,
			AggregateUID:  envelope.AggregateUID,
		})

		bus.Publish(c, envelope.AggregateName, *envelope)

	}

	return nil
}

func (s *MockEventStore) get(c context.Context, UUID string) (*events.Envelope, error) {
	for _, env := range StoredEnvelopes {
		if env.UUID == UUID {
			if debug {
				s.logger.Debug(c, "Found event: %+v", env)
			}
			return &env, nil
		}
	}

	return nil, fmt.Errorf("Object with UUID %s not found", UUID)
}

// Iterate visits every item in the store
func (s *MockEventStore) Iterate(c context.Context, callback StoredItemHandlerFunc) error {
	idx := 0
	for _, e := range StoredEnvelopes {
		callback(e)
		idx++
	}
	if debug {
		s.logger.Debug(c, "Found %d events", idx)
	}

	return nil
}

// Search looks for events related to a specific aggregate and return them sorted desc on timestamp
func (s *MockEventStore) Search(c context.Context, aggregateName string, aggregateUID string) ([]events.Envelope, error) {

	found := make([]events.Envelope, 0, 10)
	for _, e := range StoredEnvelopes {
		if e.AggregateName == aggregateName && e.AggregateUID == aggregateUID {
			found = append(found, e)
		}
	}
	if debug {
		s.logger.Debug(c, "For %s %s Found %d events: %+v", aggregateName, aggregateUID, len(found), found)
	}

	return found, nil
}
