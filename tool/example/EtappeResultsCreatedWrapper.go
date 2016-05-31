
// Generated automatically: do not edit manually

package example

import (
  "encoding/json"
  "log"
  "time"

  "code.google.com/p/go-uuid/uuid"
)

func (s *EtappeResultsCreated) Wrap() (*Envelope,error) {
    var err error
    envelope := new(Envelope)
    envelope.Uuid = uuid.New()
    envelope.SequenceNumber = 0 // Set later by event-store
    envelope.Timestamp = time.Now()
    envelope.AggregateName = "tour"
    envelope.AggregateUid = s.GetUid()
    envelope.EventTypeName = "EtappeResultsCreated"
    blob, err := json.Marshal(s)
    if err != nil {
        log.Printf("Error marshalling EtappeResultsCreated payload %+v", err)
        return nil, err
    }
    envelope.EventData = string(blob)

    return envelope, nil
}

func IsEtappeResultsCreated(envelope *Envelope) bool {
    return envelope.EventTypeName == "EtappeResultsCreated"
}

func GetIfIsEtappeResultsCreated(envelop *Envelope) (*EtappeResultsCreated, bool) {
    if IsEtappeResultsCreated(envelop) == false {
        return nil, false
    }
    event := UnWrapEtappeResultsCreated(envelop)
    return event, true
}

func UnWrapEtappeResultsCreated(envelop *Envelope) (*EtappeResultsCreated,error) {
    if IsEtappeResultsCreated(envelop) == false {
        return nil
    }
    var event EtappeResultsCreated
    err := json.Unmarshal([]byte(envelop.EventData), &event)
    if err != nil {
        log.Printf("Error unmarshalling EtappeResultsCreated payload %+v", err)
        return nil, err
    }

    return &event, nil
}
