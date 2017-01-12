// Generated automatically by golangAnnotations: do not edit manually

package event

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/Duxxie/platform/backend/lib/events"
)

const (

	// GamblerAggregateName provides constant for the name of Gambler
	GamblerAggregateName = "Gambler"

	// NewsAggregateName provides constant for the name of News
	NewsAggregateName = "News"

	// TourAggregateName provides constant for the name of Tour
	TourAggregateName = "Tour"
)

// AggregateEvents describes all aggregates with their events
var AggregateEvents = map[string][]string{

	GamblerAggregateName: {

		GamblerCreatedEventName,

		GamblerTeamCreatedEventName,
	},

	NewsAggregateName: {

		NewsItemCreatedEventName,
	},

	TourAggregateName: {

		CyclistCreatedEventName,

		EtappeCreatedEventName,

		EtappeResultsCreatedEventName,

		TourCreatedEventName,
	},
}

// GamblerAggregate provides an interface that forces all events related to an aggregate are handled
type GamblerAggregate interface {
	ApplyGamblerCreated(c context.Context, event GamblerCreated)

	ApplyGamblerTeamCreated(c context.Context, event GamblerTeamCreated)
}

// ApplyGamblerEvent applies a single event to aggregate Gambler
func ApplyGamblerEvent(c context.Context, envelope events.Envelope, aggregateRoot GamblerAggregate) error {
	switch envelope.EventTypeName {

	case GamblerCreatedEventName:
		event, err := UnWrapGamblerCreated(&envelope)
		if err != nil {
			return err
		}
		aggregateRoot.ApplyGamblerCreated(c, *event)
		break

	case GamblerTeamCreatedEventName:
		event, err := UnWrapGamblerTeamCreated(&envelope)
		if err != nil {
			return err
		}
		aggregateRoot.ApplyGamblerTeamCreated(c, *event)
		break

	default:
		return fmt.Errorf("ApplyGamblerEvent: Unexpected event %s", envelope.EventTypeName)
	}
	return nil
}

// ApplyGamblerEvents applies multiple events to aggregate Gambler
func ApplyGamblerEvents(c context.Context, envelopes []events.Envelope, aggregateRoot GamblerAggregate) error {
	var err error
	for _, envelope := range envelopes {
		err = ApplyGamblerEvent(c, envelope, aggregateRoot)
		if err != nil {
			break
		}
	}
	return err
}

// UnWrapGamblerEvent extracts the event from its envelope
func UnWrapGamblerEvent(envelope *events.Envelope) (events.Event, error) {
	switch envelope.EventTypeName {

	case GamblerCreatedEventName:
		event, err := UnWrapGamblerCreated(envelope)
		if err != nil {
			return nil, err
		}
		return event, nil

	case GamblerTeamCreatedEventName:
		event, err := UnWrapGamblerTeamCreated(envelope)
		if err != nil {
			return nil, err
		}
		return event, nil

	default:
		return nil, fmt.Errorf("UnWrapGamblerEvent: Unexpected event %s", envelope.EventTypeName)
	}
}

