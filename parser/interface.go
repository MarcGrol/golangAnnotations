package parser

import "github.com/f0rt/golangAnnotations/model"

type Parser interface {
	ParseSourceDir(dirName string, includeRegex string, excludeRegex string) (model.ParsedSources, error)
}
