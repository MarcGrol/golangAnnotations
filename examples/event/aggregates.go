// Generated automatically: do not edit manually

package event

import "fmt"

const (
	GamblerAggregateName = "Gambler"

	NewsAggregateName = "News"

	TourAggregateName = "Tour"
)

var AggregateEvents map[string][]string = map[string][]string{

	GamblerAggregateName: []string{

		GamblerCreatedEventName,

		GamblerTeamCreatedEventName,
	},

	NewsAggregateName: []string{

		NewsItemCreatedEventName,
	},

	TourAggregateName: []string{

		CyclistCreatedEventName,

		EtappeCreatedEventName,

		EtappeResultsCreatedEventName,

		TourCreatedEventName,
	},
}

type GamblerAggregate interface {
	ApplyGamblerCreated(event GamblerCreated)

	ApplyGamblerTeamCreated(event GamblerTeamCreated)
}

func ApplyGamblerEvent(envelop Envelope, aggregateRoot GamblerAggregate) error {
	switch envelop.EventTypeName {

	case GamblerCreatedEventName:
		event, err := UnWrapGamblerCreated(&envelop)
		if err != nil {
			return err
		}
		aggregateRoot.ApplyGamblerCreated(*event)
		break

	case GamblerTeamCreatedEventName:
		event, err := UnWrapGamblerTeamCreated(&envelop)
		if err != nil {
			return err
		}
		aggregateRoot.ApplyGamblerTeamCreated(*event)
		break

	default:
		return fmt.Errorf("ApplyGamblerEvent: Unexpected event %s", envelop.EventTypeName)
	}
	return nil
}

func ApplyGamblerEvents(envelopes []Envelope, aggregateRoot GamblerAggregate) error {
	var err error
	for _, envelop := range envelopes {
		err = ApplyGamblerEvent(envelop, aggregateRoot)
		if err != nil {
			break
		}
	}
	return err
}

type NewsAggregate interface {
	ApplyNewsItemCreated(event NewsItemCreated)
}

func ApplyNewsEvent(envelop Envelope, aggregateRoot NewsAggregate) error {
	switch envelop.EventTypeName {

	case NewsItemCreatedEventName:
		event, err := UnWrapNewsItemCreated(&envelop)
		if err != nil {
			return err
		}
		aggregateRoot.ApplyNewsItemCreated(*event)
		break

	default:
		return fmt.Errorf("ApplyNewsEvent: Unexpected event %s", envelop.EventTypeName)
	}
	return nil
}

func ApplyNewsEvents(envelopes []Envelope, aggregateRoot NewsAggregate) error {
	var err error
	for _, envelop := range envelopes {
		err = ApplyNewsEvent(envelop, aggregateRoot)
		if err != nil {
			break
		}
	}
	return err
}

type TourAggregate interface {
	ApplyCyclistCreated(event CyclistCreated)

	ApplyEtappeCreated(event EtappeCreated)

	ApplyEtappeResultsCreated(event EtappeResultsCreated)

	ApplyTourCreated(event TourCreated)
}

func ApplyTourEvent(envelop Envelope, aggregateRoot TourAggregate) error {
	switch envelop.EventTypeName {

	case CyclistCreatedEventName:
		event, err := UnWrapCyclistCreated(&envelop)
		if err != nil {
			return err
		}
		aggregateRoot.ApplyCyclistCreated(*event)
		break

	case EtappeCreatedEventName:
		event, err := UnWrapEtappeCreated(&envelop)
		if err != nil {
			return err
		}
		aggregateRoot.ApplyEtappeCreated(*event)
		break

	case EtappeResultsCreatedEventName:
		event, err := UnWrapEtappeResultsCreated(&envelop)
		if err != nil {
			return err
		}
		aggregateRoot.ApplyEtappeResultsCreated(*event)
		break

	case TourCreatedEventName:
		event, err := UnWrapTourCreated(&envelop)
		if err != nil {
			return err
		}
		aggregateRoot.ApplyTourCreated(*event)
		break

	default:
		return fmt.Errorf("ApplyTourEvent: Unexpected event %s", envelop.EventTypeName)
	}
	return nil
}

func ApplyTourEvents(envelopes []Envelope, aggregateRoot TourAggregate) error {
	var err error
	for _, envelop := range envelopes {
		err = ApplyTourEvent(envelop, aggregateRoot)
		if err != nil {
			break
		}
	}
	return err
}
