package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/f0rt/golangAnnotations/generator/rest"
	"github.com/f0rt/golangAnnotations/model"
)

func main() {
	inputFile, outputDir, triggerRestGenerator := processArgs()

	parsedSources, err := model.Parse(inputFile)
	if err != nil {
		log.Printf("Error parsing parsed-sources as json in %s: %s", inputFile, err)
		os.Exit(-1)
	}

	if triggerRestGenerator {
		generator := rest.NewGenerator()
		err = generator.Generate(outputDir, parsedSources)
		if err != nil {
			log.Printf("Error triggering rest-generator: %s", err)
			os.Exit(-2)
		}
	}

	os.Exit(0)
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "\nUsage:\n")
	fmt.Fprintf(os.Stderr, " %s [flags]\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func processArgs() (string, string, bool) {
	inputFile := flag.String("input-file", "", "Parsed sources (as json) to be read")
	outputDir := flag.String("output-dir", "", "Target directory where being written to")
	userRestGenerator := flag.Bool("use-generator-rest", false, "Trigger rest generator")
	help := flag.Bool("help", false, "Usage information")
	flag.Parse()

	if help != nil && *help == true {
		printUsage()
	}
	if inputFile == nil {
		*inputFile = ""
	}
	if outputDir == nil || *outputDir == "" {
		*outputDir = path.Dir(*inputFile)
	}

	return *inputFile, *outputDir, *userRestGenerator
}
