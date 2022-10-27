package generator

import (
	"github.com/f0rt/golangAnnotations/generator/annotation"
	"github.com/f0rt/golangAnnotations/model"
)

const (
	GenfilePrefix       = "gen_"
	GenfileExcludeRegex = GenfilePrefix + ".*"
)

type Generator interface {
	GetAnnotations() []annotation.AnnotationDescriptor
	Generate(inputDir string, parsedSources model.ParsedSources) error
}
