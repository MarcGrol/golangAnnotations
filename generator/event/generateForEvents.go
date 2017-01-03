package event

import (
	"fmt"
	"log"
	"text/template"

	"github.com/MarcGrol/golangAnnotations/annotation"
	"github.com/MarcGrol/golangAnnotations/generator/event/eventAnnotation"
	"github.com/MarcGrol/golangAnnotations/generator/generationUtil"
	"github.com/MarcGrol/golangAnnotations/model"
)

type AggregateMap struct {
	PackageName  string
	AggregateMap map[string]map[string]string
}

type Structs struct {
	PackageName string
	Structs     []model.Struct
}

func Generate(inputDir string, parsedSource model.ParsedSources) error {
	return generate(inputDir, parsedSource.Structs)
}

func generate(inputDir string, structs []model.Struct) error {
	eventAnnotation.Register()

	packageName, err := generationUtil.GetPackageName(structs)
	if err != nil {
		return err
	}
	aggregates := make(map[string]map[string]string)
	eventCount := 0
	for _, s := range structs {
		if IsEvent(s) {
			events, ok := aggregates[GetAggregateName(s)]
			if !ok {
				events = make(map[string]string)
			}
			events[s.Name] = s.Name
			aggregates[GetAggregateName(s)] = events
			eventCount++
		}
	}

	if eventCount > 0 {
		targetDir, err := generationUtil.DetermineTargetPath(inputDir, packageName)
		if err != nil {
			return err
		}
		{
			target := fmt.Sprintf("%s/aggregates.go", targetDir)

			data := AggregateMap{
				PackageName:  packageName,
				AggregateMap: aggregates,
			}

			err = generationUtil.GenerateFileFromTemplate(data, packageName, "aggregates", aggregateTemplate, customTemplateFuncs, target)
			if err != nil {
				log.Fatalf("Error generating aggregates (%s)", err)
				return err
			}
		}
		{
			target := fmt.Sprintf("%s/wrappers.go", targetDir)

			data := Structs{
				PackageName: packageName,
				Structs:     structs,
			}
			err = generationUtil.GenerateFileFromTemplate(data, packageName, "wrappers", wrappersTemplate, customTemplateFuncs, target)
			if err != nil {
				log.Fatalf("Error generating wrappers for structs (%s)", err)
				return err
			}
		}
		{
			target := fmt.Sprintf("%s/../store/%sEventStore.go", targetDir, packageName)

			data := Structs{
				PackageName: packageName,
				Structs:     structs,
			}
			err = generationUtil.GenerateFileFromTemplate(data, packageName, "store-events", storeEventsTemplate, customTemplateFuncs, target)
			if err != nil {
				log.Fatalf("Error generating store-events for structs (%s)", err)
				return err
			}
		}
		{
			target := fmt.Sprintf("%s/wrappers_test.go", targetDir)

			data := Structs{
				PackageName: packageName,
				Structs:     structs,
			}
			err = generationUtil.GenerateFileFromTemplate(data, packageName, "wrappers-test", wrappersTestTemplate, customTemplateFuncs, target)
			if err != nil {
				log.Fatalf("Error generating wrappers-test for structs (%s)", err)
				return err
			}
		}

	}
	return nil
}

var customTemplateFuncs = template.FuncMap{
	"IsEvent":          IsEvent,
	"IsRootEvent":      IsRootEvent,
	"GetAggregateName": GetAggregateName,
	"HasValueForField": HasValueForField,
	"ValueForField":    ValueForField,
}

func IsEvent(s model.Struct) bool {
	annotation, ok := annotation.ResolveAnnotations(s.DocLines)
	if !ok || annotation.Name != "Event" {
		return false
	}
	return ok
}

func GetAggregateName(s model.Struct) string {
	val, ok := annotation.ResolveAnnotations(s.DocLines)
	if ok {
		return val.Attributes["aggregate"]
	}
	return ""
}

func IsRootEvent(s model.Struct) bool {
	annotation, ok := annotation.ResolveAnnotations(s.DocLines)
	if ok {
		isRootEvent := annotation.Attributes["isrootevent"]
		if isRootEvent == "true" {
			return true
		}
	}
	return false
}

