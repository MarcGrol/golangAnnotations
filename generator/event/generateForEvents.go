package event

import (
	"fmt"
	"html/template"
	"log"
	"strings"

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

func Generate(inputDir string, structs []model.Struct) error {
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
	}
	return nil
}

var customTemplateFuncs = template.FuncMap{
	"IsEvent":          IsEvent,
	"GetAggregateName": GetAggregateName,
	"ToFirstUpper":     ToFirstUpper,
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

func ToFirstUpper(in string) string {
	if len(in) == 0 {
		return in
	}
	return strings.ToUpper(fmt.Sprintf("%c", in[0])) + in[1:]
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
	SequenceNumber uint64
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

{{range .Structs}}
{{if IsEvent . }}

func (s *{{.Name}}) Wrap(uid string) (*Envelope,error) {
    envelope := new(Envelope)
    envelope.Uuid = uuid.NewV1().String()
    envelope.SequenceNumber = 0 // Set later by event-store
    envelope.Timestamp = time.Now()
    envelope.AggregateName = {{GetAggregateName . }}AggregateName // from annotation!
    envelope.AggregateUid = uid
    envelope.EventTypeName = {{.Name}}EventName
    blob, err := json.Marshal(s)
    if err != nil {
        log.Printf("Error marshalling {{.Name}} payload %+v", err)
        return nil, err
    }
    envelope.EventData = string(blob)

    return envelope, nil
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
