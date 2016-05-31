
// Generated automatically: do not edit manually

package generator

import (
  "encoding/json"
  "log"
  "time"

  "code.google.com/p/go-uuid/uuid"
)

func (s *MyStruct) Wrap() *Envelope {
    var err error
    envelope := new(Envelope)
    envelope.Uuid = uuid.New()
    envelope.SequenceNumber = 0 // Set later by event-store
    envelope.Timestamp = time.Now()
    envelope.AggregateName = "person"
    envelope.AggregateUid = s.GetUid()
    envelope.EventTypeName = "MyStruct"
    blob, err := json.Marshal(s)
    if err != nil {
        log.Printf("Error marshalling MyStruct payload %+v", err)
        return nil //, err
    }
    envelope.EventData = string(blob)

    return envelope
}

func IsMyStruct(envelope *Envelope) bool {
    return envelope.EventTypeName == "MyStruct"
}

func GetIfIsMyStruct(envelop *Envelope) (*MyStruct, bool) {
    if IsMyStruct(envelop) == false {
        return nil, false
    }
    event := UnWrapMyStruct(envelop)
    return event, true
}

func UnWrapMyStruct(envelop *Envelope) *MyStruct {
    if IsMyStruct(envelop) == false {
        return nil
    }
    var event MyStruct
    err := json.Unmarshal([]byte(envelop.EventData), &event)
    if err != nil {
        log.Printf("Error unmarshalling MyStruct payload %+v", err)
        return nil
    }

    return &event
}
