// Generated automatically by golangAnnotations: do not edit manually

package event

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Duxxie/platform/backend/lib/events"
	"github.com/Duxxie/platform/backend/lib/mytime"
	uuid "github.com/satori/go.uuid"
)

const (

	// TourCreatedEventName provides a constant symbol for TourCreated
	TourCreatedEventName = "TourCreated"

	// CyclistCreatedEventName provides a constant symbol for CyclistCreated
	CyclistCreatedEventName = "CyclistCreated"

	// EtappeCreatedEventName provides a constant symbol for EtappeCreated
	EtappeCreatedEventName = "EtappeCreated"

	// EtappeResultsCreatedEventName provides a constant symbol for EtappeResultsCreated
	EtappeResultsCreatedEventName = "EtappeResultsCreated"

	// GamblerCreatedEventName provides a constant symbol for GamblerCreated
	GamblerCreatedEventName = "GamblerCreated"

	// GamblerTeamCreatedEventName provides a constant symbol for GamblerTeamCreated
	GamblerTeamCreatedEventName = "GamblerTeamCreated"

	// NewsItemCreatedEventName provides a constant symbol for NewsItemCreated
	NewsItemCreatedEventName = "NewsItemCreated"
)

var getUID = func() string {
	return uuid.NewV1().String()
}

// Wrap wraps event TourCreated into an envelope
func (s *TourCreated) Wrap(sessionUID string) (*events.Envelope, error) {
	blob, err := json.Marshal(s)
	if err != nil {
		log.Printf("Error marshalling TourCreated payload %+v", err)
		return nil, err
	}
	envelope := events.Envelope{
		UUID:             getUID(),
		IsRootEvent:      false,
		SequenceNumber:   int64(0), // Set later by event-store
		SessionUID:       sessionUID,
		Timestamp:        mytime.Now(),
		AggregateName:    TourAggregateName, // from annotation!
		AggregateUID:     s.GetUID(),
		EventTypeName:    TourCreatedEventName,
		EventTypeVersion: 0,
		EventData:        string(blob),
	}

	return &envelope, nil
}

// IsTourCreated detects of envelope carries event of type TourCreated
func IsTourCreated(envelope *events.Envelope) bool {
	return envelope.EventTypeName == TourCreatedEventName
}

// GetIfIsTourCreated detects of envelope carries event of type TourCreated and returns the event if so
func GetIfIsTourCreated(envelope *events.Envelope) (*TourCreated, bool) {
	if IsTourCreated(envelope) == false {
		return nil, false
	}
	event, err := UnWrapTourCreated(envelope)
	if err != nil {
		return nil, false
	}
	return event, true
}

// UnWrapTourCreated extracts event TourCreated from its envelope
func UnWrapTourCreated(envelope *events.Envelope) (*TourCreated, error) {
	if IsTourCreated(envelope) == false {
		return nil, fmt.Errorf("Not a TourCreated")
	}
	var event TourCreated
	err := json.Unmarshal([]byte(envelope.EventData), &event)
	if err != nil {
		log.Printf("Error unmarshalling TourCreated payload %+v", err)
		return nil, err
	}
	event.Timestamp = envelope.Timestamp.In(mytime.DutchLocation())

	return &event, nil
}

// Wrap wraps event CyclistCreated into an envelope
func (s *CyclistCreated) Wrap(sessionUID string) (*events.Envelope, error) {
	blob, err := json.Marshal(s)
	if err != nil {
		log.Printf("Error marshalling CyclistCreated payload %+v", err)
		return nil, err
	}
	envelope := events.Envelope{
		UUID:             getUID(),
		IsRootEvent:      false,
		SequenceNumber:   int64(0), // Set later by event-store
		SessionUID:       sessionUID,
		Timestamp:        mytime.Now(),
		AggregateName:    TourAggregateName, // from annotation!
		AggregateUID:     s.GetUID(),
		EventTypeName:    CyclistCreatedEventName,
		EventTypeVersion: 0,
		EventData:        string(blob),
	}

	return &envelope, nil
}

