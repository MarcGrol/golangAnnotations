package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/MarcGrol/golangAnnotations/generator/event"
	"github.com/MarcGrol/golangAnnotations/generator/rest"
	"github.com/MarcGrol/golangAnnotations/parser"
)

const (
	VERSION = "0.2"
)

var (
	inputDir *string
)

func main() {
	processArgs()

	harvest, err := parser.ParseSourceDir(*inputDir, ".*.go")
	if err != nil {
		log.Printf("Error finding structs in %s:%s", *inputDir, err)
		os.Exit(1)
	}

	err = event.Generate(*inputDir, harvest.Structs)
	if err != nil {
		log.Printf("Error generating event code:%s", err)
		os.Exit(1)
	}

	err = rest.Generate(*inputDir, harvest.Structs)
	if err != nil {
		log.Printf("Error generating rest code:%s", err)
		os.Exit(1)
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

func printVersion() {
	fmt.Fprintf(os.Stderr, "\nVersion: %s\n", VERSION)
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
