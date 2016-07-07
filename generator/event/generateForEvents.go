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

			err = generationUtil.GenerateFileFromTemplate(data, "aggregates", aggregateTemplate, customTemplateFuncs, target)
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
			err = generationUtil.GenerateFileFromTemplate(data, "wrappers", wrappersTemplate, customTemplateFuncs, target)
			if err != nil {
				log.Fatalf("Error generating wrappers for structs (%s)", err)
				return err
			}
		}
		{
			target := fmt.Sprintf("%s/wrappers_test.go", targetDir)

			data := Structs{
				PackageName: packageName,
				Structs:     structs,
			}
			err = generationUtil.GenerateFileFromTemplate(data, "wrappers-test", wrappersTestTemplate, customTemplateFuncs, target)
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
			return "[]string{" + fmt.Sprintf("\"Eaample1%s\"", field.Name) + "," +
				fmt.Sprintf("\"Eaample2%s\"", field.Name) + "}"
		} else {
			return fmt.Sprintf("\"Eaample%s\"", field.Name)
		}
	} else if field.TypeName == "bool" {
		return "true"
	}
	return ""
}

var aggregateTemplate string = `
// Generated automatically: do not edit manually

package {{.PackageName}}

import "fmt"

const (
{{range $aggr, $events := .AggregateMap}}
    {{$aggr}}AggregateName = "{{$aggr}}"
{{end}}
)
var AggregateEvents map[string][]string = map[string][]string{
{{range $aggr, $events := .AggregateMap}}
	{{$aggr}}AggregateName: []string {
	{{range $aggregName, $eventName := $events}}
		{{$eventName}}EventName,
	{{end}}
	},
{{end}}
}

{{range $aggr, $events := .AggregateMap}}
type {{$aggr}}Aggregate interface {
	ApplyAll(envelopes []Envelope)
	{{range $aggregName, $eventName := $events}}
		Apply{{$eventName}}(event {{$eventName}})
	{{end}}
}

func Apply{{$aggr}}Event(envelop Envelope, aggregateRoot {{$aggr}}Aggregate) error {
	switch envelop.EventTypeName {
	{{range $aggregName, $eventName := $events}}
		case {{$eventName}}EventName:
		event, err := 	UnWrap{{$eventName}}(&envelop)
		if err != nil {
			return err
		}
		aggregateRoot.Apply{{$eventName}}(*event)
		break
	{{end}}
	default:
		return fmt.Errorf("Apply{{$aggr}}Event: Unexpected event %s", envelop.EventTypeName)
	}
	return nil
}

func Apply{{$aggr}}Events(envelopes []Envelope, aggregateRoot {{$aggr}}Aggregate) error {
	var err error
	for _, envelop := range envelopes {
		err = Apply{{$aggr}}Event(envelop, aggregateRoot)
		if err != nil {
			break
		}
	}
	return err
}

{{end}} 
`

var wrappersTemplate string = `
// Generated automatically: do not edit manually

package {{.PackageName}}

import (
  "encoding/json"
  "fmt"
  "log"
  "time"

  "github.com/satori/go.uuid"
)

type Envelope struct {
	Uuid           string
	SequenceNumber int64
	Timestamp      time.Time
	AggregateName  string
	AggregateUid   string
	EventTypeName  string
	EventData      string
}

const (
{{range .Structs}}
{{if IsEvent . }}
	{{.Name}}EventName = "{{.Name}}"
{{end}}
{{end}}
)

type getTimeFunc func() time.Time

var getTime getTimeFunc = func() time.Time {
	return time.Now()
}

type getUidFunc func() string

var getUid getUidFunc = func() string {
	return uuid.NewV1().String()
}

{{range .Structs}}
{{if IsEvent . }}

func (s *{{.Name}}) Wrap(uid string) (*Envelope,error) {
    blob, err := json.Marshal(s)
    if err != nil {
        log.Printf("Error marshalling {{.Name}} payload %+v", err)
        return nil, err
    }
	envelope := Envelope{
		Uuid: getUid(),
		SequenceNumber: uint64(0), // Set later by event-store
		Timestamp: getTime(),
		AggregateName: {{GetAggregateName . }}AggregateName, // from annotation!
		AggregateUid: uid,
		EventTypeName: {{.Name}}EventName,
		EventData: string(blob),
    }

    return &envelope, nil
}

func Is{{.Name}}(envelope *Envelope) bool {
    return envelope.EventTypeName == {{.Name}}EventName
}

func GetIfIs{{.Name}}(envelop *Envelope) (*{{.Name}}, bool) {
    if Is{{.Name}}(envelop) == false {
        return nil, false
    }
    event,err := UnWrap{{.Name}}(envelop)
    if err != nil {
    	return nil, false
    }
    return event, true
}

func UnWrap{{.Name}}(envelop *Envelope) (*{{.Name}},error) {
    if Is{{.Name}}(envelop) == false {
        return nil, fmt.Errorf("Not a {{.Name}}")
    }
    var event {{.Name}}
    err := json.Unmarshal([]byte(envelop.EventData), &event)
    if err != nil {
        log.Printf("Error unmarshalling {{.Name}} payload %+v", err)
        return nil, err
    }

    return &event, nil
}
{{end}}
{{end}}
`
var wrappersTestTemplate string = `
// Generated automatically: do not edit manually

package {{.PackageName}}

import (
	"testing"
	"time"
	"reflect"

	"github.com/stretchr/testify/assert"
)

func testGetTime() time.Time {
	t, _ := time.Parse(time.RFC3339Nano, "2003-02-11T11:50:51.123Z")
	return t
}

func testGetUid() string {
	return "1234321"
}

{{range .Structs}}
{{if IsEvent . }}

func Test{{.Name}}Wrapper(t *testing.T) {
	getUid = testGetUid
	getTime = testGetTime

	event := {{.Name}}{
	   {{range .Fields}}
	   {{if HasValueForField .}} {{.Name}}: {{ValueForField .}}, {{end}} {{end}}
	}
	wrapped, err := event.Wrap("UID_{{.Name}}")
	assert.NoError(t, err)
	assert.True(t, Is{{.Name}}(wrapped))
    assert.Equal(t, "{{GetAggregateName . }}", wrapped.AggregateName)
    assert.Equal(t, "{{.Name}}", wrapped.EventTypeName)
	assert.Equal(t, "UID_{{.Name}}", wrapped.AggregateUid)
    assert.Equal(t, "1234321", wrapped.Uuid)
    assert.Equal(t, "2003-02-11T11:50:51.123Z", wrapped.Timestamp.Format(time.RFC3339Nano))
	assert.Equal(t, uint64(0), wrapped.SequenceNumber)
	again, ok := GetIfIs{{.Name}}(wrapped)
	assert.True(t, ok)
	assert.NotNil(t,again)
	reflect.DeepEqual(event, *again)
}
{{end}}
{{end}}
`
