package parser

import (
	"testing"

	"github.com/MarcGrol/golangAnnotations/generator"
	"github.com/MarcGrol/golangAnnotations/model"
	"github.com/stretchr/testify/assert"
)

func TestInterfacesInDir(t *testing.T) {
	parsedSources, err := New().ParseSourceDir("./interfaces", "^.*.go$", generator.GenfileExcludeRegex)
	assert.Equal(t, nil, err)
	assert.Len(t, parsedSources.Interfaces, 1)

	{
		i := parsedSources.Interfaces[0]
		assert.Equal(t, "interfaces", i.PackageName)
		assert.Equal(t, []string{"// docline for interface Doer"}, i.DocLines)
		assert.Equal(t, "Doer", i.Name)

		{
			assert.Len(t, i.Methods, 2)
			{
				m := i.Methods[0]
				assert.Equal(t, []string{"// docline for interface method doit"}, m.DocLines)
				assert.Equal(t, "doit", m.Name)
				assert.Nil(t, m.RelatedStruct)
				assert.Equal(t, 2, len(m.InputArgs))
				assertField(t, model.Field{Name: "c", PackageName: "golang.org/x/net/context", TypeName: "context.Context", IsSlice: false}, m.InputArgs[0])
				assertField(t, model.Field{Name: "req", TypeName: "Req", IsSlice: false}, m.InputArgs[1])

				assert.Equal(t, 2, len(m.OutputArgs))
				assertField(t, model.Field{TypeName: "Resp", IsSlice: false}, m.OutputArgs[0])
				assertField(t, model.Field{TypeName: "error"}, m.OutputArgs[1])
			}
			{
				m := i.Methods[1]
				assert.Equal(t, []string{"// docline for interface method dontDoit"}, m.DocLines)
				assert.Equal(t, "dontDoit", m.Name)
				assert.Nil(t, m.RelatedStruct)
				assert.Equal(t, 0, len(m.InputArgs))
				assert.Equal(t, 0, len(m.OutputArgs))
			}
		}
	}
}
