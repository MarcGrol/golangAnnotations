
// Generated automatically: do not edit manually

package example

import (
  "encoding/json"
  "log"
  "time"

  "code.google.com/p/go-uuid/uuid"
)

func (s *TourCreated) Wrap() (*Envelope,error) {
    var err error
    envelope := new(Envelope)
    envelope.Uuid = uuid.New()
    envelope.SequenceNumber = 0 // Set later by event-store
    envelope.Timestamp = time.Now()
    envelope.AggregateName = "tour"
    envelope.AggregateUid = s.GetUid()
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
    event := UnWrapTourCreated(envelop)
    return event, true
}

func UnWrapTourCreated(envelop *Envelope) (*TourCreated,error) {
    if IsTourCreated(envelop) == false {
        return nil
    }
    var event TourCreated
    err := json.Unmarshal([]byte(envelop.EventData), &event)
    if err != nil {
        log.Printf("Error unmarshalling TourCreated payload %+v", err)
        return nil, err
    }

    return &event, nil
}
