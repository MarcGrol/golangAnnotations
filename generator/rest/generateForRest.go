package rest

import (
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/MarcGrol/golangAnnotations/annotation"
	"github.com/MarcGrol/golangAnnotations/generator/generationUtil"
	"github.com/MarcGrol/golangAnnotations/generator/rest/restAnnotation"
	"github.com/MarcGrol/golangAnnotations/model"
)

func Generate(inputDir string, parsedSource model.ParsedSources) error {
	return generate(inputDir, parsedSource.Structs)
}

func generate(inputDir string, structs []model.Struct) error {
	restAnnotation.Register()

	packageName, err := generationUtil.GetPackageName(structs)
	if err != nil {
		return err
	}
	targetDir, err := generationUtil.DetermineTargetPath(inputDir, packageName)
	if err != nil {
		return err
	}
	for _, service := range structs {
		if IsRestService(service) {
			{
				target := fmt.Sprintf("%s/http%s.go", targetDir, service.Name)
				err = generationUtil.GenerateFileFromTemplate(service, "handlers", handlersTemplate, customTemplateFuncs, target)
				if err != nil {
					log.Fatalf("Error generating handlers for service %s: %s", service.Name, err)
					return err
				}
			}
			{
				target := fmt.Sprintf("%s/http%sHelpers_test.go", targetDir, service.Name)
				err = generationUtil.GenerateFileFromTemplate(service, "helpers", HelpersTemplate, customTemplateFuncs, target)
				if err != nil {
					log.Fatalf("Error generating helpers for service %s: %s", service.Name, err)
					return err
				}
			}

		}
	}
	return nil
}

var customTemplateFuncs = template.FuncMap{
	"IsRestService":          IsRestService,
	"GetRestServicePath":     GetRestServicePath,
	"IsRestOperation":        IsRestOperation,
	"GetRestOperationPath":   GetRestOperationPath,
	"GetRestOperationMethod": GetRestOperationMethod,
	"HasInput":               HasInput,
	"GetInputArgType":        GetInputArgType,
	"GetInputArgName":        GetInputArgName,
	"GetInputParamString":    GetInputParamString,
	"GetOutputArgType":       GetOutputArgType,
	"HasOutput":              HasOutput,
	"IsPrimitive":            IsPrimitive,
	"IsNumber":               IsNumber,
	"HasContext":             HasContext,
	"GetContextName":         GetContextName,
}

func IsRestService(s model.Struct) bool {
	annotation, ok := annotation.ResolveAnnotations(s.DocLines)
	if !ok || annotation.Name != "RestService" {
		return false
	}
	return ok
}

func GetRestServicePath(o model.Struct) string {
	val, ok := annotation.ResolveAnnotations(o.DocLines)
	if ok {
		return val.Attributes["path"]
	}
	return ""
}

func IsRestOperation(o model.Operation) bool {
	annotation, ok := annotation.ResolveAnnotations(o.DocLines)
	if !ok || annotation.Name != "RestOperation" {
		return false
	}
	return ok
}

func GetRestOperationPath(o model.Operation) string {
	val, ok := annotation.ResolveAnnotations(o.DocLines)
	if ok {
		return val.Attributes["path"]
	}
	return ""
}

func GetRestOperationMethod(o model.Operation) string {
	val, ok := annotation.ResolveAnnotations(o.DocLines)
	if ok {
		return val.Attributes["method"]
	}
	return ""
}

func HasInput(o model.Operation) bool {
	if GetRestOperationMethod(o) == "POST" || GetRestOperationMethod(o) == "PUT" {
		for _, arg := range o.InputArgs {
			if arg.TypeName != "context.Context" {
				return true
			}
		}
	}
	return false
}

func HasContext(o model.Operation) bool {
	for _, arg := range o.InputArgs {
		if arg.TypeName == "context.Context" {
			return true
		}
	}
	return false
}

func GetContextName(o model.Operation) string {
	for _, arg := range o.InputArgs {
		if arg.TypeName == "context.Context" {
			return arg.Name
		}
	}
	return ""
}

func GetInputArgType(o model.Operation) string {
	for _, arg := range o.InputArgs {
		if arg.TypeName != "int" && arg.TypeName != "string" && arg.TypeName != "context.Context" {
			return arg.TypeName
		}
	}
	return ""
}

func GetInputArgName(o model.Operation) string {
	for _, arg := range o.InputArgs {
		if arg.TypeName != "int" && arg.TypeName != "string" && arg.TypeName != "context.Context" {
			return arg.Name
		}
	}
	return ""
}

func GetInputParamString(o model.Operation) string {
	args := []string{}
	for _, arg := range o.InputArgs {
		args = append(args, arg.Name)
	}
	return strings.Join(args, ",")
}

func HasOutput(o model.Operation) bool {
	for _, arg := range o.OutputArgs {
		if arg.TypeName != "error" {
			return true
		}
	}
	return false
}

func GetOutputArgType(o model.Operation) string {
	for _, arg := range o.OutputArgs {
		if arg.TypeName != "error" {
			return arg.TypeName
		}
	}
	return ""
}

func IsPrimitive(f model.Field) bool {
	return f.TypeName == "int" || f.TypeName == "string"
}

func IsNumber(f model.Field) bool {
	return f.TypeName == "int"
}

