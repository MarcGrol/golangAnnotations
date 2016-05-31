// Generated automatically: do not edit manually

package example

import (
	"time"
)

type Uider interface {
	GetUid() string
}

type Envelope struct {
	Uuid           string
	SequenceNumber uint64
	Timestamp      time.Time
	AggregateName  string
	AggregateUid   string
	EventTypeName  string
	EventData      string
}
