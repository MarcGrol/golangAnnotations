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
				err = generationUtil.GenerateFileFromTemplate(service, fmt.Sprintf("%s.%s", service.PackageName, service.Name), "handlers", handlersTemplate, customTemplateFuncs, target)
				if err != nil {
					log.Fatalf("Error generating handlers for service %s: %s", service.Name, err)
					return err
				}
			}
			{
				target := fmt.Sprintf("%s/http%sHelpers_test.go", targetDir, service.Name)
				err = generationUtil.GenerateFileFromTemplate(service, fmt.Sprintf("%s.%s", service.PackageName, service.Name), "helpers", HelpersTemplate, customTemplateFuncs, target)
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
	"ExtractImports":         ExtractImports,
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
	"WithBackTicks": SurroundWithBackTicks,
	"BackTick": BackTick,
}

func BackTick() string{
	return "`"
}

func SurroundWithBackTicks(body string) string{
	return fmt.Sprintf("`%s'", body)
}

func IsRestService(s model.Struct) bool {
	annotation, ok := annotation.ResolveAnnotations(s.DocLines)
	if !ok || annotation.Name != "RestService" {
		return false
	}
	return ok
}

func isImportToBeIgnored( imp string ) bool {
	if imp == "" {
		return true
	}
	for _, i := range []string{
		"golang.org/x/net/context",
		"github.com/gorilla/mux",
	} {
		if imp == i {
			return true
		}
	}
	return false
}

func ExtractImports(s model.Struct) []string {
	importsMap := map[string]string{}
	for _, o := range s.Operations {
		for _, ia := range o.InputArgs {
			if isImportToBeIgnored(ia.PackageName ) == false {
				importsMap[ia.PackageName] = ia.PackageName
			}
		}
		for _, oa := range o.OutputArgs {
			if isImportToBeIgnored(oa.PackageName ) == false {
				importsMap[oa.PackageName] = oa.PackageName
			}
		}
	}
	importsList := []string{}
	for _, v := range importsMap {
		importsList = append(importsList, v)
	}

	return importsList
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
				return fmt.Sprintf("[]%s%s{}", pointer, arg.TypeName)

			} else {
				return fmt.Sprintf("%s%s{}", addressOf, arg.TypeName)
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
// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}

import (
	"encoding/json"
	"log"
	{{if NeedsIntegerConversion .}}
		"strconv"
	{{end}}

	{{if NeedsContext .}}
		"github.com/Duxxie/platform/backend/lib/ctx"
	{{end}}

	"github.com/MarcGrol/golangAnnotations/generator/rest/errorh"
	"github.com/gorilla/mux"

	{{if HasOperationsWithInput .}}
		{{range ExtractImports .}}
			"{{.}}"
		{{end}}
	{{else}}
		"net/http"
	{{end}}

)

{{ $structName := .Name }}

// HTTPHandler registers endpoint in new router
func (ts *{{.Name}}) HTTPHandler() http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	return ts.HTTPHandlerWithRouter(router)
}

// HTTPHandlerWithRouter registers endpoint in existing router
func (ts *{{.Name}}) HTTPHandlerWithRouter(router *mux.Router) *mux.Router {
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
// {{$oper.Name}} does the http handling for business logic method service.{{$oper.Name}}
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
					language := "nl"
					langCookie, err := r.Cookie("lang")
					if err == nil {
						language = langCookie.Value
					}
					authContext := map[string]string {
						"sessionUid": r.Header.Get("X-session-uid"),
						"enduserRole": r.Header.Get("X-enduser-role"),
						"enduserUid": r.Header.Get("X-enduser-uid"),
						"language": language,
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
	enduserUID, found := authContext["enduserUid"]
	if found == false || enduserUID == "" {
		return "", "", "", errorh.NewNotAuthorizedErrorf(0, "Missing/invalid enduser-uid %s", enduserUID)
	}
	sessionUID, found := authContext["sessionUid"]
	if found == false || sessionUID == "" {
		return "", "", "", errorh.NewNotAuthorizedErrorf(0, "Missing/invalid session-uid %s", sessionUID)
	}

	return role, enduserUID, sessionUID, nil
}
{{end}}

`

var HelpersTemplate string = `
// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}

import (
	"encoding/json"
	"net/http/httptest"
	"net/http/httputil"
	"github.com/MarcGrol/golangAnnotations/generator/rest/errorh"
	"os"
	"testing"
	"fmt"
	"log"
	"testing"

	"strings"
	"bytes"
	{{range ExtractImports .}}
		"{{.}}"
	{{end}}
)

{{ $structName := .Name }}


var fp *os.File
var logFp *os.File


func openfile( filename string) *os.File {
	fp, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Error opening rest-dump-file %s: %s", filename, err.Error())
	}
	return fp
}

