package model

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func Parse(filename string) (ParsedSources, error) {
	parsedSources := ParsedSources{}

	var reader io.Reader = bufio.NewReader(os.Stdin)
	if filename != "" {
		fp, err := os.Open(filename)
		if err != nil {
			return ParsedSources{}, fmt.Errorf("Error opening file %s: %s", filename, err)
		}
		reader = bufio.NewReader(fp)
	}

	err := json.NewDecoder(reader).Decode(&parsedSources)
	if err != nil {
		return ParsedSources{}, fmt.Errorf("Error decoding parsed-sources from stdin: %s", err)
	}
	return parsedSources, nil
}
