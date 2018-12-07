package event

const anonymizedTemplate = `// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}

{{range .Structs -}}
	{{if IsSensitiveEventOrEventPart . -}}

// Anonymizes {{if IsEvent .}}event{{else}}event-part{{end}} {{.Name}}: wipes all data marked as sensitive
func ({{EventIdentifier .}} {{.Name}}) Anonymized() {{.Name}} {
	{{$evt := EventIdentifier . -}}
	{{range .Fields -}}
		{{if IsSensitiveField . -}}
			{{if IsPointer . -}}
				{{$evt}}.{{.Name}} = nil
			{{else if IsSlice . -}}
				{{$evt}}.{{.Name}} = {{.TypeName}}{}
			{{else if IsPrimitive . -}}
				{{if IsBool . -}}
					{{$evt}}.{{.Name}} = false
				{{else if IsInt . -}}
					{{$evt}}.{{.Name}} = 0
				{{else if IsString . -}}
					{{$evt}}.{{.Name}} = ""
				{{else -}}
					Force compile error: field {{.Name}} has unsupported primitive type
				{{end -}}
			{{else if IsDate . -}}
				{{$evt}}.{{.Name}} = mydate.MyDate{}
			{{else -}}
				{{if .Name -}}
					{{$evt}}.{{.Name}} = {{$evt}}.{{.Name}}.Anonymized()
				{{else -}}
					{{$evt}}.{{.TypeName}} = {{$evt}}.{{.TypeName}}.Anonymized()
				{{end -}}
			{{end -}}
		{{else if IsDeepSensitiveField . -}}
			{{if IsPointer . -}}
				{{if IsCustom . -}}
					if {{$evt}}.{{.Name}} != nil {
						{{FieldIdentifier .}} := {{$evt}}.{{.Name}}.Anonymized()
						{{$evt}}.{{.Name}} = &{{FieldIdentifier .}}
					}
				{{end -}}
			{{else if IsSlice . -}}
				{{if IsCustom . -}}
					for idx, {{SliceFieldIdentifier .}} := range {{$evt}}.{{.Name}} {
						{{$evt}}.{{.Name}}[idx] = {{SliceFieldIdentifier .}}.Anonymized()
					}
				{{end -}}
			{{end -}}
		{{end -}}
	{{end -}}
	return {{$evt}}
}

	{{end -}}
{{end -}}
`
