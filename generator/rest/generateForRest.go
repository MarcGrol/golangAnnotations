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
	"NeedsIntegerConversion": NeedsIntegerConversion,
	"NeedsContext":           NeedsContext,
	"GetRestServicePath":     GetRestServicePath,
	"IsRestOperation":        IsRestOperation,
	"GetRestOperationPath":   GetRestOperationPath,
	"GetRestOperationMethod": GetRestOperationMethod,
	"HasOperationsWithInput": HasOperationsWithInput,
	"HasInput":               HasInput,
	"GetInputArgType":        GetInputArgType,
	"UsesQueryParams":        UsesQueryParams,
	"GetInputArgName":        GetInputArgName,
	"GetInputParamString":    GetInputParamString,
	"GetOutputArgType":       GetOutputArgType,
	"HasOutput":              HasOutput,
	"IsPrimitive":            IsPrimitive,
	"IsNumber":               IsNumber,
	"IsInputArgMandatory": 	  IsInputArgMandatory,
	"IsAuthContextArg":       IsAuthContextArg,
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

func NeedsIntegerConversion(s model.Struct) bool {
	for _, oper := range s.Operations {
		for _, arg := range oper.InputArgs {
			if arg.TypeName == "int" {
				return true
			}
		}
	}
	return false
}

func NeedsContext(s model.Struct) bool {
	for _, oper := range s.Operations {
		for _, arg := range oper.InputArgs {
			if arg.TypeName == "context.Context" {
				return true
			}
		}
	}
	return false
}

func GetRestServicePath(s model.Struct) string {
	val, ok := annotation.ResolveAnnotations(s.DocLines)
	if ok {
		return val.Attributes["path"]
	}
	return ""
}

