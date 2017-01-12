package store

import (
	"golang.org/x/net/context"

	"github.com/Duxxie/platform/backend/lib/events"
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
