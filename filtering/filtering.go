package filtering

const genfilePrefix = "gen_"

func ExcludeMatchPattern() string {
	return "^" + genfilePrefix + ".*.go$"
}
