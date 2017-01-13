package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEnumsInFile(t *testing.T) {

	harvest, err := ParseSourceFile("enums/enum.go")
	assert.Equal(t, nil, err)

	{
		assert.Equal(t, 2, len(harvest.Typedefs))

		assert.Equal(t, "// @Enum()", harvest.Typedefs[0].DocLines[0])
		assert.Equal(t, "ColorType", harvest.Typedefs[0].Name)
		assert.Equal(t, "int", harvest.Typedefs[0].Type)

		assert.Equal(t, "// @Enum()", harvest.Typedefs[1].DocLines[0])
		assert.Equal(t, "Profession", harvest.Typedefs[1].Name)
		assert.Equal(t, "string", harvest.Typedefs[1].Type)
	}

	{
		assert.Equal(t, 2, len(harvest.Enums))

		assert.Equal(t, "// @Enum()", harvest.Enums[0].DocLines[0])
		assert.Equal(t, "ColorType", harvest.Enums[0].Name)
		assert.Equal(t, "Red", harvest.Enums[0].EnumLiterals[0].Name)
		assert.Equal(t, "Green", harvest.Enums[0].EnumLiterals[1].Name)
		assert.Equal(t, "Blue", harvest.Enums[0].EnumLiterals[2].Name)

		assert.Equal(t, "// @Enum()", harvest.Enums[1].DocLines[0])
		assert.Equal(t, "Profession", harvest.Enums[1].Name)
		assert.Equal(t, "Teacher", harvest.Enums[1].EnumLiterals[0].Name)
		assert.Equal(t, "_teacher", harvest.Enums[1].EnumLiterals[0].Value)
		assert.Equal(t, "Cleaner", harvest.Enums[1].EnumLiterals[1].Name)
		assert.Equal(t, "_cleaner", harvest.Enums[1].EnumLiterals[1].Value)
	}
}
