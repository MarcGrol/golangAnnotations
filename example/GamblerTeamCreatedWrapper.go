// Generated automatically: do not edit manually

package example

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/satori/go.uuid"
)

const GamblerTeamCreatedEventName = "GamblerTeamCreated"

func (s *GamblerTeamCreated) Wrap(uid string) (*Envelope, error) {
	envelope := new(Envelope)
	envelope.Uuid = uuid.NewV1().String()
	envelope.SequenceNumber = 0 // Set later by event-store
	envelope.Timestamp = time.Now()
	envelope.AggregateName = GamblerAggregateName
	envelope.AggregateUid = uid
	envelope.EventTypeName = GamblerTeamCreatedEventName
	blob, err := json.Marshal(s)
	if err != nil {
		log.Printf("Error marshalling GamblerTeamCreated payload %+v", err)
		return nil, err
	}
	envelope.EventData = string(blob)

	return envelope, nil
}

func IsGamblerTeamCreated(envelope *Envelope) bool {
	return envelope.EventTypeName == "GamblerTeamCreated"
}

func GetIfIsGamblerTeamCreated(envelop *Envelope) (*GamblerTeamCreated, bool) {
	if IsGamblerTeamCreated(envelop) == false {
		return nil, false
	}
	event, err := UnWrapGamblerTeamCreated(envelop)
	if err != nil {
		return nil, false
	}
	return event, true
}

func UnWrapGamblerTeamCreated(envelop *Envelope) (*GamblerTeamCreated, error) {
	if IsGamblerTeamCreated(envelop) == false {
		return nil, fmt.Errorf("Not a GamblerTeamCreated")
	}
	var event GamblerTeamCreated
	err := json.Unmarshal([]byte(envelop.EventData), &event)
	if err != nil {
		log.Printf("Error unmarshalling GamblerTeamCreated payload %+v", err)
		return nil, err
	}

	return &event, nil
}
