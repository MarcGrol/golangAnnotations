// Generated automatically: do not edit manually

package example

import "fmt"

const (
	AggregateName = ""
)

var AggregateEvents map[string][]string = map[string][]string{

	AggregateName: []string{

		TourServiceEventName,
	},
}

type Aggregate interface {
	ApplyAll(envelopes []Envelope)

	ApplyTourService(event TourService)
}

func ApplyEvent(envelop Envelope, aggregateRoot Aggregate) error {
	switch envelop.EventTypeName {

	case TourServiceEventName:
		event, err := UnWrapTourService(&envelop)
		if err != nil {
			return err
		}
		aggregateRoot.ApplyTourService(*event)
		break

	default:
		return fmt.Errorf("ApplyEvent: Unexpected event %s", envelop.EventTypeName)
	}
	return nil
}

func ApplyEvents(envelopes []Envelope, aggregateRoot Aggregate) error {
	var err error
	for _, envelop := range envelopes {
		err = ApplyEvent(envelop, aggregateRoot)
		if err != nil {
			break
		}
	}
	return err
}
