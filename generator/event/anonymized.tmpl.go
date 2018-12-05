package event

const anonymizedTemplate = `// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}

{{range .Structs -}}
	{{if IsSensitiveEvent . -}}

// Anonymizes event {{.Name}}: wipes all data marked as sensitive
func ({{EventIdentifier .}} *{{.Name}}) Anonymized() *{{.Name}} {
	{{$evt := EventIdentifier . -}}
	{{range .Fields -}}
		{{if IsSensitiveField . -}}
			{{if IsPrimitive . -}}
				{{if IsInt . -}}
					{{$evt}}.{{.Name}} = 0
				{{else if IsBool . -}}
					{{$evt}}.{{.Name}} = false
				{{else if IsString . -}}
					{{$evt}}.{{.Name}} = ""
				{{else -}}
					Force compile error: field {{.Name}} has unsupported primitive type
				{{end -}}
			{{else if IsPointer . -}}
				{{$evt}}.{{.Name}} = nil
			{{else if IsStringSlice . -}}
				{{$evt}}.{{.Name}} = []string{}
			{{else if IsSlice . -}}
				{{$evt}}.{{.Name}} = {{.TypeName}}{}
			{{else if IsDate . -}}
				{{$evt}}.{{.Name}} = mydate.MyDate{}
			{{else -}}
				{{$evt}}.{{.Name}} = {{$evt}}.{{.Name}}.Anonymized()
			{{end -}}
		{{else -}}
			{{if IsStringSlice . -}}
			{{else if IsSlice . -}}
				for idx, {{SliceFieldIdentifier .}} := range {{$evt}}.{{.Name}} {
					{{$evt}}.{{.Name}}[idx] = {{SliceFieldIdentifier .}}.Anonymized()
				}
			{{end -}}
		{{end -}}
	{{end -}}
	return {{$evt}}
}

	{{end -}}
{{end -}}
`
