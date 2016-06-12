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
	v := astVisitor{}
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

	v := astVisitor{}
	for _, p := range packages {
		for _, f := range p.Files {
			//ast.Print(fset, f)
			ast.Walk(&v, f)
		}
	}
	return v.structs, nil
}

func FindOperationsInDir(dirName string, filenameRegex string) ([]model.Operation, error) {
	fset := token.NewFileSet()
	packages, err := parseDir(dirName, filenameRegex)
	if err != nil {
		return []model.Operation{}, err
	}

	v := astVisitor{}
	for _, p := range packages {
		for _, f := range p.Files {
			ast.Print(fset, f)
			ast.Walk(&v, f)
		}
	}
	return v.operations, nil
}

func FindInterfacesInDir(dirName string, filenameRegex string) ([]model.Interface, error) {
	fset := token.NewFileSet()
	packages, err := parseDir(dirName, filenameRegex)
	if err != nil {
		return []model.Interface{}, err
	}

	v := astVisitor{}
	for _, p := range packages {
		for _, f := range p.Files {
			ast.Print(fset, f)
			ast.Walk(&v, f)
		}
	}
	return v.interfaces, nil
}

func parseDir(dirName string, filenameRegex string) (map[string]*ast.Package, error) {
	packages := make(map[string]*ast.Package)
	var err error

	fset := token.NewFileSet()
	packages, err = parser.ParseDir(
		fset,
		dirName,
		nil,
		parser.ParseComments)
	if err != nil {
		log.Printf("error parsing dir %s: %s", dirName, err.Error())
		return packages, err
	}
	return packages, nil
}

type astVisitor struct {
	packageName string
	structs     []model.Struct
	operations  []model.Operation
	interfaces  []model.Interface
}

func (v *astVisitor) Visit(node ast.Node) ast.Visitor {
	if node != nil {

		// package-name is in isolated node
		pName, found := extractPackageName(node)
		if found {
			v.packageName = pName
		}

		{
			// if struct, get its fields
			str, found := extractGenDeclForStruct(node)
			if found {
				str.PackageName = v.packageName
				v.structs = append(v.structs, str)
			}
		}

		{
			// if interfaces, get its methods
			iface, found := extractGenDecForInterfacel(node)
			if found {
				for _, m := range iface.Methods {
					m.PackageName = v.packageName
				}
				iface.PackageName = v.packageName
				v.interfaces = append(v.interfaces, iface)
			}
		}

		{
			// if operation, get its signature
			operation, ok := extractOperation(node)
			if ok {
				operation.PackageName = v.packageName
				v.operations = append(v.operations, operation)
			}
		}
	}
	return v
}

func extractGenDeclForStruct(node ast.Node) (model.Struct, bool) {
	found := false
	var str model.Struct

	gd, ok := node.(*ast.GenDecl)
	if ok {
		log.Printf("*** Got ast.GenDecl:%+v", gd)

		// Continue parsing to see if it a struct
		str, found = extractSpecsForStruct(gd.Specs)
		if ok {
			// Docline of struct (that could contain annotations) appear far before the details of the struct
			str.DocLines = extractDocLines(gd.Doc)
		}
	}

	return str, found
}

func extractGenDecForInterfacel(node ast.Node) (model.Interface, bool) {
	found := false
	var iface model.Interface

	gd, ok := node.(*ast.GenDecl)
	if ok {
		log.Printf("*** Got ast.GenDecl:%+v", gd)

		// Continue parsing to see if it an interface
		iface, found = extractSpecsForInterface(gd.Specs)
		if ok {
			// Docline of interface (that could contain annotations) appear far before the details of the struct
			iface.DocLines = extractDocLines(gd.Doc)
		}
	}

	return iface, found
}

func extractSpecsForStruct(specs []ast.Spec) (model.Struct, bool) {
	found := false
	str := model.Struct{}

	if len(specs) >= 1 {
		ts, ok := specs[0].(*ast.TypeSpec)
		if ok {
			log.Printf("*** Got ast.TypeSpec:%+v", ts)

			str.Name = ts.Name.Name

			ss, ok := ts.Type.(*ast.StructType)
			if ok {
				log.Printf("*** Got ast.StructType:%+v", ss)

				str.Fields = extractFieldList(ss.Fields)
				found = true
			}
		}
	}

	return str, found
}

func extractSpecsForInterface(specs []ast.Spec) (model.Interface, bool) {
	found := false
	interf := model.Interface{}

	if len(specs) >= 1 {
		ts, ok := specs[0].(*ast.TypeSpec)
		if ok {
			log.Printf("*** Got ast.TypeSpec:%+v", ts)

			interf.Name = ts.Name.Name

			it, ok := ts.Type.(*ast.InterfaceType)
			if ok {
				log.Printf("************************************")
				log.Printf("*** Got ast.InterfaceType:%+v", interf)
				log.Printf("************************************")
				interf.Methods = extractInterfaceMethods(it.Methods)
				found = true
			}
		}
	}

	return interf, found
}

func extractPackageName(node ast.Node) (string, bool) {
	name := ""

	fil, found := node.(*ast.File)
	if found {

		if fil.Name != nil {
			name = fil.Name.Name

		}
		log.Printf("*** Got ast.File:%+v ->", fil, name)

	}
	return name, found
}

