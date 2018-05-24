package ast

import (
	"encoding/json"
	"io/ioutil"

	"github.com/Duxxie/platform/backend/lib/filegen"
	"github.com/MarcGrol/golangAnnotations/generator"
	"github.com/MarcGrol/golangAnnotations/generator/annotation"
	"github.com/MarcGrol/golangAnnotations/generator/event/eventAnnotation"
	"github.com/MarcGrol/golangAnnotations/model"
)

type Generator struct {
}

func NewGenerator() generator.Generator {
	return &Generator{}
}

func (eg *Generator) GetAnnotations() []annotation.AnnotationDescriptor {
	return eventAnnotation.Get()
}

func (eg *Generator) Generate(inputDir string, parsedSources model.ParsedSources) error {

	marshalled, err := json.MarshalIndent(parsedSources, "", "\t")
	if err != nil {
		panic(err)
	}
	targetFilename := filegen.Prefixed(inputDir + "/" + "ast.json")
	err = ioutil.WriteFile(targetFilename, marshalled, 0644)
	if err != nil {
		panic(err)
	}

	return nil
}
