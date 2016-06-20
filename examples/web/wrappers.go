// Generated automatically: do not edit manually

package example

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/satori/go.uuid"
)

type Envelope struct {
	Uuid           string
	SequenceNumber uint64
	Timestamp      time.Time
	AggregateName  string
	AggregateUid   string
	EventTypeName  string
	EventData      string
}

const (
	TourServiceEventName = "TourService"
)

func (s *TourService) Wrap(uid string) (*Envelope, error) {
	envelope := new(Envelope)
	envelope.Uuid = uuid.NewV1().String()
	envelope.SequenceNumber = 0 // Set later by event-store
	envelope.Timestamp = time.Now()
	envelope.AggregateName = AggregateName // from annotation!
	envelope.AggregateUid = uid
	envelope.EventTypeName = TourServiceEventName
	blob, err := json.Marshal(s)
	if err != nil {
		log.Printf("Error marshalling TourService payload %+v", err)
		return nil, err
	}
	envelope.EventData = string(blob)

	return envelope, nil
}

func IsTourService(envelope *Envelope) bool {
	return envelope.EventTypeName == TourServiceEventName
}

func GetIfIsTourService(envelop *Envelope) (*TourService, bool) {
	if IsTourService(envelop) == false {
		return nil, false
	}
	event, err := UnWrapTourService(envelop)
	if err != nil {
		return nil, false
	}
	return event, true
}

func UnWrapTourService(envelop *Envelope) (*TourService, error) {
	if IsTourService(envelop) == false {
		return nil, fmt.Errorf("Not a TourService")
	}
	var event TourService
	err := json.Unmarshal([]byte(envelop.EventData), &event)
	if err != nil {
		log.Printf("Error unmarshalling TourService payload %+v", err)
		return nil, err
	}

	return &event, nil
}
