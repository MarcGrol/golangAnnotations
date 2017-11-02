package main

import (
	"fmt"
	"log"

	"github.com/MarcGrol/golangAnnotations/generator/event"
	"github.com/MarcGrol/golangAnnotations/generator/eventService"
	"github.com/MarcGrol/golangAnnotations/generator/generationUtil"
	"github.com/MarcGrol/golangAnnotations/generator/jsonHelpers"
	"github.com/MarcGrol/golangAnnotations/generator/rest"
	"github.com/MarcGrol/golangAnnotations/model"
)

var registeredGenerators = make(map[string]generationUtil.Generator)

func init() {
	err := register("event", event.NewGenerator())
	if err != nil {
		log.Printf("Error registering structExample-annotation-generator")
	}

	err = register("eventService", eventService.NewGenerator())
	if err != nil {
		log.Printf("Error registering eventservice-annotation-generator")
	}

	err = register("json", jsonHelpers.NewGenerator())
	if err != nil {
		log.Printf("Error registering jsonHelpers-annotation-generator")
	}

	err = register("rest", rest.NewGenerator())
	if err != nil {
		log.Printf("Error registering rest-annotation-generator")

	}
}

func register(name string, generateFunc generationUtil.Generator) error {
	_, exists := registeredGenerators[name]
	if exists {
		return fmt.Errorf("Generator module %s already exists", name)
	}
	registeredGenerators[name] = generateFunc
	return nil
}

func runAllGenerators(inputDir string, parsedSources model.ParsedSources) error {
	for name, generator := range registeredGenerators {
		err := generator.Generate(inputDir, parsedSources)
		if err != nil {
			return fmt.Errorf("Error generating module %s: %s", name, err)
		}
	}
	return nil
}
