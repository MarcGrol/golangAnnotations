
// Generated automatically: do not edit manually

package example

import (
  "encoding/json"
  "log"
  "time"

  "code.google.com/p/go-uuid/uuid"
)

func (s *GamblerCreated) Wrap() (*Envelope,error) {
    var err error
    envelope := new(Envelope)
    envelope.Uuid = uuid.New()
    envelope.SequenceNumber = 0 // Set later by event-store
    envelope.Timestamp = time.Now()
    envelope.AggregateName = "gambler"
    envelope.AggregateUid = s.GetUid()
    envelope.EventTypeName = "GamblerCreated"
    blob, err := json.Marshal(s)
    if err != nil {
        log.Printf("Error marshalling GamblerCreated payload %+v", err)
        return nil, err
    }
    envelope.EventData = string(blob)

    return envelope, nil
}

func IsGamblerCreated(envelope *Envelope) bool {
    return envelope.EventTypeName == "GamblerCreated"
}

func GetIfIsGamblerCreated(envelop *Envelope) (*GamblerCreated, bool) {
    if IsGamblerCreated(envelop) == false {
        return nil, false
    }
    event := UnWrapGamblerCreated(envelop)
    return event, true
}

func UnWrapGamblerCreated(envelop *Envelope) (*GamblerCreated,error) {
    if IsGamblerCreated(envelop) == false {
        return nil
    }
    var event GamblerCreated
    err := json.Unmarshal([]byte(envelop.EventData), &event)
    if err != nil {
        log.Printf("Error unmarshalling GamblerCreated payload %+v", err)
        return nil, err
    }

    return &event, nil
}
