package rest

const httpHandlersTemplate = `// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}

import (
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
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
    		subRouter.HandleFunc(  "{{GetRestOperationPath . }}", {{.Name}}(ts)).Methods("{{GetRestOperationMethod . }}")
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
func {{$oper.Name}}( service *{{$service.Name}} ) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var err error

        {{if NeedsContext $oper -}}
			{{GetContextName $oper}} := ctx.New.CreateContext(r)
			preLogicHook( c, w, r )
        {{else -}}
			preLogicHook( nil, w, r )
        {{end -}}

        rc := {{ $extractRequestContextMethod }}Builder(c, r).
		       Transactional({{ IsRestOperationTransactional $service .}}).
               Build()

        {{if (not $noValidation) and (HasRequestContext $oper) -}}
        	err = validateRequestContext(c, rc, {{GetRestOperationRolesString $oper}})
        	if err != nil {
            	errorh.HandleHttpError(c, rc, err, w, r)
            	return
        	}
        {{end -}}

        {{if HasAnyPathParam $oper -}}
        	pathParams := mux.Vars(r)
        	if len(pathParams) > 0 {
            	log.Printf("pathParams:%+v", pathParams)
        	}
        {{end -}}

        {{if RequiresParamValidation . -}}
        	// extract url-params
        	validationErrors := []errorh.FieldError{}
        {{end -}}

        {{range .InputArgs -}}
			{{if IsPrimitiveArg . -}}
				{{if IsNumberArg . -}}
            		{{.Name}} := 0
                    {{if IsRestOperationForm $oper -}}
            			{{.Name}}String := r.FormValue("{{.Name}}")
            			if {{.Name}}String == "" {
                    {{else if IsQueryParam $oper . -}}
                        {{if IsSliceParam . -}}
                    	{{.Name}}String, ok := r.URL.Query()["{{.Name}}"]
                		if !ok {
					{{else -}}
                		{{.Name}}String := r.URL.Query().Get("{{.Name}}")
                		if {{.Name}}String == "" {
					{{end -}}
				{{else -}}
					{{.Name}}String, exists := pathParams["{{.Name}}"]
					if !exists {
				{{end -}}

				{{if IsInputArgMandatory $oper . -}}
					validationErrors = append(validationErrors, errorh.FieldErrorForMissingParameter("{{.Name}}"))
				{{else -}}
					// optional parameter
				{{end -}}
				} else {
					{{.Name}}, err = strconv.Atoi({{.Name}}String)
					if err != nil {
						validationErrors = append(validationErrors, errorh.FieldErrorForInvalidParameter("{{.Name}}"))
					}
				}
				{{else -}}
					{{if IsRestOperationForm $oper -}}
						{{.Name}} := r.FormValue("{{.Name}}")
						if {{.Name}} == "" {
					{{else if IsQueryParam $oper . -}}
						{{if IsSliceParam . -}}
							{{.Name}} := r.URL.Query()["{{.Name}}"]
							if len({{.Name}}) == 0 {
						{{else -}}
							{{.Name}} := r.URL.Query().Get("{{.Name}}")
							if {{.Name}} == "" {
						{{end -}}
					{{else -}}
						{{.Name}}, exists := pathParams["{{.Name}}"]
						if !exists {
					{{end -}}
					{{if IsInputArgMandatory $oper . -}}
						validationErrors = append(validationErrors, errorh.FieldErrorForMissingParameter("{{.Name}}"))
					{{else -}}
						// optional parameter
					{{end -}}
					}
				{{end -}}
			{{end -}}
		{{end -}}

		{{if RequiresParamValidation . -}}
			if len(validationErrors) > 0 {
				errorh.HandleHttpError(c, rc, errorh.NewInvalidInputErrorSpecific(0, validationErrors), w, r)
				return
			}
		{{end -}}

		{{if HasUpload . -}}
			{{GetInputArgName . }}, err := service.{{$oper.Name}}GetUpload({{GetContextName $oper }}, r)
			if err != nil {
				errorh.HandleHttpError(c, rc, err, w, r)
				return
			}
		{{else if HasInput . -}}

			// read and parse request body
			var {{GetInputArgName . }} {{GetInputArgType . }}
			err = json.NewDecoder(r.Body).Decode( &{{GetInputArgName . }} )
			if err != nil {
				errorh.HandleHttpError(c, rc, errorh.NewInvalidInputErrorf(1, "Error parsing request body: %s", err), w, r)
				return
			}
		{{end}}

        // call business logic: transactional: {{ IsRestOperationTransactional $service .}}
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
			publishStoredEnvelopes(c, rc)
			return nil
		})
        {{end -}}
        if err != nil {
            errorh.HandleHttpError(c, rc, err, w, r)
            return
        }
		{{if IsRestOperationTransactional $service . -}}
		{{else -}}
			publishStoredEnvelopes(c, rc)
		{{end -}}

        {{if HasMetaOutput .}}
			if meta != nil {
				{{if IsMetaCallback . -}}
					err = meta(c, w, r)
				{{else -}}
					err = service.{{$oper.Name}}HandleMetaData(c, w, meta)
				{{end -}}
				if err != nil {
					errorh.HandleHttpError(c, rc, err, w, r)
					return
				}
			}
        {{end -}}

       {{if HasRestOperationAfter . -}}
			err = service.{{$oper.Name}}HandleAfter(c, r.Method, r.URL, {{GetInputArgName . }}, result)
			if err != nil {
				errorh.HandleHttpError(c, rc, err, w, r)
				return
			}
        {{end -}}

        {{if NeedsContext $oper}}
        	postLogicHook( c, w, r, rc )
        {{else -}}
        	postLogicHook( nil, w, r, rc )
        {{end -}}

        // write OK response body
        {{if HasContentType . -}}
        	w.Header().Set("Content-Type", "{{GetContentType .}}")
        {{end -}}
        {{if IsRestOperationJSON . -}}
            {{if HasOutput . -}}
				json.NewEncoder(w).Encode(result) 
			{{end -}}
		{{else if IsRestOperationHTML . -}}
			{{if HasOutput . -}}
				service.{{$oper.Name}}WriteHTML(w, result)
			{{else -}}
				service.{{$oper.Name}}WriteHTML(w)
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
func {{$oper.Name}}( service *{{$service.Name}} ) http.HandlerFunc {
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

func publishStoredEnvelopes(c context.Context, rc request.Context) {
	// copy stored envelopes out of request.Context
	storedEnvelopes := rc.GetEnvelopes()
	rc.ClearEnvelopes()

	mylog.New().Info(c, "Publish %d stored envelopes", len(storedEnvelopes))
	for _, e := range storedEnvelopes {
		if rc.IsTransactional() {
			bus.PublishBackground(c, rc, e.AggregateName, e) // async
		} else {
			bus.Publish(c, rc, e.AggregateName, e) //sync
		}
	}
	if len(rc.GetEnvelopes()) > 0 {
		// recurse because sync publisher can also have stored events that need to be published
		publishStoredEnvelopes(c, rc)
	}
}

`
