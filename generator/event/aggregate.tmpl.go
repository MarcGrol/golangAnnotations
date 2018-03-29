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
		{{range $aggregName, $event := $events.Events -}}
			{{$event.Name}}EventName,
		{{end -}}
	},
{{end -}}
}

{{range $aggr, $events := .AggregateMap}}

// {{$aggr}}Aggregate provides an interface that forces all events related to an aggregate are handled
type {{$aggr}}Aggregate interface {
	idempotency.Checker
	{{range $aggregName, $event := $events.Events -}}
    	{{if $event.IsPersistent -}}
			Apply{{$event.Name}}(c context.Context, evt {{$event.Name}})
		{{end -}}
	{{end -}}
}

{{if $events.IsAnyPersistent -}}
// Apply{{$aggr}}Event applies a single event to aggregate {{$aggr}}
func Apply{{$aggr}}Event(c context.Context, envlp envelope.Envelope, aggregateRoot {{$aggr}}Aggregate) error {
	if aggregateRoot.IsEventProcessed(envlp.UUID){
		 mylog.New().Warning(c, "Event %+v already processed", envlp)
		 return nil
	}

	switch envlp.EventTypeName {
		{{range $aggregName, $event := $events.Events -}}{{if $event.IsPersistent -}}
			case {{$event.Name}}EventName:
			evt, err :=    UnWrap{{$event.Name}}(&envlp)
			if err != nil {
					return err
			}
			aggregateRoot.Apply{{$event.Name}}(c, *evt)
			break
		{{end -}}{{end -}}
		default:
		return fmt.Errorf("Apply{{$aggr}}Event: Unexpected event %s", envlp.EventTypeName)
	}

	aggregateRoot.MarkEventProcessed(envlp.UUID)
	return nil
}

// Apply{{$aggr}}Events applies multiple events to aggregate {{$aggr}}
func Apply{{$aggr}}Events(c context.Context, envelopes []envelope.Envelope, aggregateRoot {{$aggr}}Aggregate) error {
	var err error
	for _, envlp := range envelopes {
		err = Apply{{$aggr}}Event(c, envlp, aggregateRoot)
		if err != nil {
				break
		}
	}
	return err
}
{{end -}}

// UnWrap{{$aggr}}Event extracts the event from its envelope
func UnWrap{{$aggr}}Event(envlp *envelope.Envelope) (envelope.Event, error) {
	switch envlp.EventTypeName {
		{{range $aggregName, $event := $events.Events -}}
			case {{$event.Name}}EventName:
				evt, err := UnWrap{{$event.Name}}(envlp)
				if err != nil {
					return nil, err
				}
				return evt, nil
		{{end -}}
		default:
		return nil, fmt.Errorf("UnWrap{{$aggr}}Event: Unexpected event %s", envlp.EventTypeName)
	}
}

// UnWrap{{$aggr}}Events extracts the events from multiple envelopes
func UnWrap{{$aggr}}Events(envelopes []envelope.Envelope) ([]envelope.Event, error) {
	events := make([]envelope.Event, 0, len(envelopes))
	for _, envlp := range envelopes {
		evt, err := UnWrap{{$aggr}}Event(&envlp)
		if err != nil {
				return nil, err
		}
		events = append(events, evt)
	}
	return events, nil
}

{{end -}}
`
