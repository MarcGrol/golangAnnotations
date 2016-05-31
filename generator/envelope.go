// Generated automatically: do not edit manually

package generator

import "time"

type Envelope struct {
	Uuid           string
	SequenceNumber uint64
	Timestamp      time.Time
	AggregateName  string
	AggregateUid   string
	EventTypeName  string
	EventData      string
}
