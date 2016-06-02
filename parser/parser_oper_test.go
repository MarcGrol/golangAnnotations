package parser

import (
	"testing"

	"github.com/MarcGrol/astTools/model"
	"github.com/stretchr/testify/assert"
)

func TestStructOperationsInDir(t *testing.T) {
	opers, err := FindOperationsInDir("./operations", ".*")
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(opers))

	{
		o := opers[0]
		assert.Equal(t, "operations", o.PackageName)
		assert.Equal(t, []string{"// +Operation: -> methhod: GET, path: /person"}, o.DocLines)
		assert.Equal(t, "getPersons", o.Name)
		assertField(t, model.Field{Name: "serv", TypeName: "Service", IsPointer: true}, o.RelatedStruct)

		assert.Equal(t, 0, len(o.InputArgs))

		assert.Equal(t, 2, len(o.OutputArgs))
		assertField(t, model.Field{TypeName: "Person", IsSlice: true}, o.OutputArgs[0])
		assertField(t, model.Field{TypeName: "error"}, o.OutputArgs[1])
	}
	{
		o := opers[1]
		assert.Equal(t, "operations", o.PackageName)
		assert.Equal(t, []string{"// +Operation: -> methhod: GET, path: /person/:uid"}, o.DocLines)
		assert.Equal(t, "getPerson", o.Name)
		assertField(t, model.Field{Name: "s", TypeName: "Service"}, o.RelatedStruct)

		assert.Equal(t, 1, len(o.InputArgs))
		assertField(t, model.Field{Name: "uid", TypeName: "string"}, o.InputArgs[0])

		assert.Equal(t, 3, len(o.OutputArgs))
		assertField(t, model.Field{TypeName: "Person"}, o.OutputArgs[0])
		assertField(t, model.Field{TypeName: "Person", IsPointer: true}, o.OutputArgs[1])
		assertField(t, model.Field{TypeName: "error"}, o.OutputArgs[2])
	}
}