func HasValueForField(field model.Field) bool {
	if field.TypeName == "int" || field.TypeName == "string" || field.TypeName == "bool" {
		return true
	}
	return false
}

func ValueForField(field model.Field) string {
	if field.TypeName == "int" {
		if field.IsSlice {
			return "[]int{1,2}"
		} else {
			return "42"
		}
	} else if field.TypeName == "string" {
		if field.IsSlice {
			return "[]string{" + fmt.Sprintf("\"Example1%s\"", field.Name) + "," +
				fmt.Sprintf("\"Example1%s\"", field.Name) + "}"
		} else {
			return fmt.Sprintf("\"Example3%s\"", field.Name)
		}
	} else if field.TypeName == "bool" {
		return "true"
	}
	return ""
}

var aggregateTemplate string = `
// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}

import (
    "golang.org/x/net/context"
)

const (
{{range $aggr, $events := .AggregateMap}}
    // {{$aggr}}AggregateName provides constant for the name of {{$aggr}}
    {{$aggr}}AggregateName = "{{$aggr}}"
{{end}}
)

// AggregateEvents describes all aggregates with their events
var AggregateEvents = map[string][]string{
{{range $aggr, $events := .AggregateMap}}
	{{$aggr}}AggregateName: []string {
	{{range $aggregName, $eventName := $events}}
		{{$eventName}}EventName,
	{{end}}
	},
{{end}}
}

{{range $aggr, $events := .AggregateMap}}
// {{$aggr}}Aggregate provides an interface that forces all events related to an aggregate are handled
type {{$aggr}}Aggregate interface {
	{{range $aggregName, $eventName := $events}}
		Apply{{$eventName}}(c context.Context, event {{$eventName}})
	{{end}}
}

// Apply{{$aggr}}Event applies a single event to aggregate {{$aggr}}
func Apply{{$aggr}}Event(c context.Context, envelope events.Envelope, aggregateRoot {{$aggr}}Aggregate) error {
	switch envelope.EventTypeName {
	{{range $aggregName, $eventName := $events}}
	case {{$eventName}}EventName:
		event, err := 	UnWrap{{$eventName}}(&envelope)
		if err != nil {
			return err
		}
		aggregateRoot.Apply{{$eventName}}(c, *event)
		break
	{{end}}
	default:
		return fmt.Errorf("Apply{{$aggr}}Event: Unexpected event %s", envelope.EventTypeName)
	}
	return nil
}

// Apply{{$aggr}}Events applies multiple events to aggregate {{$aggr}}
func Apply{{$aggr}}Events(c context.Context, envelopes []events.Envelope, aggregateRoot {{$aggr}}Aggregate) error {
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
func UnWrap{{$aggr}}Event(envelope *events.Envelope) (events.Event, error) {
	switch envelope.EventTypeName {
	{{range $aggregName, $eventName := $events}}
	case {{$eventName}}EventName:
		event, err := UnWrap{{$eventName}}(envelope)
		if err != nil {
			return nil, err
		}
		return event, nil
	{{end}}
	default:
		return nil, fmt.Errorf("UnWrap{{$aggr}}Event: Unexpected event %s", envelope.EventTypeName)
	}
}

// UnWrap{{$aggr}}Events extracts the events from multiple envelopes
func UnWrap{{$aggr}}Events(envelopes []events.Envelope) ([]events.Event, error) {
	events := make([]events.Event, 0, len(envelopes))
	for _, envelope := range envelopes {
		event, err := UnWrap{{$aggr}}Event(&envelope)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

{{end}}
`