func HasOperationsWithInput(s model.Struct) bool {
	for _,o := range s.Operations {
		if HasInput(*o) == true {
			return true
		}
	}
	return false
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


func UsesQueryParams(o model.Operation) bool {
	if GetRestOperationMethod(o) == "GET"  {
		count := 0
		for _,arg := range o.InputArgs {
			if arg.TypeName != "context.Context" && arg.Name != "authContext" {
				count++
			}
		}
		return count > 1
	}
	return false
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


func findArgInArray( array []string, toMatch string ) bool {
	for _,p := range array {
		log.Printf("arg:%v matches: %v", p,toMatch)
		if strings.Trim(p, " ") == toMatch {
			return true
		}
	}
	return false
}

func IsInputArgMandatory(o model.Operation, arg model.Field) bool {
	annotation, ok := annotation.ResolveAnnotations(o.DocLines)
	if !ok || annotation.Name != "RestOperation" {
		return false
	}
	optionalArgsString, ok :=  annotation.Attributes["optionalargs"]
	if !ok {
		return true
	}

	return !findArgInArray(strings.Split(optionalArgsString, ","),arg.Name)
}


func IsAuthContextArg(arg model.Field) bool {
	return arg.Name == "authContext" && arg.TypeName == "map[string]string"
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
	"log"
	"net/http"
	{{if NeedsIntegerConversion .}}"strconv"{{end}}

	{{if NeedsContext .}}"github.com/Duxxie/platform/backend/lib/ctx"{{end}}
	"github.com/MarcGrol/golangAnnotations/generator/rest/errorh"
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

		{{if UsesQueryParams $oper }} {{else}}
		pathParams := mux.Vars(r)
		log.Printf("pathParams:%+v", pathParams)
		{{end}}

		// extract url-params
	    validationErrors := []errorh.FieldError{}
		{{range .InputArgs}}
			{{if IsPrimitive . }}
				{{if IsNumber . }}
					{{.Name}} := 0
					{{if UsesQueryParams $oper }}
					{{.Name}}String := r.URL.Query().Get("{{.Name}}")
					if {{.Name}}String == "" {
					{{else}}
					{{.Name}}String, exists := pathParams["{{.Name}}"]
					if !exists {
					{{end}}
					{{if IsInputArgMandatory $oper .}}
						validationErrors = append(validationErrors, errorh.FieldError{
						SubCode: 1000,
						Field:   "{{.Name}}",
						Msg:     "Missing value for mandatory parameter %s",
						Args:    []string{"{{.Name}}"},
					 })
					 {{end}}
					} else {
					{{if IsInputArgMandatory $oper .}}
					{{.Name}}, err = strconv.Atoi({{.Name}}String)
					if err != nil {
						validationErrors = append(validationErrors, errorh.FieldError{
						SubCode: 1001,
						Field:   "{{.Name}}",
						Msg:     "Invalid value for mandatory parameter %s",
						Args:    []string{"{{.Name}}"},
					 })
					 }
					 }
					 {{end}}
				{{else}}
					{{if UsesQueryParams $oper }}
					{{.Name}} := r.URL.Query().Get("{{.Name}}")
					if {{.Name}} == "" {
					{{else}}
					{{.Name}}, exists := pathParams["{{.Name}}"]
					if !exists {
					{{end}}
					{{if IsInputArgMandatory $oper .}}
						validationErrors = append(validationErrors, errorh.FieldError{
						SubCode: 1000,
						Field:   "{{.Name}}",
						Msg:     "Missing value for mandatory parameter %s",
						Args:    []string{"{{.Name}}"},
					 })
					 {{end}}
					}
					{{end}}
				{{end}}
				{{if IsAuthContextArg .}}
				authContext := map[string]string {
					"enduserRole": r.Header.Get("X-enduser-role"),
					"enduserUid": r.Header.Get("X-enduser-uid"),
				}
				{{else}}
			{{end}}
		{{end}}

        if len(validationErrors) > 0 {
            errorh.HandleHttpError(errorh.NewInvalidInputErrorSpecific(0, validationErrors), w)
            return
        }

		{{if HasInput . }}
			// read abd parse request body
			var {{GetInputArgName . }} {{GetInputArgType . }}
			err = json.NewDecoder(r.Body).Decode( &{{GetInputArgName . }} )
			if err != nil {
         		errorh.HandleHttpError(errorh.NewInvalidInputErrorf(1, "Error psrsing request body: %s", err), w)
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
			errorh.HandleHttpError(err, w)
			return
		}

		// write OK response body
		{{if HasOutput . }}
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

`

var HelpersTemplate string = `
// Generated automatically: do not edit manually

package {{.PackageName}}

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"github.com/MarcGrol/golangAnnotations/generator/rest/errorh"

	{{if HasOperationsWithInput .}}"strings"{{end}}
)

{{ $structName := .Name }}

{{range .Operations}}

{{if IsRestOperation . }}

func {{.Name}}TestHelper(url string {{if HasInput . }}, input {{GetInputArgType . }} {{end}} )  (int {{if HasOutput . }},*{{GetOutputArgType . }}{{end}},*errorh.Error,error) {

	recorder := httptest.NewRecorder()

	{{if HasInput . }}
		requestBody, _ := json.Marshal(input)
		req, err := http.NewRequest("{{GetRestOperationMethod . }}", url, strings.NewReader(string(requestBody)))
	{{else}}
		req, err := http.NewRequest("{{GetRestOperationMethod . }}", url, nil)
	{{end}}
	if err != nil {
		{{if HasOutput . }}
			return 0, nil, nil, err
		{{else}}
			return 0,  nil, err
		{{end}}
	}
	req.RequestURI = url
	{{if HasOutput . }}
		req.Header.Set("Accept", "application/json")
	{{end}}

	webservice := {{$structName}}{}
	webservice.HttpHandler().ServeHTTP(recorder, req)

	{{if HasOutput . }}
		if recorder.Code == http.StatusOK {
			var resp {{GetOutputArgType . }}
			dec := json.NewDecoder(recorder.Body)
			err = dec.Decode(&resp)
			if err != nil {
				return recorder.Code, nil, nil, err
			}
			return recorder.Code, &resp, nil, nil
		} else {
			var errorResp errorh.Error
			dec := json.NewDecoder(recorder.Body)
			err = dec.Decode(&errorResp)
			if err != nil {
				return recorder.Code, nil, nil, err
			}
			return recorder.Code, nil, &errorResp, nil
		}
	{{else}}
		return recorder.Code, nil, nil
	{{end}}
}
{{end}}
{{end}}
`
