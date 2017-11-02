package parserUtil

import "github.com/MarcGrol/golangAnnotations/model"

type Parser interface {
	ParseSourceDir(dirName string, filenameRegex string) (model.ParsedSources, error)
}
