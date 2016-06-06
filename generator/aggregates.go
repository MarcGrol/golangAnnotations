
// Generated automatically: do not edit manually

package generator

import "fmt"

const (

    TestAggregateName = "Test"

)
var AggregateEvents map[string][]string = map[string][]string{

	TestAggregateName: []string {
	
		MyStructEventName,
	
	},

}


type TestAggregate interface {
	ApplyAll(envelopes []Envelope)
	
		ApplyMyStruct(event MyStruct)
	
}

func ApplyTestEvent(envelop Envelope, aggregateRoot TestAggregate) error {
	switch envelop.EventTypeName {
	
		case MyStructEventName:
		event, err := 	UnWrapMyStruct(&envelop)
		if err != nil {
			return err
		}
		aggregateRoot.ApplyMyStruct(*event)
		break
	
	default:
		return fmt.Errorf("ApplyTestEvent: Unexpected event %s", envelop.EventTypeName)
	}
	return nil
}

func ApplyTestEvents(envelopes []Envelope, aggregateRoot TestAggregate) error {
	var err error
	for _, envelop := range envelopes {
		err = ApplyTestEvent(envelop, aggregateRoot)
		if err != nil {
			break
		}
	}
	return err
}

 
