package rest

import (
	"fmt"
	"log"
	"regexp"
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
		if isRestService(service) {
			{
				target := fmt.Sprintf("%s/$http%s.go", targetDir, toFirstUpper(service.Name))
				err = generationUtil.GenerateFileFromTemplateFile(service, fmt.Sprintf("%s.%s", service.PackageName, toFirstUpper(service.Name)), "http-handlers", "generator/rest/httpHandlers.go.tmpl", customTemplateFuncs, target)
				if err != nil {
					log.Fatalf("Error generating handlers for service %s: %s", service.Name, err)
					return err
				}
			}
			if !IsRestServiceNoTest(service) {
				{
					target := fmt.Sprintf("%s/$http%sHelpers_test.go", targetDir, toFirstUpper(service.Name))
					err = generationUtil.GenerateFileFromTemplateFile(service, fmt.Sprintf("%s.%s", service.PackageName, toFirstUpper(service.Name)), "test-helpers", "generator/rest/testHelpers.go.tmpl", customTemplateFuncs, target)
					if err != nil {
						log.Fatalf("Error generating helpers for service %s: %s", service.Name, err)
						return err
					}
				}
				{
					target := fmt.Sprintf("%s/$httpTest%s.go", targetDir, toFirstUpper(service.Name))
					err = generationUtil.GenerateFileFromTemplateFile(service, fmt.Sprintf("%s.%s", service.PackageName, toFirstUpper(service.Name)), "testService", "generator/rest/testService.go.tmpl", customTemplateFuncs, target)
					if err != nil {
						log.Fatalf("Error generating testHandler for service %s: %s", service.Name, err)
						return err
					}
				}
				{
					target := fmt.Sprintf("%s/$httpClientFor%s.go", targetDir, toFirstUpper(service.Name))
					err = generationUtil.GenerateFileFromTemplateFile(service, fmt.Sprintf("%s.%s", service.PackageName, toFirstUpper(service.Name)), "http-client", "generator/rest/httpClient.go.tmpl", customTemplateFuncs, target)
					if err != nil {
						log.Fatalf("Error generating httpClient for service %s: %s", service.Name, err)
						return err
					}
				}
			}
		}
	}
	return nil
}

var customTemplateFuncs = template.FuncMap{
	"IsRestService":                         isRestService,
	"ExtractImports":                        extractImports,
	"GetRestServicePath":                    getRestServicePath,
	"GetExtractCredentialsMethod":           getExtractCredentialsMethod,
	"IsRestServiceNoValidation":             isRestServiceNoValidation,
	"IsRestOperation":                       isRestOperation,
	"IsRestOperationNoWrap":                 isRestOperationNoWrap,
	"IsRestOperationGenerated":              isRestOperationGenerated,
	"HasRestOperationAfter":                 hasRestOperationAfter,
	"GetRestOperationPath":                  getRestOperationPath,
	"GetRestOperationMethod":                getRestOperationMethod,
	"IsRestOperationForm":                   isRestOperationForm,
	"IsRestOperationJSON":                   isRestOperationJSON,
	"IsRestOperationHTML":                   isRestOperationHTML,
	"IsRestOperationCSV":                    isRestOperationCSV,
	"IsRestOperationTXT":                    isRestOperationTXT,
	"IsRestOperationMD":                     isRestOperationMD,
	"IsRestOperationNoContent":              isRestOperationNoContent,
	"IsRestOperationCustom":                 isRestOperationCustom,
	"HasContentType":                        hasContentType,
	"GetContentType":                        getContentType,
	"GetRestOperationFilename":              getRestOperationFilename,
	"GetRestOperationRolesString":           getRestOperationRolesString,
	"GetRestOperationProducesEvents":        getRestOperationProducesEvents,
	"GetRestOperationProducesEventsAsSlice": getRestOperationProducesEventsAsSlice,
	"HasOperationsWithInput":                hasOperationsWithInput,
	"HasInput":                              hasInput,
	"GetInputArgType":                       getInputArgType,
	"GetOutputArgDeclaration":               getOutputArgDeclaration,
	"GetOutputArgName":                      getOutputArgName,
	"HasAnyPathParam":                       hasAnyPathParam,
	"IsSliceParam":                          isSliceParam,
	"IsQueryParam":                          isQueryParam,
	"GetInputArgName":                       getInputArgName,
	"GetInputParamString":                   getInputParamString,
	"GetOutputArgType":                      getOutputArgType,
	"HasOutput":                             hasOutput,
	"HasMetaOutput":                         hasMetaOutput,
	"IsPrimitiveArg":                        isPrimitiveArg,
	"IsNumberArg":                           isNumberArg,
	"RequiresParamValidation":               requiresParamValidation,
	"IsInputArgMandatory":                   isInputArgMandatory,
	"HasUpload":                             hasUpload,
	"IsUploadArg":                           isUploadArg,
	"HasCredentials":                        hasCredentials,
	"HasContext":                            hasContext,
	"ReturnsError":                          returnsError,
	"NeedsContext":                          needsContext,
	"GetContextName":                        getContextName,
	"WithBackTicks":                         surroundWithBackTicks,
	"BackTick":                              backTick,
	"ToFirstUpper":                          toFirstUpper,
}

