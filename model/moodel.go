package model

type Service struct {
	DocLines     []string
	PackageName  string
	Name         string
	Methods      []Method
	CommentLines []string
}

type Method struct {
	DocLines     []string
	PackageName  string
	Service      *Service
	Name         string
	InputArgs    []Struct
	OutputArgs   []Struct
	CommentLines []string
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

func (s Struct) GetRestServiceParamaters() (path string) {
	return ""
}

func (m Method) IsRestMethod() bool {
	return false
}

func (s Struct) GetRestMethodParamaters() (path string, method string) {
	return "", "GET"
}
