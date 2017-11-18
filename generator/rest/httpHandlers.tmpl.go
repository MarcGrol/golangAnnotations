package rest

const httpHandlersTemplate = `// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}

import (
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

{{ $serviceName := .Name }}

var (
    preLogicHook  = func(c context.Context, w http.ResponseWriter, r *http.Request) {}
    postLogicHook = func(c context.Context, w http.ResponseWriter, r *http.Request, credentials rest.Credentials) {}
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

{{ $extractCredentialsMethod := GetExtractCredentialsMethod . }}
{{ $noValidation := IsRestServiceNoValidation . }}

{{range $idxOper, $oper := .Operations}}
    {{if IsRestOperation $oper -}}
        {{if IsRestOperationGenerated . -}}

// {{$oper.Name}} does the http handling for business logic method service.{{$oper.Name}}
func {{$oper.Name}}( service *{{$serviceName}} ) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var err error

        {{if NeedsContext $oper -}}
			{{GetContextName $oper}} := ctx.New.CreateContext(r)
			preLogicHook( c, w, r )
        {{else -}}
			preLogicHook( nil, w, r )
        {{end -}}

        credentials := {{ $extractCredentialsMethod }}(c, r)
        {{if (not $noValidation) and (HasCredentials $oper) -}}
        	err = validateCredentials(c, credentials, {{GetRestOperationRolesString $oper}})
        	if err != nil {
            	rest.HandleHttpError(c, credentials, err, w, r)
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
				rest.HandleHttpError(c, credentials, errorh.NewInvalidInputErrorSpecific(0, validationErrors), w, r)
				return
			}
		{{end -}}

		{{if HasUpload . -}}
			{{GetInputArgName . }}, err := service.{{$oper.Name}}GetUpload({{GetContextName $oper }}, r)
			if err != nil {
				rest.HandleHttpError(c, credentials, err, w, r)
				return
			}
		{{else if HasInput . -}}

			// read and parse request body
			var {{GetInputArgName . }} {{GetInputArgType . }}
			err = json.NewDecoder(r.Body).Decode( &{{GetInputArgName . }} )
			if err != nil {
				rest.HandleHttpError(c, credentials, errorh.NewInvalidInputErrorf(1, "Error parsing request body: %s", err), w, r)
				return
			}
		{{end -}}

        {{if HasMetaOutput . -}}
	        // call business logic
        	result, meta, err := service.{{$oper.Name}}({{GetInputParamString . }})
        {{else if HasOutput . -}}
	        // call business logic
        	result, err := service.{{$oper.Name}}({{GetInputParamString . }})
        {{else -}}
	        // call business logic
        	err = service.{{$oper.Name}}({{GetInputParamString . }})
        {{end -}}
        if err != nil {
            rest.HandleHttpError(c, credentials, err, w, r)
            return
        }
        {{if HasMetaOutput . -}}
			if meta != nil {
				err = service.{{$oper.Name}}HandleMetaData(c, w, meta)
				if err != nil {
					rest.HandleHttpError(c, credentials, err, w, r)
					return
				}
			}
        {{end -}}

        {{if HasRestOperationAfter . -}}
			err = service.{{$oper.Name}}HandleAfter(c, r.Method, r.URL, {{GetInputArgName . }}, result)
			if err != nil {
				rest.HandleHttpError(c, credentials, err, w, r)
				return
			}
        {{end -}}

        {{if NeedsContext $oper -}}
        	postLogicHook( c, w, r, credentials )
        {{else -}}
        	postLogicHook( nil, w, r, credentials )
        {{end -}}

        // write OK response body
        {{if HasContentType . -}}
        	w.Header().Set("Content-Type", "{{GetContentType .}}")
        {{end -}}
        {{if IsRestOperationJSON . -}}
            {{if HasOutput . -}}
				err = json.NewEncoder(w).Encode(result)
				if err != nil {
					log.Printf("Error encoding response payload %+v", err)
				}
			{{end -}}
		{{else if IsRestOperationHTML . -}}
			{{if HasOutput . -}}
				err = service.{{$oper.Name}}WriteHTML(w, result)
			{{else -}}
				err = service.{{$oper.Name}}WriteHTML(w)
			{{end -}}
			if err != nil {
				log.Printf("Error encoding response payload %+v", err)
			}
		{{else if IsRestOperationCSV . -}}
        	w.Header().Set("Content-Disposition", "attachment;filename={{ GetRestOperationFilename .}}")
        	{{if HasOutput . -}}
				err = service.{{$oper.Name}}WriteCSV(w, result)
			{{else -}}
				err = {{$oper.Name}}WriteCSV(w)
			{{end -}}
			if err != nil {
				log.Printf("Error encoding response payload %+v", err)
			}
		{{else if IsRestOperationTXT . -}}
			_, err = fmt.Fprint(w, result)
			if err != nil {
				log.Printf("Error encoding response payload %+v", err)
			}
		{{else if IsRestOperationMD . -}}
			_, err = fmt.Fprint(w, result)
			if err != nil {
				log.Printf("Error encoding response payload %+v", err)
			}
		{{else if IsRestOperationNoContent . -}}
			w.WriteHeader(http.StatusNoContent)
		{{else if IsRestOperationCustom . -}}
			service.{{$oper.Name}}HandleResult({{GetContextName $oper }}, w, r, result)
		{{else -}}
			errorh.NewInternalErrorf(0, "Not implemented")
		{{end -}}// call business logic
    }
}
    {{else -}}

// {{$oper.Name}} does the http handling for business logic method service.{{$oper.Name}}
func {{$oper.Name}}( service *{{$serviceName}} ) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        {{if NeedsContext $oper -}}
			{{GetContextName $oper}} := ctx.New.CreateContext(r)
		{{end -}}
        service.{{$oper.Name}}({{GetInputParamString . }})
    }
}

        {{end -}}
    {{end -}}
{{end -}}
`
