package ast

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/MarcGrol/golangAnnotations/generator"
	"github.com/MarcGrol/golangAnnotations/generator/annotation"
	"github.com/MarcGrol/golangAnnotations/generator/event/eventAnnotation"
	"github.com/MarcGrol/golangAnnotations/generator/generationUtil"
	"github.com/MarcGrol/golangAnnotations/model"
)

type Generator struct {
	targetFilename string
}

func NewGenerator(targetFilename string) generator.Generator {
	return &Generator{
		targetFilename: targetFilename,
	}
}

func (eg *Generator) GetAnnotations() []annotation.AnnotationDescriptor {
	return eventAnnotation.Get()
}

func (eg *Generator) Generate(inputDir string, parsedSources model.ParsedSources) error {

	marshalled, err := json.MarshalIndent(parsedSources, "", "\t")
	if err != nil {
		panic(err)
	}

	if eg.targetFilename != "" {
		filenamePath := generationUtil.Prefixed(inputDir + "/" + eg.targetFilename)
		err = ioutil.WriteFile(filenamePath, marshalled, 0644)
		if err != nil {
			return fmt.Errorf("Error writing json-ast to file:%s", err)
		}
	} else {
		_, err = os.Stdout.Write(marshalled)
		if err != nil {
			return fmt.Errorf("Error writing json-ast to stdout:%s", err)
		}
	}

	return nil
}