// IsCyclistCreated detects of envelope carries event of type CyclistCreated
func IsCyclistCreated(envelope *events.Envelope) bool {
	return envelope.EventTypeName == CyclistCreatedEventName
}

// GetIfIsCyclistCreated detects of envelope carries event of type CyclistCreated and returns the event if so
func GetIfIsCyclistCreated(envelope *events.Envelope) (*CyclistCreated, bool) {
	if IsCyclistCreated(envelope) == false {
		return nil, false
	}
	event, err := UnWrapCyclistCreated(envelope)
	if err != nil {
		return nil, false
	}
	return event, true
}

// UnWrapCyclistCreated extracts event CyclistCreated from its envelope
func UnWrapCyclistCreated(envelope *events.Envelope) (*CyclistCreated, error) {
	if IsCyclistCreated(envelope) == false {
		return nil, fmt.Errorf("Not a CyclistCreated")
	}
	var event CyclistCreated
	err := json.Unmarshal([]byte(envelope.EventData), &event)
	if err != nil {
		log.Printf("Error unmarshalling CyclistCreated payload %+v", err)
		return nil, err
	}
	event.Timestamp = envelope.Timestamp.In(mytime.DutchLocation())

	return &event, nil
}

// Wrap wraps event EtappeCreated into an envelope
func (s *EtappeCreated) Wrap(sessionUID string) (*events.Envelope, error) {
	blob, err := json.Marshal(s)
	if err != nil {
		log.Printf("Error marshalling EtappeCreated payload %+v", err)
		return nil, err
	}
	envelope := events.Envelope{
		UUID:             getUID(),
		IsRootEvent:      false,
		SequenceNumber:   int64(0), // Set later by event-store
		SessionUID:       sessionUID,
		Timestamp:        mytime.Now(),
		AggregateName:    TourAggregateName, // from annotation!
		AggregateUID:     s.GetUID(),
		EventTypeName:    EtappeCreatedEventName,
		EventTypeVersion: 0,
		EventData:        string(blob),
	}

	return &envelope, nil
}

// IsEtappeCreated detects of envelope carries event of type EtappeCreated
func IsEtappeCreated(envelope *events.Envelope) bool {
	return envelope.EventTypeName == EtappeCreatedEventName
}

// GetIfIsEtappeCreated detects of envelope carries event of type EtappeCreated and returns the event if so
func GetIfIsEtappeCreated(envelope *events.Envelope) (*EtappeCreated, bool) {
	if IsEtappeCreated(envelope) == false {
		return nil, false
	}
	event, err := UnWrapEtappeCreated(envelope)
	if err != nil {
		return nil, false
	}
	return event, true
}

// UnWrapEtappeCreated extracts event EtappeCreated from its envelope
func UnWrapEtappeCreated(envelope *events.Envelope) (*EtappeCreated, error) {
	if IsEtappeCreated(envelope) == false {
		return nil, fmt.Errorf("Not a EtappeCreated")
	}
	var event EtappeCreated
	err := json.Unmarshal([]byte(envelope.EventData), &event)
	if err != nil {
		log.Printf("Error unmarshalling EtappeCreated payload %+v", err)
		return nil, err
	}
	event.Timestamp = envelope.Timestamp.In(mytime.DutchLocation())

	return &event, nil
}

// Wrap wraps event EtappeResultsCreated into an envelope
func (s *EtappeResultsCreated) Wrap(sessionUID string) (*events.Envelope, error) {
	blob, err := json.Marshal(s)
	if err != nil {
		log.Printf("Error marshalling EtappeResultsCreated payload %+v", err)
		return nil, err
	}
	envelope := events.Envelope{
		UUID:             getUID(),
		IsRootEvent:      false,
		SequenceNumber:   int64(0), // Set later by event-store
		SessionUID:       sessionUID,
		Timestamp:        mytime.Now(),
		AggregateName:    TourAggregateName, // from annotation!
		AggregateUID:     s.GetUID(),
		EventTypeName:    EtappeResultsCreatedEventName,
		EventTypeVersion: 0,
		EventData:        string(blob),
	}

	return &envelope, nil
}

