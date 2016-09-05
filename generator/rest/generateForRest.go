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
	"HasAuthContextArg":      HasAuthContextArg,
	"NeedsIntegerConversion": NeedsIntegerConversion,
	"NeedsContext":           NeedsContext,
	"GetRestServicePath":     GetRestServicePath,
	"IsRestOperation":        IsRestOperation,
	"GetRestOperationPath":   GetRestOperationPath,
	"GetRestOperationMethod": GetRestOperationMethod,
	"HasOperationsWithInput": HasOperationsWithInput,
	"HasInput":               HasInput,
	"GetInputArgType":        GetInputArgType,
	"GetOutputArgDeclaration": GetOutputArgDeclaration,
	"GetOutputArgName":			GetOutputArgName,
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


func HasAuthContextArg(s model.Struct) bool {
	for _, oper := range s.Operations {
		for _, a := range oper.InputArgs {
			if IsAuthContextArg(a) {
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
			slice := ""
			if arg.IsSlice {
				slice = "[]"
			}
			pointer := ""
			if arg.IsPointer {
				pointer = "*"
			}
			return fmt.Sprintf("%s%s%s", slice, pointer, arg.TypeName)
		}
	}
	return ""
}


func GetOutputArgDeclaration(o model.Operation) string {
	for _, arg := range o.OutputArgs {
		if arg.TypeName != "error" {
			pointer := ""
			addressOf := ""
			if arg.IsPointer {
				pointer = "*"
				addressOf = "&"
			}

			if arg.IsSlice {
				return fmt.Sprintf("[]%s%s = []%s%s{}", pointer, arg.TypeName, pointer, arg.TypeName)

			} else {
				return fmt.Sprintf("%s%s = %s%s{}", pointer, arg.TypeName, addressOf, arg.TypeName)
			}
		}
	}
	return ""
}

func GetOutputArgName(o model.Operation) string {
	for _, arg := range o.OutputArgs {
		if arg.TypeName != "error" {
			if !arg.IsPointer || arg.IsSlice {

				return "&resp"
			}
			return "resp"
		}
	}
	return ""
}

func findArgInArray( array []string, toMatch string ) bool {
	for _,p := range array {
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
					 {{else}}
					 // optional parameter
					 {{end}}
					} else {
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
					  	{{else}}
					  		// optional parameter
						 {{end}}
						}
					{{end}}
				{{end}}
				{{if IsAuthContextArg .}}
					authContext := map[string]string {
						"sessionUid": r.Header.Get("X-session-uid"),
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

{{if HasAuthContextArg .}}
func getCredentials(authContext map[string]string, expectedRole string) (string, string, string, error) {
	role, found := authContext["enduserRole"]
	if role != expectedRole {
		return "", "", "", errorh.NewNotAuthorizedErrorf(0, "Missing/invalid role %s", role)
	}
	caregiverUid, found := authContext["enduserUid"]
	if found == false || caregiverUid == "" {
		return "", "", "", errorh.NewNotAuthorizedErrorf(0, "Missing/invalid caregiver-uid %s", caregiverUid)
	}
	sessionUid, found := authContext["sessionUid"]
	if found == false || sessionUid == "" {
		return "", "", "", errorh.NewNotAuthorizedErrorf(0, "Missing/invalid session-uid %s", sessionUid)
	}

	return role, caregiverUid, sessionUid, nil
}
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

func {{.Name}}TestHelper(url string {{if HasInput . }}, input {{GetInputArgType . }} {{end}} )  (int {{if HasOutput . }},{{GetOutputArgType . }}{{end}},*errorh.Error,error) {
	return {{.Name}}TestHelperWithHeaders( url {{if HasInput . }}, input {{end}}, map[string]string{} )
}

func {{.Name}}TestHelperWithHeaders(url string {{if HasInput . }}, input {{GetInputArgType . }} {{end}}, headers map[string]string)  (int {{if HasOutput . }},{{GetOutputArgType . }}{{end}},*errorh.Error,error) {

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
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	webservice := {{$structName}}{}
	webservice.HttpHandler().ServeHTTP(recorder, req)

	{{if HasOutput . }}
		if recorder.Code == http.StatusOK {
			var resp {{GetOutputArgDeclaration . }}
			dec := json.NewDecoder(recorder.Body)
			err = dec.Decode({{GetOutputArgName . }})
			if err != nil {
				return recorder.Code, nil, nil, err
			}
			return recorder.Code, resp, nil, nil
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
