package generator

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/MarcGrol/astTools/model"
)

type AggregateMap struct {
	PackageName  string
	AggregateMap map[string]map[string]string
}

type Structs struct {
	PackageName string
	Structs     []model.Struct
}

func GenerateForStructs(inputDir string, structs []model.Struct) error {
	packageName, err := getPackageName(structs)
	if err != nil {
		return err
	}
	aggregates := make(map[string]map[string]string)
	for _, s := range structs {
		if s.IsEvent() {
			log.Printf("struct:%s -> %+v", s.GetAggregateName(), s)
			events, ok := aggregates[s.GetAggregateName()]
			if !ok {
				events = make(map[string]string)
			}
			events[s.Name] = s.Name
			aggregates[s.GetAggregateName()] = events
		}
	}

	targetDir, err := determineTargetPath(inputDir, packageName)
	if err != nil {
		return err
	}
	{
		target := fmt.Sprintf("%s/aggregates.go", targetDir)

		data := AggregateMap{
			PackageName:  packageName,
			AggregateMap: aggregates,
		}

		log.Printf("aggregates:%+v", data)
		err = generateFileFromTemplate(data, "aggregates", target)
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
		err = generateFileFromTemplate(data, "wrappers", target)
		if err != nil {
			log.Fatalf("Error generating wrappers for structs (%s)", err)
			return err
		}
	}
	return nil
}

func getPackageName(structs []model.Struct) (string, error) {
	if len(structs) == 0 {
		return "", fmt.Errorf("Need at least one struct to determine package-name")
	}
	packageName := structs[0].PackageName
	for _, s := range structs {
		if s.PackageName != packageName {
			return "", fmt.Errorf("List of structs has multiple package-names")
		}
	}
	return packageName, nil
}

func determineTargetPath(inputDir string, packageName string) (string, error) {
	log.Printf("inputDir:%s", inputDir)
	log.Printf("package:%s", packageName)

	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		return "", fmt.Errorf("GOPATH not set")
	}
	log.Printf("GOPATH:%s", goPath)

	workDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("Error getting working dir:%s", err)
	}
	log.Printf("work-dir:%s", workDir)

	if !strings.Contains(workDir, goPath) {
		return "", fmt.Errorf("Code %s lives outside GOPATH:%s", workDir, goPath)
	}

	baseDir := path.Base(inputDir)
	if baseDir == "." || baseDir == packageName {
		return inputDir, nil
	} else {
		return fmt.Sprintf("%s/%s", inputDir, packageName), nil
	}
}

func generateFileFromTemplate(data interface{}, templateName string, targetFileName string) error {
	log.Printf("Using template '%s' to generate target %s\n", templateName, targetFileName)

	err := os.MkdirAll(filepath.Dir(targetFileName), 0777)
	if err != nil {
		return err
	}
	w, err := os.Create(targetFileName)
	if err != nil {
		return err
	}

	t := template.New(templateName)
	t, err = t.Parse(templates[templateName])
	if err != nil {
		return err
	}

	defer w.Close()
	if err := t.Execute(w, data); err != nil {
		return err
	}
	return nil
}

var templates map[string]string = map[string]string{
	"aggregates": aggregateTemplate,
	"wrappers":   wrappersTemplate,
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
{{if .IsEvent}}
	{{.Name}}EventName = "{{.Name}}"
{{end}}
{{end}}
)

{{range .Structs}}
{{if .IsEvent}}

func (s *{{.Name}}) Wrap(uid string) (*Envelope,error) {
    envelope := new(Envelope)
    envelope.Uuid = uuid.NewV1().String()
    envelope.SequenceNumber = 0 // Set later by event-store
    envelope.Timestamp = time.Now()
    envelope.AggregateName = {{.GetAggregateName}}AggregateName // from annotation!
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
