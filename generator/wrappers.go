// Generated automatically: do not edit manually

package generator

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
	MyStructEventName = "MyStruct"
)

func (s *MyStruct) Wrap(uid string) (*Envelope, error) {
	envelope := new(Envelope)
	envelope.Uuid = uuid.NewV1().String()
	envelope.SequenceNumber = 0 // Set later by event-store
	envelope.Timestamp = time.Now()
	envelope.AggregateName = PersonAggregateName // from annotation!
	envelope.AggregateUid = uid
	envelope.EventTypeName = MyStructEventName
	blob, err := json.Marshal(s)
	if err != nil {
		log.Printf("Error marshalling MyStruct payload %+v", err)
		return nil, err
	}
	envelope.EventData = string(blob)

	return envelope, nil
}

func IsMyStruct(envelope *Envelope) bool {
	return envelope.EventTypeName == MyStructEventName
}

func GetIfIsMyStruct(envelop *Envelope) (*MyStruct, bool) {
	if IsMyStruct(envelop) == false {
		return nil, false
	}
	event, err := UnWrapMyStruct(envelop)
	if err != nil {
		return nil, false
	}
	return event, true
}

func UnWrapMyStruct(envelop *Envelope) (*MyStruct, error) {
	if IsMyStruct(envelop) == false {
		return nil, fmt.Errorf("Not a MyStruct")
	}
	var event MyStruct
	err := json.Unmarshal([]byte(envelop.EventData), &event)
	if err != nil {
		log.Printf("Error unmarshalling MyStruct payload %+v", err)
		return nil, err
	}

	return &event, nil
}
