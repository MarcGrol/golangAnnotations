package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/MarcGrol/golangAnnotations/parser"
)

const (
	version = "0.7"
)

var (
	inputDir *string
)

func main() {
	processArgs()

	parsedSources, err := parser.ParseSourceDir(*inputDir, "^[^\\$][^_]+\\.go$")
	if err != nil {
		log.Printf("Error parsing golang sources in %s:%s", *inputDir, err)
		os.Exit(1)
	}

	marshalled, err := json.MarshalIndent(parsedSources, "", "\t")
	if err != nil {
		panic(err)
	}
	targetFilename := *inputDir + "/$" + "ast.json"
	err = ioutil.WriteFile(targetFilename, marshalled, 0644)
	if err != nil {
		panic(err)
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
