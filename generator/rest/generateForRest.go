package rest

import (
	"fmt"
	"log"
	"strings"
	"text/template"

	"unicode"

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

	packageName, err := generationUtil.GetPackageNameForStructs(structs)
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
				target := fmt.Sprintf("%s/$http%s.go", targetDir, service.Name)
				err = generationUtil.GenerateFileFromTemplate(service, fmt.Sprintf("%s.%s", service.PackageName, service.Name), "handlers", handlersTemplate, customTemplateFuncs, target)
				if err != nil {
					log.Fatalf("Error generating handlers for service %s: %s", service.Name, err)
					return err
				}
			}
			{
				target := fmt.Sprintf("%s/$http%sHelpers_test.go", targetDir, service.Name)
				err = generationUtil.GenerateFileFromTemplate(service, fmt.Sprintf("%s.%s", service.PackageName, service.Name), "helpers", helpersTemplate, customTemplateFuncs, target)
				if err != nil {
					log.Fatalf("Error generating helpers for service %s: %s", service.Name, err)
					return err
				}
			}
			{
				target := fmt.Sprintf("%s/$httpTest%s.go", targetDir, service.Name)
				err = generationUtil.GenerateFileFromTemplate(service, fmt.Sprintf("%s.%s", service.PackageName, service.Name), "testService", testServiceTemplate, customTemplateFuncs, target)
				if err != nil {
					log.Fatalf("Error generating testHandler for service %s: %s", service.Name, err)
					return err
				}
			}
			{
				target := fmt.Sprintf("%s/$httpClientFor%s.go", targetDir, service.Name)
				err = generationUtil.GenerateFileFromTemplate(service, fmt.Sprintf("%s.%s", service.PackageName, service.Name), "httpClient", httpClientTemplate, customTemplateFuncs, target)
				if err != nil {
					log.Fatalf("Error generating httpClient for service %s: %s", service.Name, err)
					return err
				}
			}

		}
	}
	return nil
}

var customTemplateFuncs = template.FuncMap{
	"IsRestService":               IsRestService,
	"ExtractImports":              ExtractImports,
	"GetRestServicePath":          GetRestServicePath,
	"IsRestOperation":             IsRestOperation,
	"GetRestOperationPath":        GetRestOperationPath,
	"GetRestOperationMethod":      GetRestOperationMethod,
	"IsRestOperationForm":         IsRestOperationForm,
	"IsRestOperationJSON":         IsRestOperationJSON,
	"IsRestOperationHTML":         IsRestOperationHTML,
	"IsRestOperationCSV":          IsRestOperationCSV,
	"IsRestOperationTXT":          IsRestOperationTXT,
	"IsRestOperationMD":           IsRestOperationMD,
	"IsRestOperationNoContent":    IsRestOperationNoContent,
	"IsRestOperationCustom":       IsRestOperationCustom,
	"IsRestOperationGenerated":    IsRestOperationGenerated,
	"HasContentType":              HasContentType,
	"GetContentType":              GetContentType,
	"GetRestOperationFilename":    GetRestOperationFilename,
	"GetRestOperationRolesString": GetRestOperationRolesString,
	"HasOperationsWithInput":      HasOperationsWithInput,
	"HasInput":                    HasInput,
	"GetInputArgType":             GetInputArgType,
	"GetOutputArgDeclaration":     GetOutputArgDeclaration,
	"GetOutputArgName":            GetOutputArgName,
	"UsesQueryParams":             UsesQueryParams,
	"GetInputArgName":             GetInputArgName,
	"GetInputParamString":         GetInputParamString,
	"GetOutputArgType":            GetOutputArgType,
	"HasOutput":                   HasOutput,
	"IsPrimitiveArg":              IsPrimitiveArg,
	"IsNumberArg":                 IsNumberArg,
	"RequiresParamValidation":     RequiresParamValidation,
	"IsInputArgMandatory":         IsInputArgMandatory,
	"HasUpload":                   HasUpload,
	"IsUploadArg":                 IsUploadArg,
	"HasCredentials":              HasCredentials,
	"HasContext":                  HasContext,
	"ReturnsError":                ReturnsError,
	"NeedsContext":                NeedsContext,
	"GetContextName":              GetContextName,
	"WithBackTicks":               SurroundWithBackTicks,
	"BackTick":                    BackTick,
	"ToFirstUpper":                toFirstUpper,
}

