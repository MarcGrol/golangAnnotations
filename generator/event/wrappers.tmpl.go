package event

const wrappersTemplate = `// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}

import (
    "encoding/json"
    "fmt"
    "log"
)

const (
{{range .Structs -}}
    {{if IsEvent . -}}
        // {{.Name}}EventName provides a constant symbol for {{.Name}}
        {{.Name}}EventName = "{{.Name}}"
    {{end -}}
{{end -}}
)

{{range .Structs -}}
    {{if IsEvent . -}}

// Wrap wraps event {{.Name}} into an envelope
func (s *{{.Name}}) Wrap(credentials rest.Credentials) (*envelope.Envelope,error) {
	blob, err := json.Marshal(s)
    if err != nil {
        log.Printf("Error marshalling {{.Name}} payload %+v", err)
        return nil, err
    }
    envelope := envelope.Envelope{
        IsRootEvent:{{if IsRootEvent .}}true{{else}}false{{end}},
        SequenceNumber: int64(0), // Set later by event-store
        SessionUID: credentials.SessionUID,
        Timestamp: mytime.Now(),
        AggregateName: {{GetAggregateName . }}AggregateName, // from annotation!
        AggregateUID:  s.GetUID(),
        EventTypeName: {{.Name}}EventName,
        EventTypeVersion: 0,
        EventData: string(blob),
    }

	requestUID := credentials.RequestUID
	if requestUID == "" {	
		requestUID, _ = myuuid.NewV1({{GetAggregateName . }}AggregateName)
	}
	envelope.UUID = envelope.CreateRequestUID(requestUID)

    return &envelope, nil
}

// Is{{.Name}} detects of envelope carries event of type {{.Name}}
func Is{{.Name}}(envelope *envelope.Envelope) bool {
    return envelope.EventTypeName == {{.Name}}EventName
}

// GetIfIs{{.Name}} detects of envelope carries event of type {{.Name}} and returns the event if so
func GetIfIs{{.Name}}(envelope *envelope.Envelope) (*{{.Name}}, bool) {
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
func UnWrap{{.Name}}(envelope *envelope.Envelope) (*{{.Name}},error) {
    if Is{{.Name}}(envelope) == false {
        return nil, fmt.Errorf("Not a {{.Name}}")
    }
    var event {{.Name}}
    err := json.Unmarshal([]byte(envelope.EventData), &event)
    if err != nil {
        log.Printf("Error unmarshalling {{.Name}} payload %+v", err)
        return nil, err
    }

    event.Metadata = Metadata{
        UUID:          envelope.UUID,
        Timestamp:     envelope.Timestamp.In(mytime.DutchLocation),
        EventTypeName: envelope.EventTypeName,
    }

    return &event, nil
}

    {{end -}}
{{end -}}
`
