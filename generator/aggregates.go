// Generated automatically: do not edit manually

package generator

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
