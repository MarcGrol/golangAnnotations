package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/MarcGrol/astTools/generator"

	"github.com/MarcGrol/astTools/parser"
)

const (
	VERSION = "0.1"
)

var (
	inputDir *string
)

func main() {
	processArgs()

	structs, err := parser.FindStructsInDir(*inputDir, ".*.go")
	if err != nil {
		log.Printf("Error finding structs in %s:%s", *inputDir, err)
	}
	log.Printf("structs found:%+v", structs)

	generator.GenerateForStructs(*inputDir, structs)
	// for _, s := range structs {
	// 	generator.GenerateForStruct(*inputDir, s)
	// }
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
