package generationUtil

import (
	"path"

	"github.com/f0rt/golangAnnotations/generator"
)

func Prefixed(filenamePath string) string {
	dir, filename := path.Split(filenamePath)
	return dir + generator.GenfilePrefix + filename
}
