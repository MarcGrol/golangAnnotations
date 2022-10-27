package generator

import (
	"github.com/MarcGrol/golangAnnotations/generator/annotation"
	"github.com/MarcGrol/golangAnnotations/model"
)

const (
	GenfilePrefix       = "gen_"
	GenfileExcludeRegex = GenfilePrefix + ".*"
)

type Generator interface {
	GetAnnotations() []annotation.AnnotationDescriptor
	Generate(inputDir string, parsedSources model.ParsedSources) error
}
