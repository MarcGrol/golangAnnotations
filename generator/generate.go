package generator

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/MarcGrol/astTools/model"
)

func GenerateForStruct(myStruct model.Struct) error {
	err := generateEnvelope(myStruct, "envelope")
	if err != nil {
		return err
	}
	err = generateWrapperForStruct(myStruct, "wrapper")
	if err != nil {
		return err
	}
	return nil
}

func generateEnvelope(str model.Struct, templateName string) error {
	targetDir := str.PackageName
	dir, _ := path.Split(str.PackageName)
	if dir == "" {
		targetDir = "."
	} else {
		str.PackageName = path.Dir(str.PackageName)
	}
	target := fmt.Sprintf("%s/envelope.go", targetDir)

	err := generateFileFromTemplate(str, templateName, target)
	if err != nil {
		log.Fatalf("Error generating events (%s)", err)
		return err
	}
	return nil
}

func generateWrapperForStruct(str model.Struct, templateName string) error {
	if str.IsEvent() {
		targetDir := str.PackageName
		dir, _ := path.Split(str.PackageName)
		if dir == "" {
			targetDir = "."
		} else {
			str.PackageName = path.Dir(str.PackageName)
		}
		target := fmt.Sprintf("%s/%sWrapperAgain.go", targetDir, str.Name)

		err := generateFileFromTemplate(str, templateName, target)
		if err != nil {
			log.Fatalf("Error generating events (%s)", err)
			return err
		}
	}
	return nil
}

func generateFileFromTemplate(data interface{}, templateName string, targetFileName string) error {
	log.Printf("Using template %s to generate target %s\n", templateName, targetFileName)

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
	"envelope": envelopeTemplate,
	"wrapper":  wrapperTemplate,
}

var envelopeTemplate string = `
// Generated automatically: do not edit manually

package {{.PackageName}}

import (
    "time"
)

type Uider interface {
    GetUid() string  
}

type Envelope struct {
    Uuid           string 
    SequenceNumber uint64 
    Timestamp      time.Time 
    AggregateName  string 
    AggregateUid   string  
    EventTypeName  string 
    EventData      string
}
`

var wrapperTemplate string = `
// Generated automatically: do not edit manually

package {{.PackageName}}

import (
  "encoding/json"
  "log"
  "time"

  "code.google.com/p/go-uuid/uuid"
)

func (s *{{.Name}}) Wrap() *Envelope {
    var err error
    envelope := new(Envelope)
    envelope.Uuid = uuid.New()
    envelope.SequenceNumber = 0 // Set later by event-store
    envelope.Timestamp = time.Now()
    envelope.AggregateName = "{{.GetAggregateName}}"
    envelope.AggregateUid = s.GetUid()
    envelope.EventTypeName = "{{.Name}}"
    blob, err := json.Marshal(s)
    if err != nil {
        log.Printf("Error marshalling {{.Name}} payload %+v", err)
        return nil //, err
    }
    envelope.EventData = string(blob)

    return envelope
}

func Is{{.Name}}(envelope *Envelope) bool {
    return envelope.EventTypeName == "{{.Name}}"
}

func GetIfIs{{.Name}}(envelop *Envelope) (*{{.Name}}, bool) {
    if Is{{.Name}}(envelop) == false {
        return nil, false
    }
    event := UnWrap{{.Name}}(envelop)
    return event, true
}

func UnWrap{{.Name}}(envelop *Envelope) *{{.Name}} {
    if Is{{.Name}}(envelop) == false {
        return nil
    }
    var event {{.Name}}
    err := json.Unmarshal([]byte(envelop.EventData), &event)
    if err != nil {
        log.Printf("Error unmarshalling {{.Name}} payload %+v", err)
        return nil
    }

    return &event
}
`
