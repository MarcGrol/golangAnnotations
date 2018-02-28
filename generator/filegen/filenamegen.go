package filegen

import (
	"path"
)

const genfilePrefix = "gen_"

func MatchPattern() string {
	return "^[^" + genfilePrefix + "][^_]+\\.go$"
}

func Prefixed(filenamePath string) string {
	dir, filename := path.Split(filenamePath)
	return dir + genfilePrefix + filename
}
