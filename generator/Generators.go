package generator

import (
	"fmt"

	"github.com/MarcGrol/golangAnnotations/model"
)

type GenerateFunc func(inputDir string, parsedSources model.ParsedSources) error

var registeredGenerators map[string]GenerateFunc = make(map[string]GenerateFunc)

func Register(name string, generateFunc GenerateFunc) error {
	_, exists := registeredGenerators[name]
	if exists {
		return fmt.Errorf("Generator module %s already exists", name)
	}
	registeredGenerators[name] = generateFunc
	return nil
}

func RunAllGenerators(inputDir string, parsedSources model.ParsedSources) error {
	for name, generateFunc := range registeredGenerators {
		err := generateFunc(inputDir, parsedSources)
		if err != nil {
			return fmt.Errorf("Error generating module %s: %s", name, err)
		}
	}
	return nil
}
