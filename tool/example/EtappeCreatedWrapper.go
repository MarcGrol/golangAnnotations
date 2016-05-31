
// Generated automatically: do not edit manually

package example

import (
  "encoding/json"
  "log"
  "time"

  "code.google.com/p/go-uuid/uuid"
)

func (s *EtappeCreated) Wrap() (*Envelope,error) {
    var err error
    envelope := new(Envelope)
    envelope.Uuid = uuid.New()
    envelope.SequenceNumber = 0 // Set later by event-store
    envelope.Timestamp = time.Now()
    envelope.AggregateName = "tour"
    envelope.AggregateUid = s.GetUid()
    envelope.EventTypeName = "EtappeCreated"
    blob, err := json.Marshal(s)
    if err != nil {
        log.Printf("Error marshalling EtappeCreated payload %+v", err)
        return nil, err
    }
    envelope.EventData = string(blob)

    return envelope, nil
}

func IsEtappeCreated(envelope *Envelope) bool {
    return envelope.EventTypeName == "EtappeCreated"
}

func GetIfIsEtappeCreated(envelop *Envelope) (*EtappeCreated, bool) {
    if IsEtappeCreated(envelop) == false {
        return nil, false
    }
    event := UnWrapEtappeCreated(envelop)
    return event, true
}

func UnWrapEtappeCreated(envelop *Envelope) (*EtappeCreated,error) {
    if IsEtappeCreated(envelop) == false {
        return nil
    }
    var event EtappeCreated
    err := json.Unmarshal([]byte(envelop.EventData), &event)
    if err != nil {
        log.Printf("Error unmarshalling EtappeCreated payload %+v", err)
        return nil, err
    }

    return &event, nil
}