func extractOperation(node ast.Node) (model.Operation, bool) {
	found := false
	oper := model.Operation{}

	fd, found := node.(*ast.FuncDecl)
	if found {
		log.Printf("*** Got FuncDecl:%+v", fd)

		oper.DocLines = extractDocLines(fd.Doc)

		if fd.Recv != nil {
			recvd := extractFieldList(fd.Recv)
			if len(recvd) >= 1 {
				oper.RelatedStruct = &(recvd[0])
			}
			log.Printf("*** Got RelatedStruct:%+v", oper.RelatedStruct)
		}

		if fd.Name != nil {
			oper.Name = fd.Name.Name
			log.Printf("*** Got operation name:%s", oper.Name)
		}

		if fd.Type.Params != nil {
			oper.InputArgs = extractFieldList(fd.Type.Params)
			log.Printf("*** Got inputArgs:%+v", oper.InputArgs)
		}

		if fd.Type.Results != nil {
			oper.OutputArgs = extractFieldList(fd.Type.Results)
			log.Printf("*** Got outputArgs:%+v", oper.OutputArgs)
		}

		log.Printf("oper: %+v", oper)
	}
	return oper, found
}

func extractDocLines(doc *ast.CommentGroup) []string {
	docLines := []string{}
	if doc != nil {
		for _, line := range doc.List {
			docLines = append(docLines, line.Text)
		}
		log.Printf("*** Got doc-lines:%+v", docLines)
	}
	return docLines
}

func extractComments(comment *ast.CommentGroup) []string {
	lines := []string{}
	if comment != nil {
		for _, c := range comment.List {
			lines = append(lines, c.Text)
		}
		log.Printf("*** Got Comment:%+v", lines)
	}
	return lines
}

func extractTag(tag *ast.BasicLit) (string, bool) {
	if tag != nil {
		log.Printf("*** Got Tag:%+v", tag.Value)
		return tag.Value, true
	}
	return "", false
}

func extractFieldList(fl *ast.FieldList) []model.Field {
	fields := []model.Field{}
	if fl != nil {
		for _, p := range fl.List {
			flds := extractFields(p)
			fields = append(fields, flds...)
		}
	}
	return fields
}

func extractInterfaceMethods(fl *ast.FieldList) []model.Operation {
	methods := []model.Operation{}

	for _, m := range fl.List {
		if len(m.Names) > 0 {
			oper := model.Operation{DocLines: extractDocLines(m.Doc)}

			log.Printf("*** Got interface name:%+v", m.Names[0].Name)
			oper.Name = m.Names[0].Name

			ft, found := m.Type.(*ast.FuncType)
			if found {
				if ft.Params != nil {
					oper.InputArgs = extractFieldList(ft.Params)
					log.Printf("*** Got inputArgs:%+v", oper.InputArgs)
				}

				if ft.Results != nil {
					oper.OutputArgs = extractFieldList(ft.Results)
					log.Printf("*** Got outputArgs:%+v", oper.OutputArgs)
				}
				log.Printf("interface:%+v", oper)
				methods = append(methods, oper)
			}
		}
	}
	return methods
}

func extractFields(input *ast.Field) []model.Field {
	fields := []model.Field{}
	if input != nil {
		if len(input.Names) == 0 {
			field := _extractField(input)
			fields = append(fields, field)
		} else {
			// A single field can refer to multiple: example: x,y int -> x int, y int
			for _, name := range input.Names {
				field := _extractField(input)
				field.Name = name.Name
				fields = append(fields, field)
			}
		}
	}
	return fields
}

func _extractField(input *ast.Field) model.Field {
	field := model.Field{}

	field.DocLines = extractDocLines(input.Doc)

	field.CommentLines = extractComments(input.Comment)

	tag, found := extractTag(input.Tag)
	if found {
		field.Tag = tag
	}
	{
		arr, ok := input.Type.(*ast.ArrayType)
		if ok {
			log.Printf("*** Got ast.ArrayType:%+v", arr)
			field.IsSlice = true
			{
				ident, ok := arr.Elt.(*ast.Ident)
				if ok {
					log.Printf("*** Got ast.Ident:%+v", ident)
					field.TypeName = ident.Name
				}
			}
			{
				star, ok := arr.Elt.(*ast.StarExpr)
				if ok {
					log.Printf("*** Got ast.StarExpr:%+v", star)
					ident, ok := star.X.(*ast.Ident)
					if ok {
						log.Printf("*** Got ast.Ident:%+v", ident)
						field.TypeName = ident.Name
						field.IsPointer = true
					}
				}
			}
		}
	}
	{
		star, ok := input.Type.(*ast.StarExpr)
		if ok {
			log.Printf("*** Got ast.StarExpr:%+v", star)
			ident, ok := star.X.(*ast.Ident)
			if ok {
				log.Printf("*** Got ast.Ident:%+v", ident)
				field.TypeName = ident.Name
				field.IsPointer = true
			}
		}
	}
	{
		ident, ok := input.Type.(*ast.Ident)
		if ok {
			log.Printf("*** Got ast.Ident:%+v", ident)
			field.TypeName = ident.Name
		}
	}

	return field
}