var wrappersTemplate string = `
// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}

const (
{{range .Structs}}
{{if IsEvent . }}
    // {{.Name}}EventName provides a constant symbol for {{.Name}}
	{{.Name}}EventName = "{{.Name}}"
{{end}}
{{end}}
)

var getUID = func() string {
	return uuid.NewV1().String()
}

{{range .Structs}}
{{if IsEvent . }}

// Wrap wraps event {{.Name}} into an envelope
func (s *{{.Name}}) Wrap(sessionUID string) (*events.Envelope,error) {
    blob, err := json.Marshal(s)
    if err != nil {
        log.Printf("Error marshalling {{.Name}} payload %+v", err)
        return nil, err
    }
	envelope := events.Envelope{
		UUID: getUID(),
		IsRootEvent:{{if IsRootEvent .}}true{{else}}false{{end}},
		SequenceNumber: int64(0), // Set later by event-store
		SessionUID: sessionUID,
		Timestamp: mytime.Now(),
		AggregateName: {{GetAggregateName . }}AggregateName, // from annotation!
		AggregateUID:  s.GetUID(),
		EventTypeName: {{.Name}}EventName,
		EventTypeVersion: 0,
		EventData: string(blob),
    }

    return &envelope, nil
}

// Is{{.Name}} detects of envelope carries event of type {{.Name}}
func Is{{.Name}}(envelope *events.Envelope) bool {
    return envelope.EventTypeName == {{.Name}}EventName
}

// GetIfIs{{.Name}} detects of envelope carries event of type {{.Name}} and returns the event if so
func GetIfIs{{.Name}}(envelope *events.Envelope) (*{{.Name}}, bool) {
    if Is{{.Name}}(envelope) == false {
        return nil, false
    }
    event,err := UnWrap{{.Name}}(envelope)
    if err != nil {
    	return nil, false
    }
    return event, true
}

// UnWrap{{.Name}} extracts event {{.Name}} from its envelope
func UnWrap{{.Name}}(envelope *events.Envelope) (*{{.Name}},error) {
    if Is{{.Name}}(envelope) == false {
        return nil, fmt.Errorf("Not a {{.Name}}")
    }
    var event {{.Name}}
    err := json.Unmarshal([]byte(envelope.EventData), &event)
    if err != nil {
        log.Printf("Error unmarshalling {{.Name}} payload %+v", err)
        return nil, err
    }
    event.Timestamp = envelope.Timestamp.In(mytime.DutchLocation())

    return &event, nil
}
{{end}}
{{end}}
`
var wrappersTestTemplate string = `
// +build !appengine

// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}

func testGetUID() string {
	return "1234321"
}

{{range .Structs}}
{{if IsEvent . }}

func Test{{.Name}}Wrapper(t *testing.T) {
	mytime.SetMockNow()
	defer mytime.SetDefaultNow()
	getUID = testGetUID

	event := {{.Name}}{
	   {{range .Fields}}
	   {{if HasValueForField .}} {{.Name}}: {{ValueForField .}}, {{end}} {{end}}
	}
	wrapped, err := event.Wrap("test_session")
	assert.NoError(t, err)
	assert.True(t, Is{{.Name}}(wrapped))
    assert.Equal(t, "{{GetAggregateName . }}", wrapped.AggregateName)
    assert.Equal(t, "{{.Name}}", wrapped.EventTypeName)
	//	assert.Equal(t, "UID_{{.Name}}", wrapped.AggregateUID)
	assert.Equal(t, "test_session", wrapped.SessionUID)
	assert.Equal(t, "1234321", wrapped.UUID)
    assert.Equal(t, "2016-02-27T01:00:00+01:00", wrapped.Timestamp.Format(time.RFC3339))
	assert.Equal(t, int64(0), wrapped.SequenceNumber)
	again, ok := GetIfIs{{.Name}}(wrapped)
	assert.True(t, ok)
	assert.NotNil(t,again)
	reflect.DeepEqual(event, *again)
}
{{end}}
{{end}}
`

var storeEventsTemplate string = `
// Generated automatically by golangAnnotations: do not edit manually

package store

import (
	"golang.org/x/net/context"
)

{{range .Structs}}
{{if IsEvent . }}

// StoreEvent{{.Name}} is used to store event of type {{.Name}}
func StoreEvent{{.Name}}(c context.Context, event *{{.PackageName}}.{{.Name}}, sessionUID string) error {
	envlp, err := event.Wrap(sessionUID)
	if err != nil {
		return errorh.NewInternalErrorf(0, "Error wrapping %s event %s: %s", envlp.EventTypeName, event.GetUID(), err)
	}

	err = New().Put(c, envlp)
	if err != nil {
		return errorh.NewInternalErrorf(0, "Error storing %s event %s: %s", envlp.EventTypeName, event.GetUID(), err)
	}
	event.Timestamp = envlp.Timestamp
	return nil
}

{{end}}
{{end}}
`
