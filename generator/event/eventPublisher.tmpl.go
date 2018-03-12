package event

const eventPublisherTemplate = `// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}Publisher

import "golang.org/x/net/context"

{{range .Structs -}}

	{{if IsTransientEvent . -}}

// PublishEvent{{.Name}} is used to publish event of type {{.Name}}
func PublishEvent{{.Name}}(c context.Context, rc request.Context, event *{{.PackageName}}.{{.Name}}) error {
    envelope, err := event.Wrap(rc)
    if err != nil {
        return errorh.NewInternalErrorf(0, "Error wrapping %s event %s: %s", envelope.EventTypeName, event.GetUID(), err)
    }

    err = publisher.PublishEnvelope(c, rc, envelope)
    if err != nil {
        return errorh.NewInternalErrorf(0, "Error publishing %s event %s: %s", envelope.EventTypeName, event.GetUID(), err)
    }

    return nil
}
	{{end -}}
{{end -}}
`
