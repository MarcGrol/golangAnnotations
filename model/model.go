package model

type Service struct {
	DocLines     []string
	PackageName  string
	Name         string
	Operations   []Operation
	CommentLines []string
}

type Operation struct {
	DocLines      []string
	PackageName   string
	RelatedStruct *Field
	Name          string
	InputArgs     []Field
	OutputArgs    []Field
	CommentLines  []string
}

type Struct struct {
	DocLines     []string
	PackageName  string
	Name         string
	Fields       []Field
	CommentLines []string
}

type Field struct {
	DocLines     []string
	Name         string
	TypeName     string
	IsSlice      bool
	IsPointer    bool
	Tag          string
	CommentLines []string
}

func (s Struct) IsEvent() bool {
	_, ok := resolveEventAnnotation(s.DocLines)
	return ok
}

func (s Struct) GetAggregateName() string {
	val, _ := resolveEventAnnotation(s.DocLines)
	return val
}

func (s Struct) IsRestService() bool {
	return false
}

func (m Operation) IsRestOperation() bool {
	return false
}

func (s Struct) GetRestOperationParamaters() (path string, method string) {
	return "", "GET"
}
