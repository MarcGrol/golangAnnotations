// Generated automatically: do not edit manually

package example

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/satori/go.uuid"
)

func (s *CyclistCreated) Wrap(uid string) (*Envelope, error) {
	envelope := new(Envelope)
	envelope.Uuid = uuid.New()
	envelope.SequenceNumber = 0 // Set later by event-store
	envelope.Timestamp = time.Now()
	envelope.AggregateName = "tour"
	envelope.AggregateUid = uid
	envelope.EventTypeName = "CyclistCreated"
	blob, err := json.Marshal(s)
	if err != nil {
		log.Printf("Error marshalling CyclistCreated payload %+v", err)
		return nil, err
	}
	envelope.EventData = string(blob)

	return envelope, nil
}

func IsCyclistCreated(envelope *Envelope) bool {
	return envelope.EventTypeName == "CyclistCreated"
}

func GetIfIsCyclistCreated(envelop *Envelope) (*CyclistCreated, bool) {
	if IsCyclistCreated(envelop) == false {
		return nil, false
	}
	event, err := UnWrapCyclistCreated(envelop)
	if err != nil {
		return nil, false
	}
	return event, true
}

func UnWrapCyclistCreated(envelop *Envelope) (*CyclistCreated, error) {
	if IsCyclistCreated(envelop) == false {
		return nil, fmt.Errorf("Not a CyclistCreated")
	}
	var event CyclistCreated
	err := json.Unmarshal([]byte(envelop.EventData), &event)
	if err != nil {
		log.Printf("Error unmarshalling CyclistCreated payload %+v", err)
		return nil, err
	}

	return &event, nil
}
