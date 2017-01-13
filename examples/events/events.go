package events

import "time"

// Envelope wraps all types events so it can be easily serialized and stored
type Envelope struct {
	// UUID globally unique identifier if this event
	UUID string `datastore:",noindex"`
	// IsRootEvent determines if this events is the first in a sequence, all related to the same aggregate-instance
	IsRootEvent bool `datastore:",noindex"`
	// SessionUID keeps tracj in which session the event was created. This gives us auditing capabilities.
	SessionUID string
	// SequenceNumber currently not used. We use Timestamp to order events over time
	SequenceNumber int64
	// Timestamp is used to order events over time
	Timestamp time.Time
	// AggregateName is the name of the aggregate
	AggregateName string
	// AggregateUID is the unique identifier of the aggregate
	AggregateUID string
	// EventTypeName is the name of the event type
	EventTypeName string
	// EventTypeVersion allows for changing event in future
	EventTypeVersion int
	// EventData is the serialized payload of the event
	EventData string `datastore:",noindex"`
	// CheckSum is based on all data in Envelope and can indicate data has been tampered with
	CheckSum []byte `datastore:",noindex"`
}

type Event interface {
	GetUID() string
}