func backTick() string {
	return "`"
}

func surroundWithBackTicks(body string) string {
	return fmt.Sprintf("`%s'", body)
}

func isRestService(s model.Struct) bool {
	_, ok := annotation.ResolveAnnotationByName(s.DocLines, restAnnotation.TypeRestService)
	return ok
}

func getRestServicePath(s model.Struct) string {
	if ann, ok := annotation.ResolveAnnotationByName(s.DocLines, restAnnotation.TypeRestService); ok {
		return ann.Attributes[restAnnotation.ParamPath]
	}
	return ""
}

func getExtractCredentialsMethod(s model.Struct) string {
	if ann, ok := annotation.ResolveAnnotationByName(s.DocLines, restAnnotation.TypeRestService); ok {
		switch ann.Attributes[restAnnotation.ParamCredentials] {
		case "all":
			return "rest.ExtractAllCredentials"
		case "admin":
			return "rest.ExtractAdminCredentials"
		case "none":
			return "rest.ExtractNoCredentials"
		}
	}
	return "extractCredentials"
}

func isRestServiceNoValidation(s model.Struct) bool {
	if ann, ok := annotation.ResolveAnnotationByName(s.DocLines, restAnnotation.TypeRestService); ok {
		return ann.Attributes[restAnnotation.ParamNoValidation] == "true"
	}
	return false
}

