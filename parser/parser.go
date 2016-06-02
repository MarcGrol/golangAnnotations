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
	//ast.Print(fset, f)
	v := structVisitor{}
	ast.Walk(&v, f)
	return v.structs, nil
}

func FindStructsInDir(dirName string, filenameRegex string) ([]model.Struct, error) {
	var pattern = regexp.MustCompile(filenameRegex)

	fset := token.NewFileSet()
	packages, err := parser.ParseDir(
		fset,
		dirName,
		func(fi os.FileInfo) bool {
			//log.Printf("filename:%s: matches %v", fi.Name(), pattern.MatchString(fi.Name()))
			return pattern.MatchString(fi.Name())
		},
		parser.ParseComments)
	if err != nil {
		log.Printf("error parsing dir %s: %s", dirName, err.Error())
		return []model.Struct{}, err
	}

	v := structVisitor{}
	for _, p := range packages {
		for _, f := range p.Files {
			//ast.Print(fset, f)
			ast.Walk(&v, f)
		}
	}
	return v.structs, nil
}

type structVisitor struct {
	packageName string
	docLines    []string
	name        string
	structs     []model.Struct
}

func (v *structVisitor) Visit(node ast.Node) ast.Visitor {
	if node != nil {
		//log.Printf("Got node:%+v", node)
		{
			ts, ok := node.(*ast.File)
			if ok {
				//log.Printf("*** Got file:%+v", ts)
				if ts.Name != nil {
					v.packageName = ts.Name.Name
				}
				return v
			}
		}
		{
			ts, ok := node.(*ast.GenDecl)
			if ok {
				//log.Printf("*** Got GenDecl:%+v", ts)
				if ts.Doc != nil && ts.Doc.List != nil && len(ts.Doc.List) > 0 {
					for _, d := range ts.Doc.List {
						v.docLines = append(v.docLines, d.Text)
					}
				}
				return v
			}
		}
		{
			ts, ok := node.(*ast.StructType)
			if ok {
				//log.Printf("*** Got StructType:%+v", ts)
				strct := handleStruct(v.packageName, v.name, v.docLines, ts)
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

func handleStruct(packageName string, name string, docLines []string, node *ast.StructType) model.Struct {
	myStruct := model.Struct{
		PackageName: packageName,
		Name:        name,
		Fields:      make([]model.Field, 0, 10),
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

type operationVisitor struct {
	packageName string
	operations  []model.Operation
}

func FindOperationsInDir(dirName string, filenameRegex string) ([]model.Operation, error) {
	fset := token.NewFileSet()
	packages, err := parser.ParseDir(
		fset,
		dirName,
		nil,
		parser.ParseComments)
	if err != nil {
		log.Printf("error parsing dir %s: %s", dirName, err.Error())
		return []model.Operation{}, err
	}

	v := operationVisitor{}
	for _, p := range packages {
		for _, f := range p.Files {
			ast.Print(fset, f)
			ast.Walk(&v, f)
		}
	}
	return v.operations, nil
}

func (v *operationVisitor) Visit(node ast.Node) ast.Visitor {
	if node != nil {
		{
			ts, ok := node.(*ast.File)
			if ok {
				if ts.Name != nil {
					v.packageName = ts.Name.Name
				}
				return v
			}
		}

		{
			fd, ok := node.(*ast.FuncDecl)
			if ok {
				oper := model.Operation{
					PackageName: v.packageName,
				}

				log.Printf("*** Got FuncDecl:%+v", fd)

				if fd.Doc != nil && len(fd.Doc.List) > 0 {
					docLines := []string{}
					for _, line := range fd.Doc.List {
						log.Printf("*** Got doc-line:%s", line.Text)
						docLines = append(docLines, line.Text)
					}
					oper.DocLines = docLines
				}

				if fd.Recv != nil && len(fd.Recv.List) >= 1 {
					relatedStruct, _ := extractField(fd.Recv.List[0])
					log.Printf("*** recv:%+v", relatedStruct)
					oper.RelatedStruct = relatedStruct
				}

				if fd.Name != nil {
					oper.Name = fd.Name.Name
					log.Printf("*** Got operation name:%s", oper.Name)
				}

				if fd.Type.Params != nil {
					args := []model.Field{}
					for _, p := range fd.Type.Params.List {
						arg, ok := extractField(p)
						if ok {
							args = append(args, arg)
						}
					}
					oper.InputArgs = args
					log.Printf("*** Got inputArgs:%+v", args)
				}

				if fd.Type.Results != nil {
					args := []model.Field{}
					for _, p := range fd.Type.Results.List {
						arg, ok := extractField(p)
						if ok {
							args = append(args, arg)
						}
					}
					oper.OutputArgs = args
					log.Printf("*** Got outputArgs:%+v", args)

				}

				log.Printf("oper: %+v", oper)
				v.operations = append(v.operations, oper)
			}
			return v
		}

	}
	return v
}

func extractField(input *ast.Field) (model.Field, bool) {
	field := model.Field{}
	if len(input.Names) >= 1 {
		// TODO should be able to handle multiple nammes
		field.Name = input.Names[0].Name
	}
	{
		param, ok := input.Type.(*ast.ArrayType)
		if ok {
			elt, ok := param.Elt.(*ast.Ident)
			if ok {
				field.TypeName = elt.Name
				field.IsSlice = true
			}
		}
	}
	{
		star, ok := input.Type.(*ast.StarExpr)
		if ok {
			x, ok := star.X.(*ast.Ident)
			if ok {
				field.TypeName = x.Name
				field.IsPointer = true
			}
		}
	}
	{
		param, ok := input.Type.(*ast.Ident)
		if ok {
			field.TypeName = param.Name
			field.IsSlice = false
		}
	}

	isOk := false
	if field.TypeName != "" {
		isOk = true
	}

	return field, isOk
}
