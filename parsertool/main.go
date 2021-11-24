package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/MarcGrol/golangAnnotations/generator"
	"github.com/MarcGrol/golangAnnotations/generator/ast"
	"github.com/MarcGrol/golangAnnotations/parser"
)

func main() {
	inputDir, outputFile := processArgs()

	excludeMatchPattern := "^" + generator.GenfilePrefix + ".*.go$"
	parsedSources, err := parser.New().ParseSourceDir(inputDir, "^.*.go$", excludeMatchPattern)
	if err != nil {
		log.Printf("Error parsing golang sources in %s: %s", inputDir, err)
		os.Exit(1)
	}

	jsonAstGenerator := ast.NewGenerator(outputFile)
	err = jsonAstGenerator.Generate(inputDir, parsedSources)
	if err != nil {
		log.Printf("Error generating json-ast: %s", err)
		os.Exit(-1)
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

func processArgs() (string, string) {
	inputDir := flag.String("input-dir", "", "Directory to be examined")
	outputFile := flag.String("output-file", "", "File jso-ast is written to")
	help := flag.Bool("help", false, "Usage information")

	flag.Parse()

	if help != nil && *help == true {
		printUsage()
	}
	if inputDir == nil || *inputDir == "" {
		printUsage()
	}

	return *inputDir, *outputFile
}
