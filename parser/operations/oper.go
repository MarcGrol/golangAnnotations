package operations

import (
	"context"

	"github.com/f0rt/golangAnnotations/parser/structs"
)

type Person struct {
	Name string
}

// docline for Service
type Service struct {
}

// docline for getPersons
func (s *Service) getPersons(
	// ctx
	ctx context.Context, data map[string]string, slice []string) ([]Person, error) {
	return []Person{
		{Name: "Marc"},
		{Name: "Eva"},
	}, nil
}

// docline for getPerson
func (s Service) getPerson(uid string) (Person, *Person, error) {
	p := Person{
		Name: "Pien",
	}
	return p, &p, nil
}

// docline for getForeignStruct
func (s Service) getForeignStruct(in structs.YetAnotherStruct) (*structs.YetAnotherStruct, error) {
	p := structs.YetAnotherStruct{
		Y: 42,
	}
	return &p, nil
}

// docline for getForeignStructs
func (s Service) getForeignStructs(ctx context.Context) ([]*structs.YetAnotherStruct, error) {
	p := &structs.YetAnotherStruct{
		Y: 42,
	}
	return []*structs.YetAnotherStruct{p}, nil
}
