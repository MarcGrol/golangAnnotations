package operations

type Person struct {
	Name string
}

// +Service -> path: /api/
type Service struct {
}

// +Operation: -> methhod: GET, path: /person
func (serv *Service) getPersons() ([]Person, error) {
	return []Person{
		{Name: "Marc"},
		{Name: "Eva"},
	}, nil
}

// +Operation: -> methhod: GET, path: /person/:uid
func (s Service) getPerson(uid string) (Person, *Person, error) {
	p := Person{
		Name: "Pien",
	}
	return p, &p, nil
}
