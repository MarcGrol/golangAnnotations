// Generated automatically: do not edit manually

package generator

import "github.com/MarcGrol/tourApp/events"

var AggregateEvents map[string][]string = map[string][]string{

	"person": []string{

		"MyStruct",
	},
}

type Aggregateperson interface {
	ApplyAll(envelopes []Envelope)
	ApplyTourCreated(event *events.TourCreated)

	ApplyMyStruct(event MyStruct)
}
