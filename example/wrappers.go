// Generated automatically: do not edit manually

package example

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

const TourCreatedEventName = "TourCreated"

func (s *TourCreated) Wrap(uid string) (*Envelope, error) {
    envelope := new(Envelope)
    envelope.Uuid = uuid.NewV1().String()
    envelope.SequenceNumber = 0 // Set later by event-store
    envelope.Timestamp = time.Now()
    envelope.AggregateName = TourAggregateName
    envelope.AggregateUid = uid
    envelope.EventTypeName = TourCreatedEventName
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
    event, err := UnWrapTourCreated(envelop)
    if err != nil {
        return nil, false
    }
    return event, true
}

func UnWrapTourCreated(envelop *Envelope) (*TourCreated, error) {
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

const CyclistCreatedEventName = "CyclistCreated"

func (s *CyclistCreated) Wrap(uid string) (*Envelope, error) {
    envelope := new(Envelope)
    envelope.Uuid = uuid.NewV1().String()
    envelope.SequenceNumber = 0 // Set later by event-store
    envelope.Timestamp = time.Now()
    envelope.AggregateName = TourAggregateName
    envelope.AggregateUid = uid
    envelope.EventTypeName = CyclistCreatedEventName
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

const EtappeCreatedEventName = "EtappeCreated"

func (s *EtappeCreated) Wrap(uid string) (*Envelope, error) {
    envelope := new(Envelope)
    envelope.Uuid = uuid.NewV1().String()
    envelope.SequenceNumber = 0 // Set later by event-store
    envelope.Timestamp = time.Now()
    envelope.AggregateName = TourAggregateName
    envelope.AggregateUid = uid
    envelope.EventTypeName = EtappeCreatedEventName
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
    event, err := UnWrapEtappeCreated(envelop)
    if err != nil {
        return nil, false
    }
    return event, true
}

func UnWrapEtappeCreated(envelop *Envelope) (*EtappeCreated, error) {
    if IsEtappeCreated(envelop) == false {
        return nil, fmt.Errorf("Not a EtappeCreated")
    }
    var event EtappeCreated
    err := json.Unmarshal([]byte(envelop.EventData), &event)
    if err != nil {
        log.Printf("Error unmarshalling EtappeCreated payload %+v", err)
        return nil, err
    }

    return &event, nil
}

const EtappeResultsCreatedEventName = "EtappeResultsCreated"

func (s *EtappeResultsCreated) Wrap(uid string) (*Envelope, error) {
    envelope := new(Envelope)
    envelope.Uuid = uuid.NewV1().String()
    envelope.SequenceNumber = 0 // Set later by event-store
    envelope.Timestamp = time.Now()
    envelope.AggregateName = TourAggregateName
    envelope.AggregateUid = uid
    envelope.EventTypeName = EtappeResultsCreatedEventName
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
    event, err := UnWrapEtappeResultsCreated(envelop)
    if err != nil {
        return nil, false
    }
    return event, true
}

func UnWrapEtappeResultsCreated(envelop *Envelope) (*EtappeResultsCreated, error) {
    if IsEtappeResultsCreated(envelop) == false {
        return nil, fmt.Errorf("Not a EtappeResultsCreated")
    }
    var event EtappeResultsCreated
    err := json.Unmarshal([]byte(envelop.EventData), &event)
    if err != nil {
        log.Printf("Error unmarshalling EtappeResultsCreated payload %+v", err)
        return nil, err
    }

    return &event, nil
}

const GamblerCreatedEventName = "GamblerCreated"

func (s *GamblerCreated) Wrap(uid string) (*Envelope, error) {
    envelope := new(Envelope)
    envelope.Uuid = uuid.NewV1().String()
    envelope.SequenceNumber = 0 // Set later by event-store
    envelope.Timestamp = time.Now()
    envelope.AggregateName = GamblerAggregateName
    envelope.AggregateUid = uid
    envelope.EventTypeName = GamblerCreatedEventName
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
    event, err := UnWrapGamblerCreated(envelop)
    if err != nil {
        return nil, false
    }
    return event, true
}

func UnWrapGamblerCreated(envelop *Envelope) (*GamblerCreated, error) {
    if IsGamblerCreated(envelop) == false {
        return nil, fmt.Errorf("Not a GamblerCreated")
    }
    var event GamblerCreated
    err := json.Unmarshal([]byte(envelop.EventData), &event)
    if err != nil {
        log.Printf("Error unmarshalling GamblerCreated payload %+v", err)
        return nil, err
    }

    return &event, nil
}

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

const NewsItemCreatedEventName = "NewsItemCreated"

func (s *NewsItemCreated) Wrap(uid string) (*Envelope, error) {
    envelope := new(Envelope)
    envelope.Uuid = uuid.NewV1().String()
    envelope.SequenceNumber = 0 // Set later by event-store
    envelope.Timestamp = time.Now()
    envelope.AggregateName = NewsAggregateName
    envelope.AggregateUid = uid
    envelope.EventTypeName = NewsItemCreatedEventName
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
