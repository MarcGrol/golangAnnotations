package filegen

import (
	"path"

	"github.com/MarcGrol/golangAnnotations/generator"
)

func Prefixed(filenamePath string) string {
	dir, filename := path.Split(filenamePath)
	return dir + generator.GenfilePrefix + filename
}