// IsEtappeResultsCreated detects of envelope carries event of type EtappeResultsCreated
func IsEtappeResultsCreated(envelope *events.Envelope) bool {
	return envelope.EventTypeName == EtappeResultsCreatedEventName
}

// GetIfIsEtappeResultsCreated detects of envelope carries event of type EtappeResultsCreated and returns the event if so
func GetIfIsEtappeResultsCreated(envelope *events.Envelope) (*EtappeResultsCreated, bool) {
	if IsEtappeResultsCreated(envelope) == false {
		return nil, false
	}
	event, err := UnWrapEtappeResultsCreated(envelope)
	if err != nil {
		return nil, false
	}
	return event, true
}

// UnWrapEtappeResultsCreated extracts event EtappeResultsCreated from its envelope
func UnWrapEtappeResultsCreated(envelope *events.Envelope) (*EtappeResultsCreated, error) {
	if IsEtappeResultsCreated(envelope) == false {
		return nil, fmt.Errorf("Not a EtappeResultsCreated")
	}
	var event EtappeResultsCreated
	err := json.Unmarshal([]byte(envelope.EventData), &event)
	if err != nil {
		log.Printf("Error unmarshalling EtappeResultsCreated payload %+v", err)
		return nil, err
	}
	event.Timestamp = envelope.Timestamp.In(mytime.DutchLocation())

	return &event, nil
}

// Wrap wraps event GamblerCreated into an envelope
func (s *GamblerCreated) Wrap(sessionUID string) (*events.Envelope, error) {
	blob, err := json.Marshal(s)
	if err != nil {
		log.Printf("Error marshalling GamblerCreated payload %+v", err)
		return nil, err
	}
	envelope := events.Envelope{
		UUID:             getUID(),
		IsRootEvent:      false,
		SequenceNumber:   int64(0), // Set later by event-store
		SessionUID:       sessionUID,
		Timestamp:        mytime.Now(),
		AggregateName:    GamblerAggregateName, // from annotation!
		AggregateUID:     s.GetUID(),
		EventTypeName:    GamblerCreatedEventName,
		EventTypeVersion: 0,
		EventData:        string(blob),
	}

	return &envelope, nil
}

// IsGamblerCreated detects of envelope carries event of type GamblerCreated
func IsGamblerCreated(envelope *events.Envelope) bool {
	return envelope.EventTypeName == GamblerCreatedEventName
}

// GetIfIsGamblerCreated detects of envelope carries event of type GamblerCreated and returns the event if so
func GetIfIsGamblerCreated(envelope *events.Envelope) (*GamblerCreated, bool) {
	if IsGamblerCreated(envelope) == false {
		return nil, false
	}
	event, err := UnWrapGamblerCreated(envelope)
	if err != nil {
		return nil, false
	}
	return event, true
}

// UnWrapGamblerCreated extracts event GamblerCreated from its envelope
func UnWrapGamblerCreated(envelope *events.Envelope) (*GamblerCreated, error) {
	if IsGamblerCreated(envelope) == false {
		return nil, fmt.Errorf("Not a GamblerCreated")
	}
	var event GamblerCreated
	err := json.Unmarshal([]byte(envelope.EventData), &event)
	if err != nil {
		log.Printf("Error unmarshalling GamblerCreated payload %+v", err)
		return nil, err
	}
	event.Timestamp = envelope.Timestamp.In(mytime.DutchLocation())

	return &event, nil
}