func IsRestServiceNoTest(s model.Struct) bool {
	if ann, ok := annotation.ResolveAnnotationByName(s.DocLines, restAnnotation.TypeRestService); ok {
		return ann.Attributes[restAnnotation.ParamNoTest] == "true"
	}
	return false
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

func extractImports(s model.Struct) []string {
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

func hasOperationsWithInput(s model.Struct) bool {
	for _, o := range s.Operations {
		if hasInput(*o) == true {
			return true
		}
	}
	return false
}

func isRestOperation(o model.Operation) bool {
	_, ok := annotation.ResolveAnnotationByName(o.DocLines, restAnnotation.TypeRestOperation)
	return ok
}

func isRestOperationNoWrap(o model.Operation) bool {
	if ann, ok := annotation.ResolveAnnotationByName(o.DocLines, restAnnotation.TypeRestOperation); ok {
		return ann.Attributes[restAnnotation.ParamNoWrap] == "true"
	}
	return false
}

func isRestOperationGenerated(o model.Operation) bool {
	return !isRestOperationNoWrap(o)
}

func hasRestOperationAfter(o model.Operation) bool {
	if ann, ok := annotation.ResolveAnnotationByName(o.DocLines, restAnnotation.TypeRestOperation); ok {
		return ann.Attributes[restAnnotation.ParamAfter] == "true"
	}
	return false
}

func getRestOperationPath(o model.Operation) string {
	if ann, ok := annotation.ResolveAnnotationByName(o.DocLines, restAnnotation.TypeRestOperation); ok {
		return ann.Attributes[restAnnotation.ParamPath]
	}
	return ""
}

func hasAnyPathParam(o model.Operation) bool {
	return len(GetAllPathParams(o)) > 0
}

func GetAllPathParams(o model.Operation) []string {
	re, _ := regexp.Compile(`\{\w+\}`)
	path := getRestOperationPath(o)
	params := re.FindAllString(path, -1)
	for idx, param := range params {
		params[idx] = param[1 : len(param)-1]
	}
	return params
}

func getRestOperationMethod(o model.Operation) string {
	if ann, ok := annotation.ResolveAnnotationByName(o.DocLines, restAnnotation.TypeRestOperation); ok {
		return ann.Attributes[restAnnotation.ParamMethod]
	}
	return ""
}

func isRestOperationForm(o model.Operation) bool {
	if ann, ok := annotation.ResolveAnnotationByName(o.DocLines, restAnnotation.TypeRestOperation); ok {
		return ann.Attributes[restAnnotation.ParamForm] == "true"
	}
	return false
}

func GetRestOperationFormat(o model.Operation) string {
	if ann, ok := annotation.ResolveAnnotationByName(o.DocLines, restAnnotation.TypeRestOperation); ok {
		return ann.Attributes[restAnnotation.ParamFormat]
	}
	return ""
}

func isRestOperationJSON(o model.Operation) bool {
	return GetRestOperationFormat(o) == "JSON"
}

func isRestOperationHTML(o model.Operation) bool {
	return GetRestOperationFormat(o) == "HTML"
}

func isRestOperationCSV(o model.Operation) bool {
	return GetRestOperationFormat(o) == "CSV"
}

func isRestOperationTXT(o model.Operation) bool {
	return GetRestOperationFormat(o) == "TXT"
}

func isRestOperationMD(o model.Operation) bool {
	return GetRestOperationFormat(o) == "MD"
}

func isRestOperationNoContent(o model.Operation) bool {
	return GetRestOperationFormat(o) == "no_content"
}

func isRestOperationCustom(o model.Operation) bool {
	return GetRestOperationFormat(o) == "custom"
}

func hasContentType(operation model.Operation) bool {
	return getContentType(operation) != ""
}

func getContentType(operation model.Operation) string {
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

func getRestOperationFilename(o model.Operation) string {
	if ann, ok := annotation.ResolveAnnotationByName(o.DocLines, restAnnotation.TypeRestOperation); ok {
		return ann.Attributes[restAnnotation.ParamFilename]
	}
	return ""
}

func getRestOperationRolesString(o model.Operation) string {
	roles := GetRestOperationRoles(o)
	for i, r := range roles {
		roles[i] = fmt.Sprintf("\"%s\"", r)
	}
	return fmt.Sprintf("[]string{%s}", strings.Join(roles, ","))
}

func GetRestOperationRoles(o model.Operation) []string {
	if ann, ok := annotation.ResolveAnnotationByName(o.DocLines, restAnnotation.TypeRestOperation); ok {
		if rolesAttr, ok := ann.Attributes[restAnnotation.ParamRoles]; ok {
			roles := strings.Split(rolesAttr, ",")
			for i, r := range roles {
				roles[i] = strings.Trim(r, " ")
			}
			return roles
		}
	}
	return []string{}
}

func getRestOperationProducesEvents(o model.Operation) string {
	return asStringSlice(getRestOperationProducesEventsAsSlice(o))
}

func getRestOperationProducesEventsAsSlice(o model.Operation) []string {
	if ann, ok := annotation.ResolveAnnotationByName(o.DocLines, restAnnotation.TypeRestOperation); ok {
		if attrs, ok := ann.Attributes[restAnnotation.ParamProducesEvents]; ok {
			eventsProduced := []string{}
			for _, evt := range strings.Split(attrs, ",") {
				evt := strings.TrimSpace(evt)
				if evt != "" {
					eventsProduced = append(eventsProduced, evt)
				}
			}
			return eventsProduced
		}
	}
	return []string{}
}

func asStringSlice(in []string) string {
	adjusted := []string{}
	for _, i := range in {
		adjusted = append(adjusted, fmt.Sprintf("\"%s\"", i))
	}
	return fmt.Sprintf("[]string{%s}", strings.Join(adjusted, ","))
}

func hasInput(o model.Operation) bool {
	if getRestOperationMethod(o) == "POST" || getRestOperationMethod(o) == "PUT" {
		for _, arg := range o.InputArgs {
			if !isPrimitiveArg(arg) && !IsContextArg(arg) && !IsCredentialsArg(arg) {
				return true
			}
		}
	}
	return false
}

func hasCredentials(o model.Operation) bool {
	for _, arg := range o.InputArgs {
		if IsCredentialsArg(arg) {
			return true
		}
	}
	return false
}

func hasContext(o model.Operation) bool {
	for _, arg := range o.InputArgs {
		if IsContextArg(arg) {
			return true
		}
	}
	return false
}

func returnsError(o model.Operation) bool {
	for _, arg := range o.OutputArgs {
		if IsErrorArg(arg) {
			return true
		}
	}
	return false
}

func needsContext(o model.Operation) bool {
	return hasContext(o) || returnsError(o)
}

func getContextName(o model.Operation) string {
	for _, arg := range o.InputArgs {
		if IsContextArg(arg) {
			return arg.Name
		}
	}
	if returnsError(o) {
		return "c"
	}
	return ""
}

func getInputArgType(o model.Operation) string {
	for _, arg := range o.InputArgs {
		if !isPrimitiveArg(arg) && !IsContextArg(arg) && !IsCredentialsArg(arg) {
			return arg.TypeName
		}
	}
	return ""
}

func isSliceParam(arg model.Field) bool {
	return arg.IsSlice
}

func isQueryParam(o model.Operation, arg model.Field) bool {
	if IsContextArg(arg) || IsCredentialsArg(arg) {
		return false
	}
	for _, pathParam := range GetAllPathParams(o) {
		if pathParam == arg.Name {
			return false
		}
	}
	return true
}

func getInputArgName(o model.Operation) string {
	for _, arg := range o.InputArgs {
		if !isPrimitiveArg(arg) && !IsContextArg(arg) && !IsCredentialsArg(arg) {
			return arg.Name
		}
	}
	return ""
}

func getInputParamString(o model.Operation) string {
	args := []string{}
	for _, arg := range o.InputArgs {
		args = append(args, arg.Name)
	}
	return strings.Join(args, ",")
}

func hasOutput(o model.Operation) bool {
	for _, arg := range o.OutputArgs {
		if !IsErrorArg(arg) {
			return true
		}
	}
	return false
}

func getOutputArgType(o model.Operation) string {
	for _, arg := range o.OutputArgs {
		if !IsErrorArg(arg) {
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

func hasMetaOutput(o model.Operation) bool {
	var count = 0
	for _, arg := range o.OutputArgs {
		if !IsErrorArg(arg) {
			count += 1
			if count == 2 {
				return true
			}
		}
	}
	return false
}

func getOutputArgDeclaration(o model.Operation) string {
	for _, arg := range o.OutputArgs {
		if !IsErrorArg(arg) {
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

func getOutputArgName(o model.Operation) string {
	for _, arg := range o.OutputArgs {
		if !IsErrorArg(arg) {
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

func requiresParamValidation(o model.Operation) bool {
	for _, field := range o.InputArgs {
		if isNumberArg(field) || IsStringArg(field) && isInputArgMandatory(o, field) {
			return true
		}
	}
	return false
}

func isInputArgMandatory(o model.Operation, arg model.Field) bool {
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

func hasUpload(o model.Operation) bool {
	for _, f := range o.InputArgs {
		if isUploadArg(f) {
			return true
		}
	}
	return false
}

func IsErrorArg(f model.Field) bool {
	return f.TypeName == "error"
}

func isUploadArg(f model.Field) bool {
	return f.Name == "upload"
}

func IsContextArg(f model.Field) bool {
	return f.TypeName == "context.Context"
}

func IsCredentialsArg(f model.Field) bool {
	return f.TypeName == "rest.Credentials"
}

func isPrimitiveArg(f model.Field) bool {
	return isNumberArg(f) || IsStringArg(f)
}

func isNumberArg(f model.Field) bool {
	return f.TypeName == "int"
}

func IsStringArg(f model.Field) bool {
	return f.TypeName == "string"
}

func toFirstUpper(in string) string {
	a := []rune(in)
	a[0] = unicode.ToUpper(a[0])
	return string(a)
}
