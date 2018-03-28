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
func (s *{{.Name}}) Wrap(rc request.Context) (*envelope.Envelope,error) {
	blob, err := json.Marshal(s)
    if err != nil {
        log.Printf("Error marshalling {{.Name}} payload %+v", err)
        return nil, err
    }
    envlp := envelope.Envelope{
        IsRootEvent:{{if IsRootEvent .}}true{{else}}false{{end}},
        SequenceNumber: int64(0), // Set later by event-store
        SessionUID: rc.GetSessionUID(),
        Timestamp: mytime.Now(),
        AggregateName: {{GetAggregateName . }}AggregateName, // from annotation!
        AggregateUID:  s.GetUID(),
        EventTypeName: {{.Name}}EventName,
        EventTypeVersion: 0,
        EventData: string(blob),
    }

	requestUID := rc.GetRequestUID()
	if requestUID == "" {	
		requestUID, _ = myuuid.NewV1({{GetAggregateName . }}AggregateName)
	}
	envlp.UUID = envlp.CreateRequestUID(requestUID)

    return &envlp, nil
}

// Is{{.Name}} detects of envelope carries event of type {{.Name}}
func Is{{.Name}}(envlp *envelope.Envelope) bool {
    return envlp.EventTypeName == {{.Name}}EventName
}

// GetIfIs{{.Name}} detects of envelope carries event of type {{.Name}} and returns the event if so
func GetIfIs{{.Name}}(envlp *envelope.Envelope) (*{{.Name}}, bool) {
    if Is{{.Name}}(envlp) == false {
        return nil, false
    }
    evt,err := UnWrap{{.Name}}(envlp)
    if err != nil {
        return nil, false
    }
    return evt, true
}

// UnWrap{{.Name}} extracts event {{.Name}} from its envelope
func UnWrap{{.Name}}(envlp *envelope.Envelope) (*{{.Name}},error) {
    if Is{{.Name}}(envlp) == false {
        return nil, fmt.Errorf("Not a {{.Name}}")
    }
    var evt {{.Name}}
    err := json.Unmarshal([]byte(envlp.EventData), &evt)
    if err != nil {
        log.Printf("Error unmarshalling {{.Name}} payload %+v", err)
        return nil, err
    }

    evt.Metadata = Metadata{
        UUID:          envlp.UUID,
		AdminUserUID:  envlp.AdminUserUID,
        Timestamp:     envlp.Timestamp.In(mytime.DutchLocation),
        EventTypeName: envlp.EventTypeName,
    }

    return &evt, nil
}

    {{end -}}
{{end -}}
`