func TestMain(m *testing.M) {

	fp = openfile("test_results_{{.Name}}.md")
	defer func() {
		fp.Close()
	}()
	fmt.Fprintf(fp, "# HTTP request and response payloads for all test in package\n\n" )

	logFp = openfile("testResultsOf{{.Name}}.go")
	defer func() {
		logFp.Close()
	}()
	fmt.Fprintf(logFp, "package {{.PackageName}}\n\n" )
	fmt.Fprintf(logFp, "// Generated automatically based on running of api-tests\n\n" )
	fmt.Fprintf(logFp, "import \"github.com/Duxxie/platform/backend/lib/testcase\"\n")

	fmt.Fprintf(logFp, "var TestResults = testcase.TestSuiteDescriptor {\n" )
	fmt.Fprintf(logFp, "\tTestCases: []testcase.TestCaseDescriptor{\n")
	code := m.Run()

	fmt.Fprintf(logFp, "},\n" )
	fmt.Fprintf(logFp, "}\n" )

	os.Exit(code)
}

func testCase(name string, description string) {
	fmt.Fprintf(fp, "\n---\n\n## Test-case '%s'\n\n%s\n\n", name, description)

	fmt.Fprintf(logFp, "\t\ttestcase.TestCaseDescriptor{\n")
	fmt.Fprintf(logFp, "\t\tName:\"%s\",\n", name)
	fmt.Fprintf(logFp, "\t\tDescription:\"%s\",\n", description)
}

func testCaseDone() {
	fmt.Fprintf(logFp, "},\n")
}


{{range .Operations}}

{{if IsRestOperation . }}

func {{.Name}}TestHelper(url string {{if HasInput . }}, input {{GetInputArgType . }} {{end}} )  (int {{if HasOutput . }},{{GetOutputArgType . }}{{end}},*errorh.Error,error) {
	return {{.Name}}TestHelperWithHeaders( url {{if HasInput . }}, input {{end}}, map[string]string{} )
}

func {{.Name}}TestHelperWithHeaders(url string {{if HasInput . }}, input {{GetInputArgType . }} {{end}}, headers map[string]string)  (int {{if HasOutput . }},{{GetOutputArgType . }}{{end}},*errorh.Error,error) {
	fmt.Fprintf(fp, "### Operation %s\n", "{{.Name}}")

	fmt.Fprintf(logFp, "\t\tOperation:\"%s\",\n", "{{.Name}}")
	defer func() {
		fmt.Fprintf(logFp, "\t},\n")
	}()

	recorder := httptest.NewRecorder()

	{{if HasInput . }}
		rb, _ := json.Marshal(input)
		// indent for readability
		var requestBody bytes.Buffer
		json.Indent(&requestBody, rb, "", "\t")

		req, err := http.NewRequest("{{GetRestOperationMethod . }}", url, strings.NewReader(requestBody.String()))
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

	fmt.Fprintf(logFp, "\tRequest: testcase.RequestDescriptor{\n")
	fmt.Fprintf(logFp, "\tMethod:\"%s\",\n", "{{GetRestOperationMethod . }}")
	fmt.Fprintf(logFp, "\tUrl:\"%s\",\n", url)
	fmt.Fprintf(logFp, "\tHeaders:[]string{\"TODO\"},\n")
	{{if HasInput . }}
		fmt.Fprintf(logFp, "\tBody:\n" )
		fmt.Fprintf(logFp, "{{BackTick}}%s{{BackTick}}", requestBody.String() )
		fmt.Fprintf(logFp, ",\n" )
	{{end}}
	fmt.Fprintf(logFp, "},\n")

	// dump readable request
	payload, err := httputil.DumpRequest(req, true)
	if err == nil {
		fmt.Fprintf(fp, "\n### http-request:\n\n    %s\n\n    ", strings.Replace(string(payload), "\n", "\n    ", -1))
	}

	fmt.Fprintf(logFp, "\tResponse:testcase.ResponseDescriptor{\n")
	defer func() {
		fmt.Fprintf(logFp, "\t},\n")
	}()


	webservice := {{$structName}}{}
	webservice.HTTPHandler().ServeHTTP(recorder, req)


    // dump readable response
	var responseBody bytes.Buffer
	json.Indent(&responseBody, recorder.Body.Bytes(), "", "\t")
	fmt.Fprintf(fp, "\n\n### http-response:\n\n    %d\n\n    %s\n\n", recorder.Code, strings.Replace(responseBody.String(), "\n", "\n    ", -1))

	fmt.Fprintf(logFp, "\tStatus:%d,\n", recorder.Code)
	fmt.Fprintf(logFp, "\tHeaders:[]string{\"TODO\"},\n")
	fmt.Fprintf(logFp, "\tBody:\n{{BackTick}}%s{{BackTick}},\n", responseBody.String())

	{{if HasOutput . }}
		if recorder.Code != http.StatusOK {
		    // return error response
			var errorResp errorh.Error
			dec := json.NewDecoder(recorder.Body)
			err = dec.Decode(&errorResp)
			if err != nil {
				return recorder.Code, nil, nil, err
			}
			return recorder.Code, nil, &errorResp, nil
		}

		// return success response
		resp := {{GetOutputArgDeclaration . }}
		dec := json.NewDecoder(recorder.Body)
		err = dec.Decode({{GetOutputArgName . }})
		if err != nil {
			return recorder.Code, nil, nil, err
		}
		return recorder.Code, resp, nil, nil

	{{else}}
		return recorder.Code, nil, nil
	{{end}}
}
{{end}}
{{end}}
`
