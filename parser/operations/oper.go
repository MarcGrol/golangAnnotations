package operations

import (
	"golang.org/x/net/context"
	"github.com/MarcGrol/golangAnnotations/parser/structs"
)

type Person struct {
	Name string
}

// docline for Service
type Service struct {
}

// docline for getPersons
func (serv *Service) getPersons(ctx context.Context) ([]Person, error) {
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
		Y:42,
	}
	return &p, nil
}
