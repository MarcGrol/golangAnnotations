
// Generated automatically: do not edit manually

package example

import (
  "encoding/json"
  "fmt"
  "log"
  "time"

  "code.google.com/p/go-uuid/uuid"
)

func (s *TourCreated) Wrap(uid string) (*Envelope,error) {
    envelope := new(Envelope)
    envelope.Uuid = uuid.New()
    envelope.SequenceNumber = 0 // Set later by event-store
    envelope.Timestamp = time.Now()
    envelope.AggregateName = "tour"
    envelope.AggregateUid = uid
    envelope.EventTypeName = "TourCreated"
    blob, err := json.Marshal(s)
    if err != nil {
        log.Printf("Error marshalling TourCreated payload %+v", err)
        return nil, err
    }
    envelope.EventData = string(blob)

    return envelope, nil
}

func IsTourCreated(envelope *Envelope) bool {
    return envelope.EventTypeName == "TourCreated"
}

func GetIfIsTourCreated(envelop *Envelope) (*TourCreated, bool) {
    if IsTourCreated(envelop) == false {
        return nil, false
    }
    event,err := UnWrapTourCreated(envelop)
    if err != nil {
    	return nil, false
    }
    return event, true
}

func UnWrapTourCreated(envelop *Envelope) (*TourCreated,error) {
    if IsTourCreated(envelop) == false {
        return nil, fmt.Errorf("Not a TourCreated")
    }
    var event TourCreated
    err := json.Unmarshal([]byte(envelop.EventData), &event)
    if err != nil {
        log.Printf("Error unmarshalling TourCreated payload %+v", err)
        return nil, err
    }

    return &event, nil
}
