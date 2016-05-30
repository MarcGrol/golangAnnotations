package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"regexp"

	"github.com/MarcGrol/astTools/model"
)

func FindStructsInFile(srcFilename string) ([]model.Struct, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, srcFilename, nil, parser.ParseComments)
	if err != nil {
		log.Printf("error parsing src %s: %s", srcFilename, err.Error())
		return []model.Struct{}, err
	}
	ast.Print(fset, f)
	v := structVisitor{}
	ast.Walk(&v, f)
	return v.structs, nil
}

func findStructsInDir(dirName string, filenameRegex string) ([]model.Struct, error) {
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
		return []model.Struct{}, err
	}

	v := structVisitor{}
	for _, p := range packages {
		for _, f := range p.Files {
			ast.Walk(&v, f)
		}
	}
	return v.structs, nil
}

type structVisitor struct {
	docLines []string
	name     string
	structs  []model.Struct
}

func (v *structVisitor) Visit(node ast.Node) ast.Visitor {
	if node != nil {
		{
			ts, ok := node.(*ast.GenDecl)
			if ok {
				if ts.Doc != nil && ts.Doc.List != nil && len(ts.Doc.List) > 0 {
					for _, d := range ts.Doc.List {
						v.docLines = append(v.docLines, d.Text)
					}
				}
				return v
			}
		}
		{
			ts, ok := node.(*ast.TypeSpec)
			if ok {
				v.name = ts.Name.Name
				return v
			}
		}
		{
			ts, ok := node.(*ast.StructType)
			if ok {
				strct := handleStruct(v.name, v.docLines, ts)
				if len(v.docLines) > 0 {
					strct.DocLines = v.docLines
				}
				v.structs = append(v.structs, strct)
				v.name = ""
				v.docLines = []string{}
			}
			return v
		}

	}
	return v
}

func handleStruct(name string, docLines []string, node *ast.StructType) model.Struct {
	myStruct := model.Struct{
		Name:   name,
		Fields: make([]model.Field, 0, 10),
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

func handleFields(node ast.Node) ([]model.Field, bool) {

	// we are looking for a node of type ield
	ts, ok := node.(*ast.Field)
	if !ok {
		return []model.Field{}, false
	}

	docLines := []string{}
	tag := ""
	dataType := ""
	isPointer := false
	isSlice := false
	commentLines := []string{}

	if ts.Doc != nil && len(ts.Doc.List) > 0 {
		for _, d := range ts.Doc.List {
			docLines = append(docLines, d.Text)
		}
	}

	if ts.Comment != nil && len(ts.Comment.List) > 0 {
		for _, c := range ts.Comment.List {
			commentLines = append(commentLines, c.Text)
		}
	}

	if ts.Tag != nil {
		tag = ts.Tag.Value
	}

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

	fields := make([]model.Field, 0, 10)
	for _, f := range ts.Names {
		field := model.Field{
			DocLines:     docLines,
			Name:         f.Name,
			TypeName:     dataType,
			IsSlice:      isSlice,
			IsPointer:    isPointer,
			Tag:          tag,
			CommentLines: commentLines,
		}
		fields = append(fields, field)
	}
	return fields, true
}
