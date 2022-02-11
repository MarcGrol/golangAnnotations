package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/f0rt/golangAnnotations/generator"
	"github.com/f0rt/golangAnnotations/generator/ast"
	"github.com/f0rt/golangAnnotations/generator/event"
	"github.com/f0rt/golangAnnotations/generator/eventService"
	"github.com/f0rt/golangAnnotations/generator/jsonHelpers"
	"github.com/f0rt/golangAnnotations/generator/repository"
	"github.com/f0rt/golangAnnotations/generator/rest"
	"github.com/f0rt/golangAnnotations/model"
	"github.com/f0rt/golangAnnotations/parser"
)

const (
	version = "0.8"

	excludeMatchPattern = "^" + generator.GenfilePrefix + ".*.go$"
)

var inputDir *string

func main() {
	processArgs()

	parsedSources, err := parser.New().ParseSourceDir(*inputDir, "^.*.go$", excludeMatchPattern)
	if err != nil {
		log.Printf("Error parsing golang sources in %s: %s", *inputDir, err)
		os.Exit(1)
	}

	runAllGenerators(*inputDir, parsedSources)

	os.Exit(0)
}

func runAllGenerators(inputDir string, parsedSources model.ParsedSources) {
	for name, g := range map[string]generator.Generator{
		"ast":           ast.NewGenerator("ast.json"),
		"event":         event.NewGenerator(),
		"event-service": eventService.NewGenerator(),
		"json-helpers":  jsonHelpers.NewGenerator(),
		"rest":          rest.NewGenerator(),
		"repository":    repository.NewGenerator(),
	} {
		err := g.Generate(inputDir, parsedSources)
		if err != nil {
			log.Printf("Error generating module %s: %s", name, err)
			os.Exit(-1)
		}
	}
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
