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
        On{{.Name}}( c context.Context, rc request.Context, event {{.Name}}) error
    {{end -}}
{{end -}}
}

/*
// Copy these empty implementations to your package to be able to easily detect missing methods

{{range .Structs -}}
    {{if IsEvent .}}
func (es *{{GetAggregateNameLowerCase .}}EventService)On{{.Name}}( _ context.Context, rc request.Context, _ {{$packageName}}.{{.Name}}) error {
	return nil
}
    {{end -}}
{{end -}}
*/

`
