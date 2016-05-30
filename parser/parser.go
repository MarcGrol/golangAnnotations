package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"regexp"
)

type Struct struct {
	PackageName string
	Name        string
	Fields      []Field
}

type Field struct {
	Name      string
	TypeName  string
	IsSlice   bool
	IsPointer bool
}

func FindStructsInFile(srcFilename string) ([]Struct, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, srcFilename, nil, parser.ParseComments)
	if err != nil {
		log.Printf("error parsing src %s: %s", srcFilename, err.Error())
		return []Struct{}, err
	}
	//ast.Print(fset, f)
	v := structVisitor{}
	ast.Walk(&v, f)
	return v.Structs, nil
}

func FindStructsInDir(dirName string, filenameRegex string) ([]Struct, error) {
	var pattern = regexp.MustCompile(filenameRegex)

	fset := token.NewFileSet()
	packages, err := parser.ParseDir(
		fset,
		dirName,
		func(fi os.FileInfo) bool {
			//log.Printf("filename:%s", fi.Name())
			return pattern.MatchString(fi.Name())
		},
		0)
	if err != nil {
		log.Printf("error parsing dir %s: %s", dirName, err.Error())
		return []Struct{}, err
	}

	v := structVisitor{}
	for _, p := range packages {
		for _, f := range p.Files {
			ast.Walk(&v, f)
		}
	}
	return v.Structs, nil
}

type structVisitor struct {
	Name    string
	Structs []Struct
}

func (v *structVisitor) Visit(node ast.Node) ast.Visitor {
	if node != nil {
		{
			ts, ok := node.(*ast.TypeSpec)
			if ok {
				v.Name = ts.Name.Name
				return v
			}
		}
		{
			ts, ok := node.(*ast.StructType)
			if ok {
				v.Structs = append(v.Structs, handleStruct(v.Name, ts))
				v.Name = ""
			}
			return v
		}

	}
	return v
}

func handleStruct(name string, node *ast.StructType) Struct {
	myStruct := Struct{
		Name:   name,
		Fields: make([]Field, 0, 10),
	}

	for _, rawField := range node.Fields.List {
		fields, ok := handleFields(rawField)
		if ok {
			for _, f := range fields {
				myStruct.Fields = append(myStruct.Fields, f)
			}
		}
	}

	return myStruct
}

func handleFields(node ast.Node) ([]Field, bool) {

	// we are looking for a node of type ield
	ts, ok := node.(*ast.Field)
	if !ok {
		return []Field{}, false
	}

	dataType := ""
	isPointer := false
	isSlice := false

	{
		// array
		slice, ok := ts.Type.(*ast.ArrayType)
		if ok {
			isSlice = true
			{
				elt, ok := slice.Elt.(*ast.StarExpr)
				if ok {
					isPointer = true
					sliceDataType, ok := elt.X.(*ast.Ident)
					if ok {
						dataType = sliceDataType.Name
					}
				}
			}
			{
				elt, ok := slice.Elt.(*ast.Ident)
				if ok {
					dataType = elt.Name
				}
			}
		}
	}

	{
		// pointer
		star, ok := ts.Type.(*ast.StarExpr)
		if ok {
			isPointer = true
			pointerDataType, ok := star.X.(*ast.Ident)
			if ok {
				dataType = pointerDataType.Name
			}
		}
	}
	{
		// no pointer, no array
		t, ok := ts.Type.(*ast.Ident)
		if ok {
			dataType = t.Name
		}
	}

	fields := make([]Field, 0, 10)
	for _, f := range ts.Names {
		field := Field{
			Name:      f.Name,
			TypeName:  dataType,
			IsSlice:   isSlice,
			IsPointer: isPointer,
		}
		fields = append(fields, field)
	}
	return fields, true
}
