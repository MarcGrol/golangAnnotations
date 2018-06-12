package rest

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"text/template"
	"unicode"

	"github.com/MarcGrol/golangAnnotations/generator"
	"github.com/MarcGrol/golangAnnotations/generator/annotation"
	"github.com/MarcGrol/golangAnnotations/generator/generationUtil"
	"github.com/MarcGrol/golangAnnotations/generator/rest/restAnnotation"
	"github.com/MarcGrol/golangAnnotations/model"
)

type Generator struct {
}

func NewGenerator() generator.Generator {
	return &Generator{}
}

func (eg *Generator) GetAnnotations() []annotation.AnnotationDescriptor {
	return restAnnotation.Get()
}

func (eg *Generator) Generate(inputDir string, parsedSource model.ParsedSources) error {
	return generate(inputDir, parsedSource.Structs)
}

type generateContext struct {
	targetDir   string
	packageName string
	service     model.Struct
}

func generate(inputDir string, structs []model.Struct) error {

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
			ctx := generateContext{
				targetDir:   targetDir,
				packageName: packageName,
				service:     service,
			}
			err = generateHttpService(ctx)
			if err != nil {
				return err
			}

			if !IsRestServiceNoTest(service) {
				err = generateHttpTestHelpers(ctx)
				if err != nil {
					return err
				}
				err = generateHttpTestService(ctx)
				if err != nil {
					return err
				}
				err = generateHttpClient(ctx)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func generateHttpService(ctx generateContext) error {
	err := generationUtil.Generate(generationUtil.Info{
		Src:            fmt.Sprintf("%s.%s", ctx.service.PackageName, ToFirstUpper(ctx.service.Name)),
		TargetFilename: generationUtil.Prefixed(fmt.Sprintf("%s/http%s.go", ctx.targetDir, ToFirstUpper(ctx.service.Name))),
		TemplateName:   "http-handlers",
		TemplateString: httpHandlersTemplate,
		FuncMap:        customTemplateFuncs,
		Data:           ctx.service,
	})
	if err != nil {
		log.Fatalf("Error generating handlers for service %s: %s", ctx.service.Name, err)
		return err
	}
	return nil
}

func generateHttpTestHelpers(ctx generateContext) error {
	err := generationUtil.Generate(generationUtil.Info{
		Src:            fmt.Sprintf("%s.%s", ctx.service.PackageName, ToFirstUpper(ctx.service.Name)),
		TargetFilename: generationUtil.Prefixed(fmt.Sprintf("%s/http%sHelpers_test.go", ctx.targetDir, ToFirstUpper(ctx.service.Name))),
		TemplateName:   "http-test-helpers",
		TemplateString: testHelpersTemplate,
		FuncMap:        customTemplateFuncs,
		Data:           ctx.service,
	})
	if err != nil {
		log.Fatalf("Error generating helpers for service %s: %s", ctx.service.Name, err)
		return err
	}
	return nil
}

func generateHttpTestService(ctx generateContext) error {
	// create this file within a subdirectoty
	ctx.packageName = ctx.packageName + "TestLog"

	ctx.service.PackageName = ctx.packageName
	target := generationUtil.Prefixed(fmt.Sprintf("%s/%s/httpTest%s.go", ctx.targetDir, ctx.packageName, ToFirstUpper(ctx.service.Name)))
	err := generationUtil.Generate(generationUtil.Info{
		Src:            fmt.Sprintf("%s.%s", ctx.service.PackageName, ToFirstUpper(ctx.service.Name)),
		TargetFilename: target,
		TemplateName:   "testService",
		TemplateString: testServiceTemplate,
		FuncMap:        customTemplateFuncs,
		Data:           ctx.service,
	})
	if err != nil {
		log.Fatalf("Error generating testHandler for service %s: %s", ctx.service.Name, err)
		return err
	}
	return nil
}

func generateHttpClient(ctx generateContext) error {
	err := generationUtil.Generate(generationUtil.Info{
		Src:            fmt.Sprintf("%s.%s", ctx.service.PackageName, ToFirstUpper(ctx.service.Name)),
		TargetFilename: generationUtil.Prefixed(fmt.Sprintf("%s/httpClientFor%s.go", ctx.targetDir, ToFirstUpper(ctx.service.Name))),
		TemplateName:   "http-client",
		TemplateString: httpClientTemplate,
		FuncMap:        customTemplateFuncs,
		Data:           ctx.service,
	})

	if err != nil {
		log.Fatalf("Error generating httpClient for service %s: %s", ctx.service.Name, err)
		return err
	}
	return nil
}

var customTemplateFuncs = template.FuncMap{
	"IsRestService":                         IsRestService,
	"ExtractImports":                        ExtractImports,
	"GetRestServicePath":                    GetRestServicePath,
	"GetExtractRequestContextMethod":        GetExtractRequestContextMethod,
	"IsRestServiceNoValidation":             IsRestServiceNoValidation,
	"IsRestOperation":                       IsRestOperation,
	"IsRestOperationNoWrap":                 IsRestOperationNoWrap,
	"IsRestOperationGenerated":              IsRestOperationGenerated,
	"HasRestOperationAfter":                 HasRestOperationAfter,
	"GetRestOperationPath":                  GetRestOperationPath,
	"GetRestOperationMethod":                GetRestOperationMethod,
	"IsRestOperationTransactional":          IsRestOperationTransactional,
	"IsRestOperationForm":                   IsRestOperationForm,
	"IsRestOperationJSON":                   IsRestOperationJSON,
	"IsRestOperationHTML":                   IsRestOperationHTML,
	"IsRestOperationCSV":                    IsRestOperationCSV,
	"IsRestOperationTXT":                    IsRestOperationTXT,
	"IsRestOperationMD":                     IsRestOperationMD,
	"IsRestOperationNoContent":              IsRestOperationNoContent,
	"IsRestOperationCustom":                 IsRestOperationCustom,
	"HasContentType":                        HasContentType,
	"GetContentType":                        GetContentType,
	"GetRestOperationFilename":              GetRestOperationFilename,
	"GetRestOperationRolesString":           GetRestOperationRolesString,
	"GetRestOperationProducesEvents":        GetRestOperationProducesEvents,
	"GetRestOperationProducesEventsAsSlice": GetRestOperationProducesEventsAsSlice,
	"HasOperationsWithInput":                HasOperationsWithInput,
	"HasInput":                              HasInput,
	"GetInputArgType":                       GetInputArgType,
	"GetOutputArgDeclaration":               GetOutputArgDeclaration,
	"GetOutputArgsDeclaration":              GetOutputArgsDeclaration,
	"GetOutputArgName":                      GetOutputArgName,
	"HasAnyPathParam":                       HasAnyPathParam,
	"IsSliceParam":                          IsSliceParam,
	"IsQueryParam":                          IsQueryParam,
	"GetInputArgName":                       GetInputArgName,
	"GetInputParamString":                   GetInputParamString,
	"GetOutputArgType":                      GetOutputArgType,
	"HasOutput":                             HasOutput,
	"HasMetaOutput":                         HasMetaOutput,
	"IsMetaCallback":                        IsMetaCallback,
	"IsPrimitiveArg":                        IsPrimitiveArg,
	"IsIntArg":                              IsIntArg,
	"IsBoolArg":                             IsBoolArg,
	"IsStringArg":                           IsStringArg,
	"IsStringSliceArg":                      IsStringSliceArg,
	"IsDateArg":                             IsDateArg,
	"RequiresParamValidation":               RequiresParamValidation,
	"IsInputArgMandatory":                   IsInputArgMandatory,
	"HasUpload":                             HasUpload,
	"IsUploadArg":                           IsUploadArg,
	"HasRequestContext":                     HasRequestContext,
	"HasContext":                            HasContext,
	"ReturnsError":                          ReturnsError,
	"NeedsContext":                          NeedsContext,
	"GetContextName":                        GetContextName,
	"WithBackTicks":                         SurroundWithBackTicks,
	"BackTick":                              BackTick,
	"ToFirstUpper":                          ToFirstUpper,
}

func BackTick() string {
	return "`"
}

func SurroundWithBackTicks(body string) string {
	return fmt.Sprintf("`%s'", body)
}

func IsRestService(s model.Struct) bool {
	annotations := annotation.NewRegistry(restAnnotation.Get())
	_, ok := annotations.ResolveAnnotationByName(s.DocLines, restAnnotation.TypeRestService)
	return ok
}

func IsRestOperationTransactional(s model.Struct, o model.Operation) bool {
	annotations := annotation.NewRegistry(restAnnotation.Get())
	if ann, ok := annotations.ResolveAnnotationByName(o.DocLines, restAnnotation.TypeRestOperation); ok {
		return ann.Attributes[restAnnotation.ParamTransactional] == "true"
	}
	return false
}

func IsRestServiceUnprotected(s model.Struct) bool {
	annotations := annotation.NewRegistry(restAnnotation.Get())
	ann, ok := annotations.ResolveAnnotationByName(s.DocLines, restAnnotation.TypeRestService)
	return ok && ann.Attributes[restAnnotation.ParamProtected] != "true"
}

func GetRestServicePath(s model.Struct) string {
	annotations := annotation.NewRegistry(restAnnotation.Get())
	if ann, ok := annotations.ResolveAnnotationByName(s.DocLines, restAnnotation.TypeRestService); ok {
		return ann.Attributes[restAnnotation.ParamPath]
	}
	return ""
}

func GetExtractRequestContextMethod(s model.Struct) string {
	annotations := annotation.NewRegistry(restAnnotation.Get())
	if ann, ok := annotations.ResolveAnnotationByName(s.DocLines, restAnnotation.TypeRestService); ok {
		switch ann.Attributes[restAnnotation.ParamCredentials] {
		case "all":
			return "request.NewContext"
		case "admin":
			return "request.NewAdminContext"
		case "none":
			return "request.NewMinimalContext"
		}
	}
	return "extractRequestContext"
}

func IsRestServiceNoValidation(s model.Struct) bool {
	annotations := annotation.NewRegistry(restAnnotation.Get())
	if ann, ok := annotations.ResolveAnnotationByName(s.DocLines, restAnnotation.TypeRestService); ok {
		return ann.Attributes[restAnnotation.ParamNoValidation] == "true"
	}
	return false
}

func IsRestServiceNoTest(s model.Struct) bool {
	annotations := annotation.NewRegistry(restAnnotation.Get())
	if ann, ok := annotations.ResolveAnnotationByName(s.DocLines, restAnnotation.TypeRestService); ok {
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

func ExtractImports(s model.Struct) []string {
	importsMap := map[string]bool{}
	for _, o := range s.Operations {
		for _, ia := range o.InputArgs {
			if isImportToBeIgnored(ia.PackageName) == false {
				importsMap[ia.PackageName] = true
			}
		}
		for _, oa := range o.OutputArgs {
			if isImportToBeIgnored(oa.PackageName) == false {
				importsMap[oa.PackageName] = true
			}
		}
	}
	return mapToSlice(importsMap)
}

func mapToSlice(importsMap map[string]bool) []string {
	importsList := make([]string, 0)
	for k := range importsMap {
		importsList = append(importsList, k)
	}
	return importsList
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
	annotations := annotation.NewRegistry(restAnnotation.Get())
	_, ok := annotations.ResolveAnnotationByName(o.DocLines, restAnnotation.TypeRestOperation)
	return ok
}

func IsRestOperationNoWrap(o model.Operation) bool {
	annotations := annotation.NewRegistry(restAnnotation.Get())
	if ann, ok := annotations.ResolveAnnotationByName(o.DocLines, restAnnotation.TypeRestOperation); ok {
		return ann.Attributes[restAnnotation.ParamNoWrap] == "true"
	}
	return false
}

func IsRestOperationGenerated(o model.Operation) bool {
	return !IsRestOperationNoWrap(o)
}

func HasRestOperationAfter(o model.Operation) bool {
	annotations := annotation.NewRegistry(restAnnotation.Get())
	if ann, ok := annotations.ResolveAnnotationByName(o.DocLines, restAnnotation.TypeRestOperation); ok {
		return ann.Attributes[restAnnotation.ParamAfter] == "true"
	}
	return false
}

func GetRestOperationPath(o model.Operation) string {
	annotations := annotation.NewRegistry(restAnnotation.Get())
	if ann, ok := annotations.ResolveAnnotationByName(o.DocLines, restAnnotation.TypeRestOperation); ok {
		return ann.Attributes[restAnnotation.ParamPath]
	}
	return ""
}

func HasAnyPathParam(o model.Operation) bool {
	return len(getAllPathParams(o)) > 0
}

func getAllPathParams(o model.Operation) []string {
	re, _ := regexp.Compile(`\{\w+\}`)
	path := GetRestOperationPath(o)
	params := re.FindAllString(path, -1)
	for idx, param := range params {
		params[idx] = param[1 : len(param)-1]
	}
	return params
}

func GetRestOperationMethod(o model.Operation) string {
	annotations := annotation.NewRegistry(restAnnotation.Get())
	if ann, ok := annotations.ResolveAnnotationByName(o.DocLines, restAnnotation.TypeRestOperation); ok {
		return ann.Attributes[restAnnotation.ParamMethod]
	}
	return ""
}

func IsRestOperationForm(o model.Operation) bool {
	annotations := annotation.NewRegistry(restAnnotation.Get())
	if ann, ok := annotations.ResolveAnnotationByName(o.DocLines, restAnnotation.TypeRestOperation); ok {
		return ann.Attributes[restAnnotation.ParamForm] == "true"
	}
	return false
}

func GetRestOperationFormat(o model.Operation) string {
	annotations := annotation.NewRegistry(restAnnotation.Get())
	if ann, ok := annotations.ResolveAnnotationByName(o.DocLines, restAnnotation.TypeRestOperation); ok {
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
	annotations := annotation.NewRegistry(restAnnotation.Get())
	if ann, ok := annotations.ResolveAnnotationByName(o.DocLines, restAnnotation.TypeRestOperation); ok {
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
	annotations := annotation.NewRegistry(restAnnotation.Get())
	if ann, ok := annotations.ResolveAnnotationByName(o.DocLines, restAnnotation.TypeRestOperation); ok {
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

func GetRestOperationProducesEvents(o model.Operation) string {
	return asStringSlice(GetRestOperationProducesEventsAsSlice(o))
}

func GetRestOperationProducesEventsAsSlice(o model.Operation) []string {
	annotations := annotation.NewRegistry(restAnnotation.Get())
	if ann, ok := annotations.ResolveAnnotationByName(o.DocLines, restAnnotation.TypeRestOperation); ok {
		if attrs, ok := ann.Attributes[restAnnotation.ParamProducesEvents]; ok {
			eventsProduced := make([]string, 0)
			for _, e := range strings.Split(attrs, ",") {
				evt := strings.TrimSpace(e)
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
	adjusted := make([]string, 0)
	for _, i := range in {
		adjusted = append(adjusted, fmt.Sprintf("\"%s\"", i))
	}
	return fmt.Sprintf("[]string{%s}", strings.Join(adjusted, ","))
}

func HasInput(o model.Operation) bool {
	if GetRestOperationMethod(o) == "POST" || GetRestOperationMethod(o) == "PUT" {
		for _, arg := range o.InputArgs {
			if !IsPrimitiveArg(arg) && !IsContextArg(arg) && !IsRequestContextArg(arg) {
				return true
			}
		}
	}
	return false
}

func HasRequestContext(o model.Operation) bool {
	for _, arg := range o.InputArgs {
		if IsRequestContextArg(arg) {
			return true
		}
	}
	return false
}

func HasContext(o model.Operation) bool {
	for _, arg := range o.InputArgs {
		if IsContextArg(arg) {
			return true
		}
	}
	return false
}

func ReturnsError(o model.Operation) bool {
	for _, arg := range o.OutputArgs {
		if IsErrorArg(arg) {
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
		if IsContextArg(arg) {
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
		if !IsPrimitiveArg(arg) && !IsContextArg(arg) && !IsRequestContextArg(arg) {
			return arg.TypeName
		}
	}
	return ""
}

func IsSliceParam(arg model.Field) bool {
	return arg.IsSlice()
}

func IsQueryParam(o model.Operation, arg model.Field) bool {
	if IsContextArg(arg) || IsRequestContextArg(arg) {
		return false
	}
	for _, pathParam := range getAllPathParams(o) {
		if pathParam == arg.Name {
			return false
		}
	}
	return true
}

func GetInputArgName(o model.Operation) string {
	for _, arg := range o.InputArgs {
		if !IsPrimitiveArg(arg) && !IsContextArg(arg) && !IsRequestContextArg(arg) {
			return arg.Name
		}
	}
	return ""
}

func GetInputParamString(o model.Operation) string {
	args := make([]string, 0)
	for _, arg := range o.InputArgs {
		args = append(args, arg.Name)
	}
	return strings.Join(args, ",")
}

func HasOutput(o model.Operation) bool {
	for _, arg := range o.OutputArgs {
		if !IsErrorArg(arg) {
			return true
		}
	}
	return false
}

func GetOutputArgType(o model.Operation) string {
	for _, arg := range o.OutputArgs {
		if !IsErrorArg(arg) {
			if arg.IsSlice() {
				return arg.TypeName
			}
			pointer := ""
			if arg.IsPointer {
				pointer = "*"
			}
			return fmt.Sprintf("%s%s", pointer, arg.TypeName)
		}
	}
	return ""
}

func HasMetaOutput(o model.Operation) bool {
	return GetMetaArg(o) != nil
}

func IsMetaCallback(o model.Operation) bool {
	arg := GetMetaArg(o)
	return arg != nil && IsMetaCallbackArg(*arg)
}

func GetMetaArg(o model.Operation) *model.Field {
	var count = 0
	for _, arg := range o.OutputArgs {
		if !IsErrorArg(arg) {
			count++
			if count == 2 {
				return &arg
			}
		}
	}
	return nil
}

func GetOutputArgDeclaration(o model.Operation) string {
	for _, arg := range o.OutputArgs {
		if !IsErrorArg(arg) {
			addressOf := ""
			if arg.IsPointer {
				addressOf = "&"
			}

			if arg.IsSlice() {
				return fmt.Sprintf("%s{}", arg.TypeName)

			}
			return fmt.Sprintf("%s%s{}", addressOf, arg.TypeName)
		}
	}
	return ""
}

func GetOutputArgsDeclaration(o model.Operation) []string {
	args := make([]string, 0)
	for idx, arg := range o.OutputArgs {
		if !IsErrorArg(arg) {
			name := ""
			pointer := ""
			if arg.IsPointer {
				pointer = "*"
			}
			switch idx {
			case 0:
				name = "result"
			case 1:
				name = "meta"
			}

			if arg.IsSlice() {
				args = append(args, fmt.Sprintf("var %s %s", name, arg.TypeName))
			} else {
				args = append(args, fmt.Sprintf("var %s %s%s", name, pointer, arg.TypeName))
			}
		}
	}
	return args
}

func GetOutputArgName(o model.Operation) string {
	for _, arg := range o.OutputArgs {
		if !IsErrorArg(arg) {
			if !arg.IsPointer || arg.IsSlice() {
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
		if (IsIntArg(field) || IsBoolArg(field) || IsStringSliceArg(field) || IsStringArg(field)) && IsInputArgMandatory(o, field) {
			return true
		}
	}
	return false
}

func IsInputArgMandatory(o model.Operation, arg model.Field) bool {
	annotations := annotation.NewRegistry(restAnnotation.Get())
	ann, ok := annotations.ResolveAnnotationByName(o.DocLines, restAnnotation.TypeRestOperation)
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

func IsErrorArg(f model.Field) bool {
	return f.TypeName == "error"
}

func IsUploadArg(f model.Field) bool {
	return f.Name == "upload"
}

func IsContextArg(f model.Field) bool {
	return f.TypeName == "context.Context"
}

func IsRequestContextArg(f model.Field) bool {
	return f.TypeName == "request.Context"
}

func IsMetaCallbackArg(f model.Field) bool {
	return f.TypeName == "errorh.MetaCallback"
}

func IsPrimitiveArg(f model.Field) bool {
	return IsBoolArg(f) || IsIntArg(f) || IsStringArg(f) || IsStringSliceArg(f) || IsDateArg(f)
}

func IsBoolArg(f model.Field) bool {
	return f.IsBool()
}

func IsIntArg(f model.Field) bool {
	return f.IsInt()
}

func IsStringArg(f model.Field) bool {
	return f.IsString()
}

func IsStringSliceArg(f model.Field) bool {
	return f.IsStringSlice()
}

func IsDateArg(f model.Field) bool {
	return f.TypeName == "mydate.MyDate"
}

func ToFirstUpper(in string) string {
	a := []rune(in)
	a[0] = unicode.ToUpper(a[0])
	return string(a)
}
