// Generated automatically: do not edit manually

package generator

import "fmt"

const (
	PersonAggregateName = "Person"
)

var AggregateEvents map[string][]string = map[string][]string{

	PersonAggregateName: []string{

		MyStructEventName,
	},
}

type PersonAggregate interface {
	ApplyAll(envelopes []Envelope)

	ApplyMyStruct(event MyStruct)
}

func ApplyPersonEvent(envelop Envelope, aggregateRoot PersonAggregate) error {
	switch envelop.EventTypeName {

	case MyStructEventName:
		event, err := UnWrapMyStruct(&envelop)
		if err != nil {
			return err
		}
		aggregateRoot.ApplyMyStruct(*event)
		break

	default:
		return fmt.Errorf("ApplyPersonEvent: Unexpected event %s", envelop.EventTypeName)
	}
	return nil
}

func ApplyPersonEvents(envelopes []Envelope, aggregateRoot PersonAggregate) error {
	var err error
	for _, envelop := range envelopes {
		err = ApplyPersonEvent(envelop, aggregateRoot)
		if err != nil {
			break
		}
	}
	return err
}
