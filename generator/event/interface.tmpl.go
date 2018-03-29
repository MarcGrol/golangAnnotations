package event

const interfaceTemplate = `// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}

import (
	"golang.org/x/net/context"

    "github.com/Duxxie/platform/backend/lib/request"
)

{{$packageName := .PackageName}}

type EventHandlerInterface interface {
{{range .Structs -}}
    {{if IsEvent . -}}
        on{{.Name}}( c context.Context, rc request.Context, event {{.Name}})
    {{end -}}
{{end -}}
}

/*
{{range .Structs -}}
    {{if IsEvent .}}
func (es *{{GetAggregateNameLowerCase .}}EventService)on{{.Name}}( c context.Context, rc request.Context, event {{$packageName}}.{{.Name}}) error {
	return nil
}
    {{end -}}
{{end -}}
*/

`
