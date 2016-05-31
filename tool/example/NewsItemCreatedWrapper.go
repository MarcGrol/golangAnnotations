// Generated automatically: do not edit manually

package example

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"code.google.com/p/go-uuid/uuid"
)

func (s *NewsItemCreated) Wrap() (*Envelope, error) {
	envelope := new(Envelope)
	envelope.Uuid = uuid.New()
	envelope.SequenceNumber = 0 // Set later by event-store
	envelope.Timestamp = time.Now()
	envelope.AggregateName = "news"
	envelope.AggregateUid = s.GetUid()
	envelope.EventTypeName = "NewsItemCreated"
	blob, err := json.Marshal(s)
	if err != nil {
		log.Printf("Error marshalling NewsItemCreated payload %+v", err)
		return nil, err
	}
	envelope.EventData = string(blob)

	return envelope, nil
}

func IsNewsItemCreated(envelope *Envelope) bool {
	return envelope.EventTypeName == "NewsItemCreated"
}

func GetIfIsNewsItemCreated(envelop *Envelope) (*NewsItemCreated, bool) {
	if IsNewsItemCreated(envelop) == false {
		return nil, false
	}
	event, err := UnWrapNewsItemCreated(envelop)
	if err != nil {
		return nil, false
	}
	return event, true
}

func UnWrapNewsItemCreated(envelop *Envelope) (*NewsItemCreated, error) {
	if IsNewsItemCreated(envelop) == false {
		return nil, fmt.Errorf("Not a NewsItemCreated")
	}
	var event NewsItemCreated
	err := json.Unmarshal([]byte(envelop.EventData), &event)
	if err != nil {
		log.Printf("Error unmarshalling NewsItemCreated payload %+v", err)
		return nil, err
	}

	return &event, nil
}
