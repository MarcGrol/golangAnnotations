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

func GenerateForStruct(myStruct model.Struct) error {

	path, err := determineTargetPath(myStruct.PackageName)
	log.Printf("target:%s %v\n", path, err)

	err = generateEnvelope(myStruct)
	if err != nil {
		return err
	}

	err = generateWrapperForStruct(myStruct)
	if err != nil {
		return err
	}
	return nil
}

func determineTargetPath(packageName string) (string, error) {
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

	baseDir := path.Base(workDir)
	if baseDir == packageName {
		return ".", nil
	} else {
		return packageName, nil
	}
}

func generateEnvelope(data model.Struct) error {
	targetDir, err := determineTargetPath(data.PackageName)
	if err != nil {
		return err
	}

	target := fmt.Sprintf("%s/envelope.go", targetDir)

	err = generateFileFromTemplate(data, "envelope", target)
	if err != nil {
		log.Fatalf("Error generating events (%s)", err)
		return err
	}
	return nil
}

func generateWrapperForStruct(data model.Struct) error {
	if data.IsEvent() {
		targetDir, err := determineTargetPath(data.PackageName)
		if err != nil {
			return err
		}
		target := fmt.Sprintf("%s/%sWrapper.go", targetDir, data.Name)

		err = generateFileFromTemplate(data, "wrapper", target)
		if err != nil {
			log.Fatalf("Error generating events (%s)", err)
			return err
		}
	}
	return nil
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
  "fmt"
  "log"
  "time"

  "code.google.com/p/go-uuid/uuid"
)

func (s *{{.Name}}) Wrap() (*Envelope,error) {
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
        return nil, err
    }
    envelope.EventData = string(blob)

    return envelope, nil
}

func Is{{.Name}}(envelope *Envelope) bool {
    return envelope.EventTypeName == "{{.Name}}"
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
`
