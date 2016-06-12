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
	RelatedStruct *Field // optional
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

type Interface struct {
	DocLines     []string
	PackageName  string
	Name         string
	Methods      []Operation
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
