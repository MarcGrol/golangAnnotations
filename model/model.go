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
	_, ok := resolveRestServiceAnnotation(s.DocLines)
	return ok
}

func (o Operation) IsRestOperation() bool {
	_, ok := resolveRestOperationAnnotation(o.DocLines)
	return ok
}

func (o Operation) GetRestOperationParamaters() (path string, method string) {
	val, _ := resolveRestOperationAnnotation(o.DocLines)
	return val["Method"], val["Path"]
}
