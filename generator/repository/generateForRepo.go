package repository

import (
	"fmt"
	"log"
	"strings"
	"text/template"
	"unicode"

	"github.com/f0rt/golangAnnotations/generator"
	"github.com/f0rt/golangAnnotations/generator/annotation"
	"github.com/f0rt/golangAnnotations/generator/generationUtil"
	"github.com/f0rt/golangAnnotations/generator/repository/repositoryAnnotation"
	"github.com/f0rt/golangAnnotations/model"
)

type Generator struct {
}

func NewGenerator() generator.Generator {
	return &Generator{}
}

func (eg *Generator) GetAnnotations() []annotation.AnnotationDescriptor {
	return repositoryAnnotation.Get()
}

func (eg *Generator) Generate(inputDir string, parsedSource model.ParsedSources) error {
	structs := parsedSource.Structs

	packageName, err := generationUtil.GetPackageNameForStructs(structs)
	if packageName == "" || err != nil {
		return err
	}
	targetDir, err := generationUtil.DetermineTargetPath(inputDir, packageName)
	if err != nil {
		return err
	}
	for _, repository := range structs {
		if IsRepository(repository) {
			err = generationUtil.Generate(generationUtil.Info{
				Src:            fmt.Sprintf("%s.%s", repository.PackageName, repository.Name),
				TargetFilename: generationUtil.Prefixed(fmt.Sprintf("%s/%s.go", targetDir, toFirstLower(repository.Name))),
				TemplateName:   "repository",
				TemplateString: repositoryTemplate,
				FuncMap:        customTemplateFuncs,
				Data:           repository,
			})
			if err != nil {
				log.Fatalf("Error generating repository %s: %s", repository.Name, err)
				return err
			}
		}
	}
	return nil
}

var customTemplateFuncs = template.FuncMap{
	"IsRepository":              IsRepository,
	"AggregateNameConst":        AggregateNameConst,
	"LowerAggregateName":        LowerAggregateName,
	"UpperAggregateName":        UpperAggregateName,
	"GetPackageName":            GetPackageName,
	"LowerModelName":            LowerModelName,
	"UpperModelName":            UpperModelName,
	"ModelPackageName":          ModelPackageName,
	"HasMethodFind":             HasMethodFind,
	"HasMethodFilterByEvent":    HasMethodFilterByEvent,
	"HasMethodFilterByMoment":   HasMethodFilterByMoment,
	"HasMethodFindStates":       HasMethodFindStates,
	"HasMethodExists":           HasMethodExists,
	"HasMethodAllAggregateUIDs": HasMethodAllAggregateUIDs,
	"HasMethodGetAllAggregates": HasMethodGetAllAggregates,
	"HasMethodPurgeOnEventUIDs": HasMethodPurgeOnEventUIDs,
	"HasMethodPurgeOnEventType": HasMethodPurgeOnEventType,
	"HasMethodPurgeAll":         HasMethodPurgeAll,
}

func IsRepository(s model.Struct) bool {
	annotations := annotation.NewRegistry(repositoryAnnotation.Get())
	_, ok := annotations.ResolveAnnotationByName(s.DocLines, repositoryAnnotation.TypeRepository)
	return ok
}

func AggregateNameConst(s model.Struct) string {
	return fmt.Sprintf("%sAggregateName", UpperAggregateName(s))
}

func LowerAggregateName(s model.Struct) string {
	return toFirstLower(GetAggregateName(s))
}

func UpperAggregateName(s model.Struct) string {
	return toFirstUpper(GetAggregateName(s))
}

func GetAggregateName(s model.Struct) string {
	annotations := annotation.NewRegistry(repositoryAnnotation.Get())
	if ann, ok := annotations.ResolveAnnotationByName(s.DocLines, repositoryAnnotation.TypeRepository); ok {
		return ann.Attributes[repositoryAnnotation.ParamAggregate]
	}
	return ""
}

func GetPackageName(s model.Struct) string {
	annotations := annotation.NewRegistry(repositoryAnnotation.Get())
	if ann, ok := annotations.ResolveAnnotationByName(s.DocLines, repositoryAnnotation.TypeRepository); ok {
		packageName := ann.Attributes[repositoryAnnotation.ParamPackage]
		if packageName != "" {
			return packageName
		}
	}
	return fmt.Sprintf("%sEvents", LowerAggregateName(s))
}

func LowerModelName(s model.Struct) string {
	return toFirstLower(GetModelName(s))
}

func UpperModelName(s model.Struct) string {
	return toFirstUpper(GetModelName(s))
}

func ModelPackageName(s model.Struct) string {
	return toFirstLower(GetModelName(s)) + "Model"
}

func GetModelName(s model.Struct) string {
	annotations := annotation.NewRegistry(repositoryAnnotation.Get())
	if ann, ok := annotations.ResolveAnnotationByName(s.DocLines, repositoryAnnotation.TypeRepository); ok {
		m := ann.Attributes[repositoryAnnotation.ParamModel]
		if m != "" {
			return m
		}
	}
	return GetAggregateName(s)
}

func HasMethodFind(s model.Struct) bool {
	return HasMethod(s, "find")
}

func HasMethodFilterByEvent(s model.Struct) bool {
	return HasMethod(s, "filterByEvent")
}

func HasMethodFilterByMoment(s model.Struct) bool {
	return HasMethod(s, "filterByMoment")
}

func HasMethodFindStates(s model.Struct) bool {
	return HasMethod(s, "findStates")
}

func HasMethodExists(s model.Struct) bool {
	return HasMethod(s, "exists")
}

func HasMethodAllAggregateUIDs(s model.Struct) bool {
	return HasMethod(s, "allAggregateUIDs")
}

func HasMethodGetAllAggregates(s model.Struct) bool {
	return HasMethod(s, "allAggregates")
}

func HasMethodPurgeOnEventUIDs(s model.Struct) bool {
	return HasMethod(s, "purgeOnEventUIDs")
}

func HasMethodPurgeOnEventType(s model.Struct) bool {
	return HasMethod(s, "purgeOnEventType")
}

func HasMethodPurgeAll(s model.Struct) bool {
	return HasMethod(s, "purgeAll")
}

func HasMethod(s model.Struct, methodName string) bool {
	annotations := annotation.NewRegistry(repositoryAnnotation.Get())
	if ann, ok := annotations.ResolveAnnotationByName(s.DocLines, repositoryAnnotation.TypeRepository); ok {
		methods := strings.Split(ann.Attributes[repositoryAnnotation.ParamMethods], ",")
		for _, method := range methods {
			if strings.TrimSpace(method) == methodName {
				return true
			}
		}
	}
	return false
}

func toFirstLower(in string) string {
	a := []rune(in)
	a[0] = unicode.ToLower(a[0])
	return string(a)
}

func toFirstUpper(in string) string {
	a := []rune(in)
	a[0] = unicode.ToUpper(a[0])
	return string(a)
}
