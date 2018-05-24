package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/MarcGrol/golangAnnotations/generator/ast"
	"github.com/MarcGrol/golangAnnotations/generator/event"
	"github.com/MarcGrol/golangAnnotations/generator/eventService"
	"github.com/MarcGrol/golangAnnotations/generator/filtering"
	"github.com/MarcGrol/golangAnnotations/generator/generationUtil"
	"github.com/MarcGrol/golangAnnotations/generator/jsonHelpers"
	"github.com/MarcGrol/golangAnnotations/generator/repository"
	"github.com/MarcGrol/golangAnnotations/generator/rest"
	"github.com/MarcGrol/golangAnnotations/model"
	"github.com/MarcGrol/golangAnnotations/parser"
)

const (
	version = "0.7"
)

var inputDir *string

func main() {
	processArgs()

	parsedSources, err := parser.New().ParseSourceDir(*inputDir, "^.*.go$", filtering.ExcludeMatchPattern())
	if err != nil {
		log.Printf("Error parsing golang sources in %s:%s", *inputDir, err)
		os.Exit(1)
	}

	runAllGenerators(*inputDir, parsedSources)

	os.Exit(0)
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "\nUsage:\n")
	fmt.Fprintf(os.Stderr, " %s [flags]\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func printVersion() {
	fmt.Fprintf(os.Stderr, "\nVersion: %s\n", version)
	os.Exit(1)
}

func processArgs() {
	inputDir = flag.String("input-dir", "", "Directory to be examined")
	help := flag.Bool("help", false, "Usage information")
	version := flag.Bool("version", false, "Version information")

	flag.Parse()

	if help != nil && *help == true {
		printUsage()
	}
	if version != nil && *version == true {
		printVersion()
	}
	if inputDir == nil || *inputDir == "" {
		printUsage()
	}
}

func runAllGenerators(inputDir string, parsedSources model.ParsedSources) error {
	for name, g := range map[string]generationUtil.Generator{
		"ast":           ast.NewGenerator(),
		"event":         event.NewGenerator(),
		"event-service": eventService.NewGenerator(),
		"json-helpers":  jsonHelpers.NewGenerator(),
		"rest":          rest.NewGenerator(),
		"repository":    repository.NewGenerator(),
	} {
		err := g.Generate(inputDir, parsedSources)
		if err != nil {
			return fmt.Errorf("Error generating module %s: %s", name, err)
		}
	}
	return nil
}