// Wrap wraps event GamblerTeamCreated into an envelope
func (s *GamblerTeamCreated) Wrap(sessionUID string) (*events.Envelope, error) {
	blob, err := json.Marshal(s)
	if err != nil {
		log.Printf("Error marshalling GamblerTeamCreated payload %+v", err)
		return nil, err
	}
	envelope := events.Envelope{
		UUID:             getUID(),
		IsRootEvent:      false,
		SequenceNumber:   int64(0), // Set later by event-store
		SessionUID:       sessionUID,
		Timestamp:        mytime.Now(),
		AggregateName:    GamblerAggregateName, // from annotation!
		AggregateUID:     s.GetUID(),
		EventTypeName:    GamblerTeamCreatedEventName,
		EventTypeVersion: 0,
		EventData:        string(blob),
	}

	return &envelope, nil
}

// IsGamblerTeamCreated detects of envelope carries event of type GamblerTeamCreated
func IsGamblerTeamCreated(envelope *events.Envelope) bool {
	return envelope.EventTypeName == GamblerTeamCreatedEventName
}

// GetIfIsGamblerTeamCreated detects of envelope carries event of type GamblerTeamCreated and returns the event if so
func GetIfIsGamblerTeamCreated(envelope *events.Envelope) (*GamblerTeamCreated, bool) {
	if IsGamblerTeamCreated(envelope) == false {
		return nil, false
	}
	event, err := UnWrapGamblerTeamCreated(envelope)
	if err != nil {
		return nil, false
	}
	return event, true
}

// UnWrapGamblerTeamCreated extracts event GamblerTeamCreated from its envelope
func UnWrapGamblerTeamCreated(envelope *events.Envelope) (*GamblerTeamCreated, error) {
	if IsGamblerTeamCreated(envelope) == false {
		return nil, fmt.Errorf("Not a GamblerTeamCreated")
	}
	var event GamblerTeamCreated
	err := json.Unmarshal([]byte(envelope.EventData), &event)
	if err != nil {
		log.Printf("Error unmarshalling GamblerTeamCreated payload %+v", err)
		return nil, err
	}
	event.Timestamp = envelope.Timestamp.In(mytime.DutchLocation())

	return &event, nil
}

// Wrap wraps event NewsItemCreated into an envelope
func (s *NewsItemCreated) Wrap(sessionUID string) (*events.Envelope, error) {
	blob, err := json.Marshal(s)
	if err != nil {
		log.Printf("Error marshalling NewsItemCreated payload %+v", err)
		return nil, err
	}
	envelope := events.Envelope{
		UUID:             getUID(),
		IsRootEvent:      false,
		SequenceNumber:   int64(0), // Set later by event-store
		SessionUID:       sessionUID,
		Timestamp:        mytime.Now(),
		AggregateName:    NewsAggregateName, // from annotation!
		AggregateUID:     s.GetUID(),
		EventTypeName:    NewsItemCreatedEventName,
		EventTypeVersion: 0,
		EventData:        string(blob),
	}

	return &envelope, nil
}

// IsNewsItemCreated detects of envelope carries event of type NewsItemCreated
func IsNewsItemCreated(envelope *events.Envelope) bool {
	return envelope.EventTypeName == NewsItemCreatedEventName
}

// GetIfIsNewsItemCreated detects of envelope carries event of type NewsItemCreated and returns the event if so
func GetIfIsNewsItemCreated(envelope *events.Envelope) (*NewsItemCreated, bool) {
	if IsNewsItemCreated(envelope) == false {
		return nil, false
	}
	event, err := UnWrapNewsItemCreated(envelope)
	if err != nil {
		return nil, false
	}
	return event, true
}

// UnWrapNewsItemCreated extracts event NewsItemCreated from its envelope
func UnWrapNewsItemCreated(envelope *events.Envelope) (*NewsItemCreated, error) {
	if IsNewsItemCreated(envelope) == false {
		return nil, fmt.Errorf("Not a NewsItemCreated")
	}
	var event NewsItemCreated
	err := json.Unmarshal([]byte(envelope.EventData), &event)
	if err != nil {
		log.Printf("Error unmarshalling NewsItemCreated payload %+v", err)
		return nil, err
	}
	event.Timestamp = envelope.Timestamp.In(mytime.DutchLocation())

	return &event, nil
}
