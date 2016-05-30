package model

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
