// Generated automatically: do not edit manually

package generator

import "fmt"

const (
	personAggregateName = "person"
)

var AggregateEvents map[string][]string = map[string][]string{

	personAggregateName: []string{

		MyStructEventName,
	},
}

type Aggregateperson interface {
	ApplyAll(envelopes []Envelope)

	ApplyMyStruct(event MyStruct)
}

func ApplypersonEvent(envelop Envelope, aggregateRoot Aggregateperson) error {
	switch envelop.EventTypeName {

	case MyStructEventName:
		event, err := UnWrapMyStruct(&envelop)
		if err != nil {
			return err
		}
		aggregateRoot.ApplyMyStruct(*event)
		break

	default:
		return fmt.Errorf("ApplypersonEvent: Unexpected event %s", envelop.EventTypeName)
	}
	return nil
}

func ApplypersonEvents(envelopes []Envelope, aggregateRoot Aggregateperson) error {
	var err error
	for _, envelop := range envelopes {
		err = ApplypersonEvent(envelop, aggregateRoot)
		if err != nil {
			break
		}
	}
	return err
}