// UnWrapGamblerEvents extracts the events from multiple envelopes
func UnWrapGamblerEvents(envelopes []events.Envelope) ([]events.Event, error) {
	events := make([]events.Event, 0, len(envelopes))
	for _, envelope := range envelopes {
		event, err := UnWrapGamblerEvent(&envelope)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

// NewsAggregate provides an interface that forces all events related to an aggregate are handled
type NewsAggregate interface {
	ApplyNewsItemCreated(c context.Context, event NewsItemCreated)
}

// ApplyNewsEvent applies a single event to aggregate News
func ApplyNewsEvent(c context.Context, envelope events.Envelope, aggregateRoot NewsAggregate) error {
	switch envelope.EventTypeName {

	case NewsItemCreatedEventName:
		event, err := UnWrapNewsItemCreated(&envelope)
		if err != nil {
			return err
		}
		aggregateRoot.ApplyNewsItemCreated(c, *event)
		break

	default:
		return fmt.Errorf("ApplyNewsEvent: Unexpected event %s", envelope.EventTypeName)
	}
	return nil
}

// ApplyNewsEvents applies multiple events to aggregate News
func ApplyNewsEvents(c context.Context, envelopes []events.Envelope, aggregateRoot NewsAggregate) error {
	var err error
	for _, envelope := range envelopes {
		err = ApplyNewsEvent(c, envelope, aggregateRoot)
		if err != nil {
			break
		}
	}
	return err
}

// UnWrapNewsEvent extracts the event from its envelope
func UnWrapNewsEvent(envelope *events.Envelope) (events.Event, error) {
	switch envelope.EventTypeName {

	case NewsItemCreatedEventName:
		event, err := UnWrapNewsItemCreated(envelope)
		if err != nil {
			return nil, err
		}
		return event, nil

	default:
		return nil, fmt.Errorf("UnWrapNewsEvent: Unexpected event %s", envelope.EventTypeName)
	}
}

// UnWrapNewsEvents extracts the events from multiple envelopes
func UnWrapNewsEvents(envelopes []events.Envelope) ([]events.Event, error) {
	events := make([]events.Event, 0, len(envelopes))
	for _, envelope := range envelopes {
		event, err := UnWrapNewsEvent(&envelope)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

// TourAggregate provides an interface that forces all events related to an aggregate are handled
type TourAggregate interface {
	ApplyCyclistCreated(c context.Context, event CyclistCreated)

	ApplyEtappeCreated(c context.Context, event EtappeCreated)

	ApplyEtappeResultsCreated(c context.Context, event EtappeResultsCreated)

	ApplyTourCreated(c context.Context, event TourCreated)
}

// ApplyTourEvent applies a single event to aggregate Tour
func ApplyTourEvent(c context.Context, envelope events.Envelope, aggregateRoot TourAggregate) error {
	switch envelope.EventTypeName {

	case CyclistCreatedEventName:
		event, err := UnWrapCyclistCreated(&envelope)
		if err != nil {
			return err
		}
		aggregateRoot.ApplyCyclistCreated(c, *event)
		break

	case EtappeCreatedEventName:
		event, err := UnWrapEtappeCreated(&envelope)
		if err != nil {
			return err
		}
		aggregateRoot.ApplyEtappeCreated(c, *event)
		break

	case EtappeResultsCreatedEventName:
		event, err := UnWrapEtappeResultsCreated(&envelope)
		if err != nil {
			return err
		}
		aggregateRoot.ApplyEtappeResultsCreated(c, *event)
		break

	case TourCreatedEventName:
		event, err := UnWrapTourCreated(&envelope)
		if err != nil {
			return err
		}
		aggregateRoot.ApplyTourCreated(c, *event)
		break

	default:
		return fmt.Errorf("ApplyTourEvent: Unexpected event %s", envelope.EventTypeName)
	}
	return nil
}

// ApplyTourEvents applies multiple events to aggregate Tour
func ApplyTourEvents(c context.Context, envelopes []events.Envelope, aggregateRoot TourAggregate) error {
	var err error
	for _, envelope := range envelopes {
		err = ApplyTourEvent(c, envelope, aggregateRoot)
		if err != nil {
			break
		}
	}
	return err
}

// UnWrapTourEvent extracts the event from its envelope
func UnWrapTourEvent(envelope *events.Envelope) (events.Event, error) {
	switch envelope.EventTypeName {

	case CyclistCreatedEventName:
		event, err := UnWrapCyclistCreated(envelope)
		if err != nil {
			return nil, err
		}
		return event, nil

	case EtappeCreatedEventName:
		event, err := UnWrapEtappeCreated(envelope)
		if err != nil {
			return nil, err
		}
		return event, nil

	case EtappeResultsCreatedEventName:
		event, err := UnWrapEtappeResultsCreated(envelope)
		if err != nil {
			return nil, err
		}
		return event, nil

	case TourCreatedEventName:
		event, err := UnWrapTourCreated(envelope)
		if err != nil {
			return nil, err
		}
		return event, nil

	default:
		return nil, fmt.Errorf("UnWrapTourEvent: Unexpected event %s", envelope.EventTypeName)
	}
}

// UnWrapTourEvents extracts the events from multiple envelopes
func UnWrapTourEvents(envelopes []events.Envelope) ([]events.Event, error) {
	events := make([]events.Event, 0, len(envelopes))
	for _, envelope := range envelopes {
		event, err := UnWrapTourEvent(&envelope)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}
