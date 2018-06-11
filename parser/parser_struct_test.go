package parser

import (
	"testing"

	"github.com/MarcGrol/golangAnnotations/model"
	"github.com/stretchr/testify/assert"
)

func TestParseStructsInFile(t *testing.T) {

	parsedSources, err := parseSourceFile("structs/example.go")
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(parsedSources.Structs))

	assertStruct(t,
		model.Struct{PackageName: "structs", Name: "Person", DocLines: []string{"// Struct comment before type"}},
		parsedSources.Structs[0])
	assert.Equal(t, 13, len(parsedSources.Structs[0].Fields))

	{
		s := parsedSources.Structs[0]

		assertField(t,
			model.Field{Name: "FirstName", TypeName: "string", IsPointer: false, IsSlice: false},
			s.Fields[0])

		assertField(t,
			model.Field{Name: "LastName", TypeName: "string", IsPointer: false, IsSlice: false},
			s.Fields[1])

		assertField(t,
			model.Field{Name: "Age", TypeName: "int", IsPointer: false, IsSlice: false, CommentLines: []string{"// Age comment"}},
			s.Fields[2])

		assertField(t,
			model.Field{Name: "Nice", TypeName: "bool", IsPointer: true, IsSlice: false, DocLines: []string{"// Before nice comment"}, CommentLines: []string{"// After Nice comment"}},
			s.Fields[3])

		assertField(t,
			model.Field{Name: "Color", TypeName: "ColorType", IsPointer: false, IsSlice: false, DocLines: []string{"// Before Color comment"}, Tag: "`json:\"COLOR_TYPE\"`"},
			s.Fields[4])

		assertField(t,
			model.Field{Name: "OptionalColor", TypeName: "ColorType", IsPointer: true, IsSlice: false},
			s.Fields[5])

		assertField(t,
			model.Field{Name: "Father", TypeName: "Person", IsPointer: true, IsSlice: false},
			s.Fields[6])

		assertField(t,
			model.Field{Name: "Uncles", TypeName: "Person", IsPointer: true, IsSlice: true},
			s.Fields[7])

		assertField(t,
			model.Field{Name: "Children", TypeName: "Person", IsPointer: false, IsSlice: true},
			s.Fields[8])

		assertField(t,
			model.Field{Name: "ChildMap", TypeName: "map[string]Person", IsPointer: false, IsSlice: false},
			s.Fields[9])

		assertField(t,
			model.Field{Name: "ChildPointerMap", TypeName: "map[string]*Person", IsPointer: false, IsSlice: false},
			s.Fields[10])

		assertField(t,
			model.Field{Name: "ChildrenMap", TypeName: "map[string][]Person", IsPointer: false, IsSlice: false},
			s.Fields[11])

		assertField(t,
			model.Field{Name: "ChildrenPointerMap", TypeName: "map[string][]*Person", IsPointer: false, IsSlice: false},
			s.Fields[12])
	}
}

func TestParseStructsInDir(t *testing.T) {
	parsedSources, err := New().ParseSourceDir("structs", "^.*xample.go$", "gen_.*")
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, len(parsedSources.Structs))

	// Order is undetermined
	for _, s := range parsedSources.Structs {
		if s.Name == "Person" {
			assert.Equal(t, 13, len(s.Fields))
		}
		if s.Name == "MyStruct" {
			assert.Equal(t, 1, len(s.Fields))
		}
		if s.Name == "OtherStruct" {
			assert.Equal(t, 2, len(s.Fields))
		}
	}
}

func assertStruct(t *testing.T, expected model.Struct, actual model.Struct) {
	//t.Logf("expected: %+v, actual: %+v", expected, actual)
	assertStringSlice(t, expected.DocLines, actual.DocLines)
	assert.Equal(t, expected.PackageName, actual.PackageName)
	assert.Equal(t, expected.Name, actual.Name)
	assertStringSlice(t, expected.CommentLines, actual.CommentLines)
}

func assertField(t *testing.T, expected model.Field, actual model.Field) {
	//t.Logf("expected: %+v, actual: %+v", expected, actual)
	assertStringSlice(t, expected.DocLines, actual.DocLines)

	assert.Equal(t, expected.PackageName, actual.PackageName)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.TypeName, actual.TypeName)
	assert.Equal(t, expected.IsPointer, actual.IsPointer)
	assert.Equal(t, expected.IsSlice, actual.IsSlice)
	assert.Equal(t, expected.Tag, actual.Tag)
	assert.Equal(t, len(expected.CommentLines), len(actual.CommentLines))
	assertStringSlice(t, expected.CommentLines, actual.CommentLines)
}

func assertStringSlice(t *testing.T, expected []string, actual []string) {
	//t.Logf("expected: %+v, actual: %+v", expected, actual)
	actualHas := false
	if actual != nil && len(actual) > 0 {
		actualHas = true
	}
	expectedHas := false
	if expected != nil && len(expected) > 0 {
		expectedHas = true
	}

	assert.Equal(t, expectedHas, actualHas)
	if expected != nil && actual != nil {
		assert.Equal(t, len(expected), len(actual))
		for idx, s := range expected {
			assert.Equal(t, s, actual[idx])
		}
	}
}
