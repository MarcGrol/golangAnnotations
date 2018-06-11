package model

//go:generate golangAnnotations -input-dir .

// @JsonStruct()
type ParsedSources struct {
	Structs    []Struct    `json:"structs,omitempty"`
	Operations []Operation `json:"operations,omitempty"`
	Interfaces []Interface `json:"interfaces,omitempty"`
	Typedefs   []Typedef   `json:"typedefs,omitempty"`
	Enums      []Enum      `json:"enums,omitempty"`
}

// @JsonStruct()
type Operation struct {
	PackageName   string   `json:"packageName,omitempty"`
	Filename      string   `json:"filename,omitempty"`
	DocLines      []string `json:"docLines,omitempty"`
	RelatedStruct *Field   `json:"relatedStruct,omitempty"` // optional
	Name          string   `json:"name"`
	InputArgs     []Field  `json:"inputArgs,omitempty"`
	OutputArgs    []Field  `json:"outputArgs,omitempty"`
	CommentLines  []string `json:"commentLines,omitempty"`
}

// @JsonStruct()
type Struct struct {
	PackageName  string       `json:"packageName"`
	Filename     string       `json:"filename"`
	DocLines     []string     `json:"docLines,omitempty"`
	Name         string       `json:"name"`
	Fields       []Field      `json:"fields,omitempty"`
	Operations   []*Operation `json:"operations,omitempty"`
	CommentLines []string     `json:"commentLines,omitempty"`
}

// @JsonStruct()
type Interface struct {
	PackageName  string      `json:"packageName"`
	Filename     string      `json:"filename"`
	DocLines     []string    `json:"docLines,omitempty"`
	Name         string      `json:"name"`
	Methods      []Operation `json:"methods,omitempty"`
	CommentLines []string    `json:"commentLines,omitempty"`
}

// @JsonStruct()
type Field struct {
	PackageName string   `json:"packageName,omitempty"`
	DocLines    []string `json:"docLines,omitempty"`
	Name        string   `json:"name,omitempty"`
	TypeName    string   `json:"typeName,omitempty"`
	IsSlice     bool     `json:"isSlice,omitempty"`
	IsPointer   bool     `json:"isPointer,omitempty"`
	Tag          string   `json:"tag,omitempty"`
	CommentLines []string `json:"commentLines,omitempty"`
}

// @JsonStruct()
type Typedef struct {
	PackageName string   `json:"packageName"`
	Filename    string   `json:"filename"`
	DocLines    []string `json:"docLines,omitempty"`
	Name        string   `json:"name"`
	Type        string   `json:"type,omitempty"`
}

// @JsonStruct()
type Enum struct {
	PackageName  string        `json:"packageName"`
	Filename     string        `json:"filename"`
	DocLines     []string      `json:"docLines,omitempty"`
	Name         string        `json:"name,omitempty"`
	EnumLiterals []EnumLiteral `json:"enumLiterals,omitempty"`
	CommentLines []string      `json:"commentLines,omitempty"`
}

// @JsonStruct()
type EnumLiteral struct {
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
}