func BackTick() string {
	return "`"
}

func SurroundWithBackTicks(body string) string {
	return fmt.Sprintf("`%s'", body)
}

func IsRestService(s model.Struct) bool {
	_, ok := annotation.ResolveAnnotationByName(s.DocLines, restAnnotation.TypeRestService)
	return ok
}

func isImportToBeIgnored(imp string) bool {
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
			if isImportToBeIgnored(ia.PackageName) == false {
				importsMap[ia.PackageName] = ia.PackageName
			}
		}
		for _, oa := range o.OutputArgs {
			if isImportToBeIgnored(oa.PackageName) == false {
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

func GetRestServicePath(s model.Struct) string {
	ann, ok := annotation.ResolveAnnotationByName(s.DocLines, restAnnotation.TypeRestService)
	if ok {
		return ann.Attributes[restAnnotation.ParamPath]
	}
	return ""
}

func HasOperationsWithInput(s model.Struct) bool {
	for _, o := range s.Operations {
		if HasInput(*o) == true {
			return true
		}
	}
	return false
}

func IsRestOperation(o model.Operation) bool {
	_, ok := annotation.ResolveAnnotationByName(o.DocLines, restAnnotation.TypeRestOperation)
	return ok
}

func GetRestOperationPath(o model.Operation) string {
	ann, ok := annotation.ResolveAnnotationByName(o.DocLines, restAnnotation.TypeRestOperation)
	if ok {
		return ann.Attributes[restAnnotation.ParamPath]
	}
	return ""
}

func GetRestOperationMethod(o model.Operation) string {
	ann, ok := annotation.ResolveAnnotationByName(o.DocLines, restAnnotation.TypeRestOperation)
	if ok {
		return ann.Attributes[restAnnotation.ParamMethod]
	}
	return ""
}

func IsRestOperationForm(o model.Operation) bool {
	ann, ok := annotation.ResolveAnnotationByName(o.DocLines, restAnnotation.TypeRestOperation)
	if ok {
		return ann.Attributes[restAnnotation.ParamForm] == "true"
	}
	return false
}

func GetRestOperationFormat(o model.Operation) string {
	ann, ok := annotation.ResolveAnnotationByName(o.DocLines, restAnnotation.TypeRestOperation)
	if ok {
		return ann.Attributes[restAnnotation.ParamFormat]
	}
	return ""
}

func IsRestOperationJSON(o model.Operation) bool {
	return GetRestOperationFormat(o) == "JSON"
}

func IsRestOperationHTML(o model.Operation) bool {
	return GetRestOperationFormat(o) == "HTML"
}

func IsRestOperationCSV(o model.Operation) bool {
	return GetRestOperationFormat(o) == "CSV"
}

func IsRestOperationTXT(o model.Operation) bool {
	return GetRestOperationFormat(o) == "TXT"
}

func IsRestOperationMD(o model.Operation) bool {
	return GetRestOperationFormat(o) == "MD"
}

func IsRestOperationNoContent(o model.Operation) bool {
	return GetRestOperationFormat(o) == "no_content"
}

func IsRestOperationCustom(o model.Operation) bool {
	return GetRestOperationFormat(o) == "custom"
}

func IsRestOperationGenerated(o model.Operation) bool {
	return IsRestOperationJSON(o) || IsRestOperationHTML(o) || IsRestOperationCSV(o) || IsRestOperationTXT(o) || IsRestOperationMD(o) || IsRestOperationNoContent(o) || IsRestOperationCustom(o)
}

func HasContentType(operation model.Operation) bool {
	return GetContentType(operation) != ""
}

func GetContentType(operation model.Operation) string {
	switch GetRestOperationFormat(operation) {
	case "JSON":
		return "application/json"
	case "HTML":
		return "text/html; charset=UTF-8"
	case "CSV":
		return "text/csv; charset=UTF-8"
	case "TXT":
		return "text/plain; charset=UTF-8"
	case "MD":
		return "text/markdown; charset=UTF-8"
	default:
		return ""
	}
}

func GetRestOperationFilename(o model.Operation) string {
	ann, ok := annotation.ResolveAnnotationByName(o.DocLines, restAnnotation.TypeRestOperation)
	if ok {
		return ann.Attributes[restAnnotation.ParamFilename]
	}
	return ""
}

func GetRestOperationRolesString(o model.Operation) string {
	roles := GetRestOperationRoles(o)
	for i, r := range roles {
		roles[i] = fmt.Sprintf("\"%s\"", r)
	}
	return fmt.Sprintf("[]string{%s}", strings.Join(roles, ","))
}

func GetRestOperationRoles(o model.Operation) []string {
	ann, ok := annotation.ResolveAnnotationByName(o.DocLines, restAnnotation.TypeRestOperation)
	if ok {
		rolesAttr, ok := ann.Attributes[restAnnotation.ParamRoles]
		if ok {
			roles := strings.Split(rolesAttr, ",")
			for i, r := range roles {
				roles[i] = strings.Trim(r, " ")
			}
			return roles
		}
	}
	return []string{}
}

func HasInput(o model.Operation) bool {
	if GetRestOperationMethod(o) == "POST" || GetRestOperationMethod(o) == "PUT" {
		for _, arg := range o.InputArgs {
			if arg.TypeName != "int" && arg.TypeName != "string" && arg.TypeName != "context.Context" && arg.TypeName != "rest.Credentials" {
				return true
			}
		}
	}
	return false
}

func HasCredentials(o model.Operation) bool {
	for _, arg := range o.InputArgs {
		if arg.Name == "credentials" {
			return true
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

func ReturnsError(o model.Operation) bool {
	for _, arg := range o.OutputArgs {
		if arg.TypeName == "error" {
			return true
		}
	}
	return false
}

func NeedsContext(o model.Operation) bool {
	return HasContext(o) || ReturnsError(o)
}

func GetContextName(o model.Operation) string {
	for _, arg := range o.InputArgs {
		if arg.TypeName == "context.Context" {
			return arg.Name
		}
	}
	if ReturnsError(o) {
		return "c"
	}
	return ""
}

func GetInputArgType(o model.Operation) string {
	for _, arg := range o.InputArgs {
		if arg.TypeName != "int" && arg.TypeName != "string" && arg.TypeName != "context.Context" && arg.TypeName != "rest.Credentials" {
			return arg.TypeName
		}
	}
	return ""
}

func UsesQueryParams(o model.Operation) bool {
	if GetRestOperationMethod(o) == "GET" {
		count := 0
		for _, arg := range o.InputArgs {
			if arg.TypeName != "context.Context" && arg.TypeName != "rest.Credentials" {
				count++
			}
		}
		return count > 1
	}
	return false
}

func GetInputArgName(o model.Operation) string {
	for _, arg := range o.InputArgs {
		if arg.TypeName != "int" && arg.TypeName != "string" && arg.TypeName != "context.Context" && arg.TypeName != "rest.Credentials" {
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

func findArgInArray(array []string, toMatch string) bool {
	for _, p := range array {
		if strings.Trim(p, " ") == toMatch {
			return true
		}
	}
	return false
}

func RequiresParamValidation(o model.Operation) bool {
	for _, field := range o.InputArgs {
		if field.TypeName == "int" || field.TypeName == "string" && IsInputArgMandatory(o, field) {
			return true
		}
	}
	return false
}

func IsInputArgMandatory(o model.Operation, arg model.Field) bool {
	ann, ok := annotation.ResolveAnnotationByName(o.DocLines, restAnnotation.TypeRestOperation)
	if !ok {
		return false
	}
	optionalArgsString, ok := ann.Attributes[restAnnotation.ParamOptional]
	if !ok {
		return true
	}

	return !findArgInArray(strings.Split(optionalArgsString, ","), arg.Name)
}

func HasUpload(o model.Operation) bool {
	for _, f := range o.InputArgs {
		if IsUploadArg(f) {
			return true
		}
	}
	return false
}

func IsUploadArg(arg model.Field) bool {
	return arg.Name == "upload"
}

func IsPrimitiveArg(f model.Field) bool {
	return f.TypeName == "int" || f.TypeName == "string"
}

func IsNumberArg(f model.Field) bool {
	return f.TypeName == "int"
}

func toFirstUpper(in string) string {
	a := []rune(in)
	a[0] = unicode.ToUpper(a[0])
	return string(a)
}

var handlersTemplate string = `
// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}

import "golang.org/x/net/context"

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
			{{if IsRestOperationGenerated . }}
				subRouter.HandleFunc(  "{{GetRestOperationPath . }}", {{.Name}}(ts)).Methods("{{GetRestOperationMethod . }}")
			{{else}}
				subRouter.HandleFunc(  "{{GetRestOperationPath . }}", ts.{{.Name}}()).Methods("{{GetRestOperationMethod . }}")
			{{end}}
		{{end}}
	{{end}}

	return router
}

{{range $idxOper, $oper := .Operations}}

{{if IsRestOperation $oper}}
{{if IsRestOperationGenerated . }}
// {{$oper.Name}} does the http handling for business logic method service.{{$oper.Name}}
func {{$oper.Name}}( service *{{$structName}} ) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		{{if NeedsContext $oper }}
			{{GetContextName $oper}} := context.WithValue(ctx.New.CreateContext(r), "sessionUid", r.Header.Get("X-session-uid"))
		{{end}}

		{{if HasCredentials $oper}}
			language := "nl"
			langCookie, err := r.Cookie("lang")
			if err == nil {
				language = langCookie.Value
			}
			credentials := rest.ExtractCredentials(language, r)
			err = validateCredentials(credentials, "{{GetRestOperationPath . }}", {{GetRestOperationRolesString $oper}})
			if err != nil {
				errorhandling.HandleHttpError(c, err, w)
				return
			}
		{{end}}

		{{if UsesQueryParams $oper }} {{else}}
		pathParams := mux.Vars(r)
			if len(pathParams) > 0 {
				log.Printf("pathParams:%+v", pathParams)
			}
		{{end}}

		{{if RequiresParamValidation .}}
		// extract url-params
	    validationErrors := []errorh.FieldError{}
	    {{end}}
		{{range .InputArgs}}
			{{if IsPrimitiveArg . }}
				{{if IsNumberArg . }}
					{{.Name}} := 0
					{{if IsRestOperationForm $oper }}
						{{.Name}}String := r.FormValue("{{.Name}}")
						if {{.Name}}String == "" {
					{{else if UsesQueryParams $oper }}
						{{.Name}}String := r.URL.Query().Get("{{.Name}}")
						if {{.Name}}String == "" {
					{{else}}
						{{.Name}}String, exists := pathParams["{{.Name}}"]
						if !exists {
					{{end}}
					{{if IsInputArgMandatory $oper .}}
						validationErrors = append(validationErrors, errorh.FieldErrorForMissingParameter("{{.Name}}"))
					{{else}}
						// optional parameter
					{{end}}
					} else {
						{{.Name}}, err = strconv.Atoi({{.Name}}String)
						if err != nil {
							validationErrors = append(validationErrors, errorh.FieldErrorForInvalidParameter("{{.Name}}"))
						}
					 }
				{{else}}
					{{if IsRestOperationForm $oper }}
						{{.Name}} := r.FormValue("{{.Name}}")
						if {{.Name}} == "" {
					{{else if UsesQueryParams $oper }}
						{{.Name}} := r.URL.Query().Get("{{.Name}}")
						if {{.Name}} == "" {
					{{else}}
						{{.Name}}, exists := pathParams["{{.Name}}"]
						if !exists {
					{{end}}
						{{if IsInputArgMandatory $oper .}}
							validationErrors = append(validationErrors, errorh.FieldErrorForMissingParameter("{{.Name}}"))
					  	{{else}}
					  		// optional parameter
						 {{end}}
						}
					{{end}}
				{{end}}

		{{end}}

		{{if RequiresParamValidation .}}
        if len(validationErrors) > 0 {
            errorhandling.HandleHttpError(c, errorh.NewInvalidInputErrorSpecific(0, validationErrors), w)
            return
        }
        {{end}}

		{{if HasUpload . }}
			{{GetInputArgName . }}, err := {{$oper.Name}}GetUpload({{GetContextName $oper }}, r)
			if err != nil {
				errorhandling.HandleHttpError(c, err, w)
				return
			}
		{{else if HasInput . }}
			// read and parse request body
			var {{GetInputArgName . }} {{GetInputArgType . }}
			err = json.NewDecoder(r.Body).Decode( &{{GetInputArgName . }} )
			if err != nil {
				errorhandling.HandleHttpError(c, errorh.NewInvalidInputErrorf(1, "Error parsing request body: %s", err), w)
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
			errorhandling.HandleHttpError(c, err, w)
			return
		}

		// write OK response body
		{{if HasContentType .}}
			w.Header().Set("Content-Type", "{{GetContentType .}}")
		{{end}}
		{{if IsRestOperationJSON .}}
			{{if HasOutput . }}
				err = json.NewEncoder(w).Encode(result)
				if err != nil {
					log.Printf("Error encoding response payload %+v", err)
				}
			{{end}}
		{{else if IsRestOperationHTML .}}
			{{if HasOutput . }}err = {{$oper.Name}}WriteHTML(w, result){{else}}err = {{$oper.Name}}WriteHTML(w){{end}}
			if err != nil {
				log.Printf("Error encoding response payload %+v", err)
			}
		{{else if IsRestOperationCSV .}}
			w.Header().Set("Content-Disposition", "attachment;filename={{ GetRestOperationFilename .}}")
			{{if HasOutput . }}err = {{$oper.Name}}WriteCSV(w, result){{else}}err = {{$oper.Name}}WriteCSV(w){{end}}
			if err != nil {
				log.Printf("Error encoding response payload %+v", err)
			}
		{{else if IsRestOperationTXT .}}
			_, err = fmt.Fprint(w, result)
			if err != nil {
				log.Printf("Error encoding response payload %+v", err)
			}
		{{else if IsRestOperationMD .}}
			_, err = fmt.Fprint(w, result)
			if err != nil {
				log.Printf("Error encoding response payload %+v", err)
			}
		{{else if IsRestOperationNoContent .}}
			w.WriteHeader(http.StatusNoContent)
		{{else if IsRestOperationCustom .}}
			{{$oper.Name}}HandleResult({{GetContextName $oper }}, w, r, result)
		{{else}}
			errorh.NewInternalErrorf(0, "Not implemented")
		{{end}}
      }
 }
{{end}}
{{end}}
{{end}}

`

var helpersTemplate string = `
// +build !appengine

// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}

import (
    "golang.org/x/net/context"
)

{{ $structName := .Name }}


var logFp *os.File


func openfile( filename string) *os.File {
	fp, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Error opening rest-dump-file %s: %s", filename, err.Error())
	}
	return fp
}

func TestMain(m *testing.M) {

	dirname := "{{.PackageName}}TestLog"
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
    	os.Mkdir(dirname, os.ModePerm)
	}
	logFp = openfile(dirname + "/$testResults.go")
	defer func() {
		logFp.Close()
	}()
	fmt.Fprintf(logFp, "package %s\n\n", dirname )
	fmt.Fprintf(logFp, "// Generated automatically based on running of api-tests\n\n" )
	fmt.Fprintf(logFp, "import (\n")
	fmt.Fprintf(logFp, "\"github.com/MarcGrol/golangAnnotations/generator/rest/testcase\"\n")
	fmt.Fprintf(logFp, ")\n")

	fmt.Fprintf(logFp, "var TestResults = testcase.TestSuiteDescriptor {\n" )
	fmt.Fprintf(logFp, "\tTestCases: []testcase.TestCaseDescriptor{\n")

	beforeAll()

	code := m.Run()

    afterAll()

	fmt.Fprintf(logFp, "},\n" )
	fmt.Fprintf(logFp, "}\n" )

	os.Exit(code)
}

func beforeAll() {
	mytime.SetMockNow()
}

func afterAll() {
    mytime.SetDefaultNow()
}


func testCase(name string, description string) {
	fmt.Fprintf(logFp, "\t\ttestcase.TestCaseDescriptor{\n")
	fmt.Fprintf(logFp, "\t\tName:\"%s\",\n", name)
	fmt.Fprintf(logFp, "\t\tDescription:\"%s\",\n", description)
}

func testCaseDone() {
	fmt.Fprintf(logFp, "},\n")
}


{{range .Operations}}

{{if IsRestOperation . }}
func {{.Name}}TestHelper(url string {{if HasInput . }}, input {{GetInputArgType . }} {{end}} )  ({{if IsRestOperationJSON . }}int {{if HasOutput . }},{{GetOutputArgType . }}{{end}},*errorh.Error{{else}}*httptest.ResponseRecorder{{end}},error) {
	return {{.Name}}TestHelperWithHeaders( url {{if HasInput . }}, input {{end}}, map[string]string{} )
}

func {{.Name}}TestHelperWithHeaders(url string {{if HasInput . }}, input {{GetInputArgType . }} {{end}}, headers map[string]string)  ({{if IsRestOperationJSON . }}int {{if HasOutput . }},{{GetOutputArgType . }}{{end}},*errorh.Error{{else}}*httptest.ResponseRecorder{{end}},error) {

	fmt.Fprintf(logFp, "\t\tOperation:\"%s\",\n", "{{.Name}}")
	defer func() {
		fmt.Fprintf(logFp, "\t},\n")
	}()

	recorder := httptest.NewRecorder()

	{{if HasUpload . }}
		{{.Name}}SetUpload(input)
		req, err := http.NewRequest("{{GetRestOperationMethod . }}", url, nil)
	{{else if HasInput . }}
		rb, _ := json.Marshal(input)
		// indent for readability
		var requestBody bytes.Buffer
		json.Indent(&requestBody, rb, "", "\t")
		req, err := http.NewRequest("{{GetRestOperationMethod . }}", url, strings.NewReader(requestBody.String()))
	{{else}}
		req, err := http.NewRequest("{{GetRestOperationMethod . }}", url, nil)
	{{end}}
	if err != nil {
		{{if IsRestOperationJSON . }}
			{{if HasOutput . }} return 0, nil, nil, err{{else}}return 0, nil, err{{end}}
		{{else}}return nil, err{{end}}
	}
	req.RequestURI = url
	{{if HasUpload . }}
	{{else if HasInput . }}
		req.Header.Set("Content-type", "application/json")
	{{end}}
	{{if HasOutput . }}
		req.Header.Set("Accept", "application/json")
	{{end}}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	headersToBeSorted := []string{}
	for key, values := range req.Header {
		for _, value := range values {
			headersToBeSorted = append(headersToBeSorted, fmt.Sprintf("%s:%s", key, value))
		}
	}
	sort.Strings(headersToBeSorted)

	fmt.Fprintf(logFp, "\tRequest: testcase.RequestDescriptor{\n")
	fmt.Fprintf(logFp, "\tMethod:\"%s\",\n", "{{GetRestOperationMethod . }}")
	fmt.Fprintf(logFp, "\tUrl:\"%s\",\n", url)
	fmt.Fprintf(logFp, "\tHeaders: []string{\n")
	for _, h := range headersToBeSorted {
		fmt.Fprintf(logFp, "\"%s\",\n", h)
	}
	fmt.Fprintf(logFp, "\t},\n")

	{{if HasUpload . }}
	{{else if HasInput . }}
		fmt.Fprintf(logFp, "\tBody:\n" )
		fmt.Fprintf(logFp, "{{BackTick}}%s{{BackTick}}", requestBody.String() )
		fmt.Fprintf(logFp, ",\n" )
	{{end}}
	fmt.Fprintf(logFp, "},\n")

	// dump readable request
	//payload, err := httputil.DumpRequest(req, true)

	fmt.Fprintf(logFp, "\tResponse:testcase.ResponseDescriptor{\n")
	defer func() {
		fmt.Fprintf(logFp, "\t},\n")
	}()

	webservice := {{$structName}}{}
	webservice.HTTPHandler().ServeHTTP(recorder, req)

	{{if IsRestOperationJSON . }}
		// dump readable response
		var responseBody bytes.Buffer
		json.Indent(&responseBody, recorder.Body.Bytes(), "", "\t")
	{{end}}

	fmt.Fprintf(logFp, "\tStatus:%d,\n", recorder.Code)

	headersToBeSorted = []string{}
	for key, values := range recorder.Header() {
		for _, value := range values {
			headersToBeSorted = append(headersToBeSorted, fmt.Sprintf("%s:%s", key, value))
		}
	}
	sort.Strings(headersToBeSorted)

	fmt.Fprintf(logFp, "\tHeaders:[]string{\n")
	for _, h := range headersToBeSorted {
		fmt.Fprintf(logFp, "\"%s\",\n", h)
	}
	fmt.Fprintf(logFp, "\t},\n")
	fmt.Fprintf(logFp, "\tBody:\n{{BackTick}}%s{{BackTick}},\n", {{if IsRestOperationJSON . }}responseBody.String(){{else}}recorder.Body.Bytes(){{end}})

	{{if IsRestOperationJSON . }}
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
	{{else}}
		return recorder, nil
	{{end}}
}
{{end}}
{{end}}
`

var httpClientTemplate string = `
// +build !appengine

// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}

import "golang.org/x/net/context"

{{ $structName := .Name }}

var debug = false


type HTTPClient struct {
	hostName string
}

func NewHTTPClient(host string) *HTTPClient {
	return &HTTPClient{
		hostName: host,
	}
}

{{range .Operations}}

{{if IsRestOperation . }}
{{if IsRestOperationJSON . }}

// {{ToFirstUpper .Name}} can be used by external clients to interact with the system
func (c *HTTPClient) {{ToFirstUpper .Name}}(ctx context.Context, url string {{if HasInput . }}, input {{GetInputArgType . }} {{end}}, cookie *http.Cookie, requestUID string, timeout time.Duration)  (int {{if HasOutput . }},{{GetOutputArgType . }}{{end}},*errorh.Error,error) {

	{{if HasInput . }}
		requestBody, _ := json.Marshal(input)
		req, err := http.NewRequest("{{GetRestOperationMethod . }}", c.hostName+url, strings.NewReader(string(requestBody)))
	{{else}}
		req, err := http.NewRequest("{{GetRestOperationMethod . }}", c.hostName+url, nil)
	{{end}}
	if err != nil {
		{{if HasOutput . }} return 0, nil, nil, err
		{{else}} return 0,  nil, err
		{{end}}
	}
	if cookie != nil {
		req.AddCookie(cookie)
	}
	if requestUID != "" {
		req.Header.Set("X-request-uid", requestUID)
	}
	{{if HasInput . }}
		req.Header.Set("Content-type", "application/json")
	{{end}}
	{{if HasOutput . }}
		req.Header.Set("Accept", "application/json")
	{{end}}
	req.Header.Set("X-CSRF-Token", "true")

    if debug {
		dump, err := httputil.DumpRequest(req, true)
		if err == nil {
			logging.New().Debug(ctx, "HTTP request-payload:\n %s", dump)
		}
    }

	cl := http.Client{}
	cl.Timeout = timeout
	res, err := cl.Do(req)
	if err != nil {
	{{if HasOutput . }}
		return -1, nil, nil, err
	{{else}}
		return -1	, nil, nil
	{{end}}
	}
	defer res.Body.Close()

	if debug {
		respDump, err := httputil.DumpResponse(res, true)
		if err == nil {
			logging.New().Debug(ctx,"HTTP response-payload:\n%s", string(respDump))
		}
	}

	{{if HasOutput . }}
		if res.StatusCode >= http.StatusMultipleChoices {
		    // return error response
			var errorResp errorh.Error
			dec := json.NewDecoder(res.Body)
			err = dec.Decode(&errorResp)
			if err != nil {
				return res.StatusCode, nil, nil, err
			}
			return res.StatusCode, nil, &errorResp, nil
		}

		// return success response
		resp := {{GetOutputArgDeclaration . }}
		dec := json.NewDecoder(res.Body)
		err = dec.Decode({{GetOutputArgName . }})
		if err != nil {
			return res.StatusCode, nil, nil, err
		}
		return res.StatusCode, resp, nil, nil

	{{else}}
		return res.StatusCode, nil, nil
	{{end}}
}
{{end}}
{{end}}
{{end}}
`

var testServiceTemplate = `
// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}

import (
	"github.com/gorilla/mux"
	"github.com/MarcGrol/golangAnnotations/generator/rest/testcase"
)

// HTTPTestHandlerWithRouter registers endpoint in existing router
func HTTPTestHandlerWithRouter(router *mux.Router, results testcase.TestSuiteDescriptor) *mux.Router {
	subRouter := router.PathPrefix("{{GetRestServicePath . }}").Subrouter()

	subRouter.HandleFunc("/logs.md", testcase.WriteTestLogsAsMarkdown(results)).Methods("GET")

	return router
}

`
