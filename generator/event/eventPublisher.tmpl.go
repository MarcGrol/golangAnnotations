package event

const eventPublisherTemplate = `// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}Publisher

import "golang.org/x/net/context"

{{range .Structs -}}

	{{if IsTransientEvent . -}}

// PublishEvent{{.Name}} is used to publish event of type {{.Name}}
func Publish{{.Name}}(c context.Context, rc request.Context, evt *{{.PackageName}}.{{.Name}}) error {
    envlp, err := evt.Wrap(rc)
    if err != nil {
        return errorh.NewInternalErrorf(0, "Error wrapping %s event %s: %s", envlp.EventTypeName, evt.GetUID(), err)
    }

    err = publisher.PublishEnvelope(c, rc, envlp)
    if err != nil {
        return errorh.NewInternalErrorf(0, "Error publishing %s event %s: %s", envlp.EventTypeName, evt.GetUID(), err)
    }

    return nil
}
	{{end -}}
{{end -}}
`