var handlersTemplate string = `
// Generated automatically: do not edit manually

package {{.PackageName}}

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Duxxie/platform/lib/ctx"
	"github.com/MarcGrol/microgen/lib/myerrors"
	"github.com/gorilla/mux"
)

{{ $structName := .Name }}

func (ts *{{.Name}}) HttpHandler() http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	return ts.HttpHandlerWithRouter(router)
}

func (ts *{{.Name}}) HttpHandlerWithRouter(router *mux.Router) *mux.Router {
	subRouter := router.PathPrefix("{{GetRestServicePath . }}").Subrouter()

	{{range .Operations}}
		{{if IsRestOperation . }}
			subRouter.HandleFunc(  "{{GetRestOperationPath . }}", {{.Name}}(ts)).Methods("{{GetRestOperationMethod . }}")
		{{end}}
	{{end}}
	return router
}

{{range $idxOper, $oper := .Operations}}

{{if IsRestOperation $oper}}
func {{$oper.Name}}( service *{{$structName}} ) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		{{if HasContext $oper }}
			{{GetContextName $oper }} := ctx.New.CreateContext(r)
		{{end}}
		pathParams := mux.Vars(r)
		log.Printf("pathParams:%+v", pathParams)

		// extract url-params
		{{range .InputArgs}}
			{{if IsPrimitive . }}
				{{if IsNumber . }}
					{{.Name}}String, exists := pathParams["{{.Name}}"]
					if !exists {
						handleError(myerrors.NewInvalidInputError(fmt.Errorf("Missing path param '{{.Name}}'")), w)
						return
					}
					{{.Name}}, err := strconv.Atoi({{.Name}}String)
					if err != nil {
						handleError(myerrors.NewInvalidInputError(fmt.Errorf("Invalid path param '{{.Name}}'")), w)
						return
					}
				{{else}}
					{{.Name}}, exists := pathParams["{{.Name}}"]
					if !exists {
						handleError(myerrors.NewInvalidInputError(fmt.Errorf("Missing path param '{{.Name}}'")), w)
						return
					}
				{{end}}
			{{end}}
		{{end}}

		{{if HasInput . }}
			// read abd parse request body
			var {{GetInputArgName . }} {{GetInputArgType . }}
			err = json.NewDecoder(r.Body).Decode( &{{GetInputArgName . }} )
			if err != nil {
				handleError(myerrors.NewInvalidInputError(fmt.Errorf("Error decoding request payload:%s", err)), w)
				return
			}
		{{end}}

		// call business logic
		{{if HasOutput . }}
			result, err := service.{{$oper.Name}}({{GetInputParamString . }})
		{{else}}
			err = service.{{$oper.Name}}({{GetInputParamString . }})
		{{end}}
		if err != nil {
			handleError(err, w)
			return
		}

		// write response body
		{{if HasOutput . }}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(result)
			if err != nil {
				log.Printf("Error encoding response payload %+v", err)
			}
		{{else}}
			w.WriteHeader(http.StatusNoContent)
		{{end}}
      }
 }
{{end}}
{{end}}


func handleError(err error, w http.ResponseWriter) {
	errorBody := struct {
		ErrorMessage string
	}{
		err.Error(),
	}
	blob, err := json.Marshal(errorBody)
	if err != nil {
		log.Printf("Error marshalling error response payload %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(determineHttpCode(err))
	w.Header().Set("Content-Type", "application/json")
	w.Write(blob)
}

func determineHttpCode(err error) int {
	if myerrors.IsNotFoundError(err) {
		return http.StatusNotFound
	} else if myerrors.IsInternalError(err) {
		return http.StatusInternalServerError
	} else if myerrors.IsInvalidInputError(err) {
		return http.StatusBadRequest
	} else if myerrors.IsNotAuthorizedError(err) {
		return http.StatusForbidden
	} else {
		return http.StatusInternalServerError
	}
}

`

var HelpersTemplate string = `
// Generated automatically: do not edit manually

package {{.PackageName}}

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
)

{{ $structName := .Name }}

{{range .Operations}}

{{if IsRestOperation . }}
func {{.Name}}TestHelper(url string {{if HasInput . }}, input {{GetInputArgType . }} {{end}} )  (int {{if HasOutput . }},*{{GetOutputArgType . }}{{end}},error) {

	recorder := httptest.NewRecorder()

	{{if HasInput . }}
		requestBody, _ := json.Marshal(input)
		req, err := http.NewRequest("{{GetRestOperationMethod . }}", url, strings.NewReader(string(requestBody)))
	{{else}}
		req, err := http.NewRequest("{{GetRestOperationMethod . }}", url, nil)
	{{end}}
	if err != nil {
		{{if HasOutput . }}
			return 0, nil, err
		{{else}}
			return 0,  err
		{{end}}
	}
	{{if HasOutput . }}
		req.Header.Set("Accept", "application/json")
	{{end}}

	webservice := {{$structName}}{}
	webservice.HttpHandler().ServeHTTP(recorder, req)

	{{if HasOutput . }}
		var resp {{GetOutputArgType . }}
		dec := json.NewDecoder(recorder.Body)
		err = dec.Decode(&resp)
		if err != nil {
			return recorder.Code, nil, err
		}
		return recorder.Code, &resp, nil
	{{else}}
		return recorder.Code, nil
	{{end}}
}
{{end}}
{{end}}
`
