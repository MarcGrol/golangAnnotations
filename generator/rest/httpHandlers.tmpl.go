package rest

const httpHandlersTemplate = `// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

{{ $service := . }}

var (
	preLogicHook  = func(c context.Context, w http.ResponseWriter, r *http.Request) {}
	postLogicHook = func(c context.Context, w http.ResponseWriter, r *http.Request, rc request.Context) {}
)

// HTTPHandler registers endpoint in new router
func (ts *{{.Name}}) HTTPHandler() http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	return ts.HTTPHandlerWithRouter(router)
}

// HTTPHandlerWithRouter registers endpoint in existing router
func (ts *{{.Name}}) HTTPHandlerWithRouter(router *mux.Router) *mux.Router {
	subRouter := router.PathPrefix("{{GetRestServicePath . }}").Subrouter()

	{{range .Operations -}}
		{{if IsRestOperation . -}}
			subRouter.HandleFunc("{{GetRestOperationPath . }}", {{.Name}}(ts)).Methods("{{GetRestOperationMethod . }}")
		{{end -}}
	{{end -}}

	return router
}

{{ $extractRequestContextMethod := GetExtractRequestContextMethod . }}
{{ $noValidation := IsRestServiceNoValidation . }}

{{range $idxOper, $oper := .Operations}}
	{{if IsRestOperation $oper -}}
		{{if IsRestOperationGenerated . -}}

// {{$oper.Name}} does the http handling for business logic method service.{{$oper.Name}}
func {{$oper.Name}}(service *{{$service.Name}}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		{{if NeedsContext $oper -}}
			{{GetContextName $oper}} := ctx.New.CreateContext(r)
			preLogicHook(c, w, r)
		{{else -}}
			preLogicHook(nil, w, r)
		{{end -}}

		rc := {{ $extractRequestContextMethod }}(c, r)

		{{if (not $noValidation) and (HasRequestContext $oper) -}}

			err = validateRequestContext(c, rc, {{GetRestOperationRolesString $oper}})
			if err != nil {
				errorh.HandleHTTPError(c, rc, err, w, r)
				return
			}

		{{end -}}

		{{if HasUpload . -}}

			// Note: blobstore.ParseUpload must be called before parsing request POST-params
			{{GetInputArgName . }}, err := service.{{$oper.Name}}GetUpload({{GetContextName $oper }}, r)
			if err != nil {
				errorh.HandleHTTPError(c, rc, err, w, r)
				return
			}

		{{else if HasInput . -}}

			// read and parse request body
			var {{GetInputArgName . }} {{GetInputArgType . }}
			err = json.NewDecoder(r.Body).Decode(&{{GetInputArgName . }})
			if err != nil {
				errorh.HandleHTTPError(c, rc, errorh.NewInvalidInputErrorf(1, "Error parsing request body: %s", err), w, r)
				return
			}

		{{end -}}

		{{if RequiresParamValidation . -}}

			// start parameter validation
			validationErrors := []errorh.FieldError{}
		{{end -}}

		{{range .InputArgs -}}

			{{if not (IsCustomArg .) }}
				{{if IsIntArg . -}}
					{{if IsInputArgMandatory $oper . -}}
						{{.Name}}, fieldError := httpparser.ExtractNumber(r, "{{Uncapitalized .Name}}", true)
						if fieldError != nil {
							validationErrors = append(validationErrors, *fieldError)
						}
					{{else -}}
						{{.Name}}, _ := httpparser.ExtractNumber(r, "{{Uncapitalized .Name}}",false)
					{{end -}}
				{{else if IsBoolArg . -}}
					{{if IsInputArgMandatory $oper . -}}
						{{.Name}}, fieldError := httpparser.ExtractBool(r, "{{Uncapitalized .Name}}", true)
						if fieldError != nil {
							validationErrors = append(validationErrors, *fieldError)
						}
					{{else -}}
						{{.Name}}, _ := httpparser.ExtractBool(r, "{{Uncapitalized .Name}}", false)
					{{end -}}
				{{else if IsDateArg . -}}
					{{if IsInputArgMandatory $oper . -}}
						{{.Name}}, fieldError := httpparser.ExtractDate(r, "{{Uncapitalized .Name}}", true)
						if fieldError != nil {
							validationErrors = append(validationErrors, *fieldError)
						}
					{{else -}}
						{{.Name}}, _ := httpparser.ExtractDate(r, "{{Uncapitalized .Name}}", false)
					{{end -}}
				{{else if IsStringArg . -}}
					{{if IsInputArgMandatory $oper . -}}
						{{.Name}}, fieldError := httpparser.ExtractString(r, "{{Uncapitalized .Name}}", true)
						if fieldError != nil {
							validationErrors = append(validationErrors, *fieldError)
						}
					{{else -}}
						{{.Name}}, _ := httpparser.ExtractString(r, "{{Uncapitalized .Name}}", false)
					{{end -}}
				{{else if IsStringSliceArg . -}}
					{{if IsInputArgMandatory $oper . -}}
						{{.Name}}, fieldError := httpparser.ExtractStringSlice(r, "{{Uncapitalized .Name}}", true)
						if fieldError != nil {
							validationErrors = append(validationErrors, *fieldError)
						}
					{{else -}}
						{{.Name}}, _ := httpparser.ExtractStringSlice(r, "{{Uncapitalized .Name}}", false)
					{{end -}}
				{{else}}
					Force compile error: Input arg {{.}} has unsupported primitive type
				{{end -}}
			{{end -}}
		{{end -}}

		{{if RequiresParamValidation . -}}

			if len(validationErrors) > 0 {
				errorh.HandleHTTPError(c, rc, errorh.NewInvalidInputErrorSpecific(0, validationErrors), w, r)
				return
			}
			// end of parameter validation

		{{end -}}

		// call business logic
		rc.Set(request.Transactional({{ IsRestOperationTransactional $service .}}))
		{{range GetOutputArgsDeclaration . -}}
			{{.}}
		{{end -}}
		{{if IsRestOperationTransactional $service . -}}
		err = eventStore.RunInTransaction(c, rc, func(c context.Context) error {
		{{end -}}
		{{if HasMetaOutput . -}}
			result, meta, err = service.{{$oper.Name}}({{GetInputParamString . }})
		{{else if HasOutput . -}}
			result, err = service.{{$oper.Name}}({{GetInputParamString . }})
		{{else -}}
			err = service.{{$oper.Name}}({{GetInputParamString . }})
		{{end -}}
		{{if IsRestOperationTransactional $service . -}}
			if err != nil {
				return err
			}
			return nil
		})
		{{end -}}
		{{if HasMetaOutput . -}}
			if meta != nil {
				{{if IsMetaCallback . -}}
					metaErr := meta(c, w, r)
				{{else -}}
					metaErr := service.{{$oper.Name}}HandleMetaData(c, w, meta)
				{{end -}}
				if metaErr != nil {
					if err != nil {
						metaErr = err
					}
					errorh.HandleHTTPError(c, rc, metaErr, w, r)
					return
				}
			}
		{{end -}}
		if err != nil {
			errorh.HandleHTTPError(c, rc, err, w, r)
			return
		}

	   {{if HasRestOperationAfter . -}}
			err = service.{{$oper.Name}}HandleAfter(c, r.Method, r.URL, {{GetInputArgName . }}, result)
			if err != nil {
				errorh.HandleHTTPError(c, rc, err, w, r)
				return
			}
		{{end -}}

		{{if NeedsContext $oper}}
			postLogicHook(c, w, r, rc)
		{{else -}}
			postLogicHook(nil, w, r, rc)
		{{end -}}

		// write OK response body
		{{if HasContentType . -}}
			w.Header().Set("Content-Type", "{{GetContentType .}}")
		{{end -}}
		{{if IsRestOperationJSON . -}}
			{{if HasOutput . -}}
				err = json.NewEncoder(w).Encode(result)
				if err != nil {
					mylog.New().Warning(c, "Error writing json-response: %s", err)
				}
			{{end -}}
		{{else if IsRestOperationHTML . -}}
			{{if HasOutput . -}}
				err = service.{{$oper.Name}}WriteHTML(w, result)
				if err != nil {
					mylog.New().Warning(c, "Error writing html-response: %s", err)
				}
			{{else -}}
				err = service.{{$oper.Name}}WriteHTML(w)
				if err != nil {
					mylog.New().Warning(c, "Error writing html-response: %s", err)
				}
			{{end -}}
		{{else if IsRestOperationCSV . -}}
			w.Header().Set("Content-Disposition", "attachment;filename={{ GetRestOperationFilename .}}")
			{{if HasOutput . -}}
				service.{{$oper.Name}}WriteCSV(w, result)
			{{else -}}
				{{$oper.Name}}WriteCSV(w)
			{{end -}}
		{{else if IsRestOperationTXT . -}}
			fmt.Fprint(w, result)
		{{else if IsRestOperationMD . -}}
			fmt.Fprint(w, result)
		{{else if IsRestOperationNoContent . -}}
			w.WriteHeader(http.StatusNoContent)
		{{else if IsRestOperationCustom . -}}
			service.{{$oper.Name}}HandleResult({{GetContextName $oper }}, w, r, result)
		{{else -}}
			errorh.NewInternalErrorf(0, "Not implemented")
		{{end -}}
	}
}
	{{else -}}

// {{$oper.Name}} does the http handling for business logic method service.{{$oper.Name}}
func {{$oper.Name}}(service *{{$service.Name}}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		{{if NeedsContext $oper -}}
			{{GetContextName $oper}} := ctx.New.CreateContext(r)
		{{end -}}
		service.{{$oper.Name}}({{GetInputParamString . }})
	}
}
		{{end -}}
	{{end -}}
{{end}}
`
