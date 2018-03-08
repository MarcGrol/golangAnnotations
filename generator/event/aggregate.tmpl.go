package event

const aggregateTemplate = `// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}

import (
	"fmt"
	"golang.org/x/net/context"
)

const (
{{range $aggr, $events := .AggregateMap -}}
	// {{$aggr}}AggregateName provides constant for the name of {{$aggr}}
	{{$aggr}}AggregateName = "{{$aggr}}"
{{end -}}
)

// AggregateEvents describes all aggregates with their events
var AggregateEvents = map[string][]string{
{{range $aggr, $events := .AggregateMap -}}
	{{$aggr}}AggregateName: []string {
		{{range $aggregName, $event := $events -}}
			{{$event.Name}}EventName,
		{{end -}}
	},
{{end -}}
}

{{range $aggr, $events := .AggregateMap}}

// {{$aggr}}Aggregate provides an interface that forces all events related to an aggregate are handled
type {{$aggr}}Aggregate interface {
	idempotency.Checker
	{{range $aggregName, $event := $events -}}
    	{{if $event.IsPersistent -}}
			Apply{{$event.Name}}(c context.Context, event {{$event.Name}})
		{{end -}}
	{{end -}}
}

// Apply{{$aggr}}Event applies a single event to aggregate {{$aggr}}
func Apply{{$aggr}}Event(c context.Context, envelope envelope.Envelope, aggregateRoot {{$aggr}}Aggregate) error {
	if aggregateRoot.IsEventProcessed(envelope.UUID){
		 mylog.New().Warning(c, "Event %+v already processed", envelope)
		 return nil
	}

	switch envelope.EventTypeName {
		{{range $aggregName, $event := $events -}}{{if $event.IsPersistent -}}
			case {{$event.Name}}EventName:
			event, err :=    UnWrap{{$event.Name}}(&envelope)
			if err != nil {
					return err
			}
			aggregateRoot.Apply{{$event.Name}}(c, *event)
			break
		{{end -}}{{end -}}
		default:
		return fmt.Errorf("Apply{{$aggr}}Event: Unexpected event %s", envelope.EventTypeName)
	}

	aggregateRoot.MarkEventProcessed(envelope.UUID)
	return nil
}

// Apply{{$aggr}}Events applies multiple events to aggregate {{$aggr}}
func Apply{{$aggr}}Events(c context.Context, envelopes []envelope.Envelope, aggregateRoot {{$aggr}}Aggregate) error {
	var err error
	for _, envelope := range envelopes {
		err = Apply{{$aggr}}Event(c, envelope, aggregateRoot)
		if err != nil {
				break
		}
	}
	return err
}

// UnWrap{{$aggr}}Event extracts the event from its envelope
func UnWrap{{$aggr}}Event(envelope *envelope.Envelope) (envelope.Event, error) {
	switch envelope.EventTypeName {
		{{range $aggregName, $event := $events -}}
			case {{$event.Name}}EventName:
				event, err := UnWrap{{$event.Name}}(envelope)
				if err != nil {
					return nil, err
				}
				return event, nil
		{{end -}}
		default:
		return nil, fmt.Errorf("UnWrap{{$aggr}}Event: Unexpected event %s", envelope.EventTypeName)
	}
}

// UnWrap{{$aggr}}Events extracts the events from multiple envelopes
func UnWrap{{$aggr}}Events(envelopes []envelope.Envelope) ([]envelope.Event, error) {
	events := make([]envelope.Event, 0, len(envelopes))
	for _, envelope := range envelopes {
		event, err := UnWrap{{$aggr}}Event(&envelope)
		if err != nil {
				return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

{{end -}}
`
