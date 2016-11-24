package generator

import (
	"fmt"
	"log"

	"github.com/MarcGrol/golangAnnotations/generator/event"
	"github.com/MarcGrol/golangAnnotations/generator/eventService"
	"github.com/MarcGrol/golangAnnotations/generator/rest"
	"github.com/MarcGrol/golangAnnotations/model"
)

type GenerateFunc func(inputDir string, parsedSources model.ParsedSources) error

var registeredGenerators map[string]GenerateFunc = make(map[string]GenerateFunc)

func init() {
	err := register("event", event.Generate)
	if err != nil {
		log.Printf("Error registering event-annotation-generator")
	}

	err = register("rest", rest.Generate)
	if err != nil {
		log.Printf("Error registering rest-annotation-generator")

	}

	err = register("eventService", eventService.Generate)
	if err != nil {
		log.Printf("Error registering eventservice-annotation-generator")
	}
}

func register(name string, generateFunc GenerateFunc) error {
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
		log.Printf("Run generation module %s", name)

	}
	return nil
}
