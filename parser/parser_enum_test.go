package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEnumsInFile(t *testing.T) {

	parsedSources, err := ParseSourceFile("enums/enum.go")
	assert.Equal(t, nil, err)

	{
		assert.Equal(t, 2, len(parsedSources.Typedefs))

		assert.Equal(t, "// @Enum()", parsedSources.Typedefs[0].DocLines[0])
		assert.Equal(t, "ColorType", parsedSources.Typedefs[0].Name)
		assert.Equal(t, "int", parsedSources.Typedefs[0].Type)
		assert.Equal(t, "enums/enum.go", parsedSources.Typedefs[0].Filename)
		assert.Equal(t, "enums", parsedSources.Typedefs[0].PackageName)

		assert.Equal(t, "// @Enum()", parsedSources.Typedefs[1].DocLines[0])
		assert.Equal(t, "Profession", parsedSources.Typedefs[1].Name)
		assert.Equal(t, "string", parsedSources.Typedefs[1].Type)
		assert.Equal(t, "enums/enum.go", parsedSources.Typedefs[1].Filename)
		assert.Equal(t, "enums", parsedSources.Typedefs[1].PackageName)
	}

	{
		assert.Equal(t, 2, len(parsedSources.Enums))

		assert.Equal(t, "// @Enum()", parsedSources.Enums[0].DocLines[0])
		assert.Equal(t, "ColorType", parsedSources.Enums[0].Name)
		assert.Equal(t, "Red", parsedSources.Enums[0].EnumLiterals[0].Name)
		assert.Equal(t, "Green", parsedSources.Enums[0].EnumLiterals[1].Name)
		assert.Equal(t, "Blue", parsedSources.Enums[0].EnumLiterals[2].Name)
		assert.Equal(t, "enums/enum.go", parsedSources.Enums[0].Filename)
		assert.Equal(t, "enums", parsedSources.Enums[0].PackageName)

		assert.Equal(t, "// @Enum()", parsedSources.Enums[1].DocLines[0])
		assert.Equal(t, "Profession", parsedSources.Enums[1].Name)
		assert.Equal(t, "Teacher", parsedSources.Enums[1].EnumLiterals[0].Name)
		assert.Equal(t, "_teacher", parsedSources.Enums[1].EnumLiterals[0].Value)
		assert.Equal(t, "Cleaner", parsedSources.Enums[1].EnumLiterals[1].Name)
		assert.Equal(t, "_cleaner", parsedSources.Enums[1].EnumLiterals[1].Value)
		assert.Equal(t, "enums/enum.go", parsedSources.Enums[1].Filename)
		assert.Equal(t, "enums", parsedSources.Enums[1].PackageName)
	}
}
