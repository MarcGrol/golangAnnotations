package operations

type Person struct {
	Name string
}

// docline for Service
type Service struct {
}

// docline for getPersons
func (serv *Service) getPersons() ([]Person, error) {
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
