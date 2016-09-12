package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"regexp"

	"github.com/MarcGrol/golangAnnotations/model"
	"path/filepath"
	"strings"
)

var (
	debugAstOfSources = false
)

func ParseSourceFile(srcFilename string) (model.ParsedSources, error) {
	if debugAstOfSources {
		dumpFile(srcFilename)
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, srcFilename, nil, parser.ParseComments)
	if err != nil {
		log.Printf("error parsing src %s: %s", srcFilename, err.Error())
		return model.ParsedSources{}, err
	}
	v := astVisitor{
		Imports:map[string]string{},
	}
	ast.Walk(&v, f)

	embedMethodsInStructs(&v)

	result := model.ParsedSources{
		Structs:    v.Structs,
		Operations: v.Operations,
		Interfaces: v.Interfaces,
	}
	return result, nil
}

func ParseSourceDir(dirName string, filenameRegex string) (model.ParsedSources, error) {
	if debugAstOfSources {
		dumpFilesInDir(dirName)
	}
	packages, err := parseDir(dirName, filenameRegex)
	if err != nil {
		log.Printf("error parsing dir %s: %s", dirName, err.Error())
		return model.ParsedSources{}, err
	}

	v := astVisitor{
		Imports: map[string]string{},
	}
	for _, p := range packages {
		for _, f := range p.Files {
			ast.Walk(&v, f)
		}
	}

	embedMethodsInStructs(&v)

	result := model.ParsedSources{
		Structs:    v.Structs,
		Operations: v.Operations,
		Interfaces: v.Interfaces,
	}

	return result, nil
}

func embedMethodsInStructs(visitor *astVisitor) {
	allStructs := make(map[string]*model.Struct)
	for idx, _ := range visitor.Structs {
		allStructs[(&visitor.Structs[idx]).Name] = &visitor.Structs[idx]
	}
	for idx, _ := range visitor.Operations {
		oper := visitor.Operations[idx]
		if oper.RelatedStruct != nil {
			found, exists := allStructs[(*oper.RelatedStruct).TypeName]
			if exists {
				found.Operations = append(found.Operations, &oper)
			}
		}
	}

}

func parseDir(dirName string, filenameRegex string) (map[string]*ast.Package, error) {
	var pattern = regexp.MustCompile(filenameRegex)

	packages := make(map[string]*ast.Package)
	var err error

	fset := token.NewFileSet()
	packages, err = parser.ParseDir(
		fset,
		dirName,
		func(fi os.FileInfo) bool {
			return pattern.MatchString(fi.Name())
		},
		parser.ParseComments)
	if err != nil {
		log.Printf("error parsing dir %s: %s", dirName, err.Error())
		return packages, err
	}

	return packages, nil
}

func dumpFile(srcFilename string) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, srcFilename, nil, parser.ParseComments)
	if err != nil {
		log.Printf("error parsing src %s: %s", srcFilename, err.Error())
		return
	}
	ast.Print(fset, f)
}

func dumpFilesInDir(dirName string) {
	fset := token.NewFileSet()
	packages, err := parser.ParseDir(
		fset,
		dirName,
		nil,
		parser.ParseComments)
	if err != nil {
		log.Printf("error parsing dir %s: %s", dirName, err.Error())
	}
	for _, p := range packages {
		for _, f := range p.Files {
			ast.Print(fset, f)
		}
	}
}

type astVisitor struct {
	PackageName string
	Imports 	map[string]string
	Structs     []model.Struct
	Operations  []model.Operation
	Interfaces  []model.Interface
}

func (v *astVisitor) Visit(node ast.Node) ast.Visitor {
	if node != nil {

		// package-name is in isolated node
		pName, found := extractPackageName(node)
		if found {
			v.PackageName = pName
		}

		// extract all imports into a map
		v.extractGenDeclImports(node)

		{
			// if struct, get its fields
			str, found := extractGenDeclForStruct(node, v.Imports )
			if found {
				str.PackageName = v.PackageName
				v.Structs = append(v.Structs, str)
			}
		}

		{
			// if interfaces, get its methods
			iface, found := extractGenDecForInterface(node, v.Imports )
			if found {
				iface.PackageName = v.PackageName
				v.Interfaces = append(v.Interfaces, iface)
			}
		}

		{
			// if operation, get its signature
			operation, ok := extractOperation(node, v.Imports )
			if ok {
				operation.PackageName = v.PackageName
				v.Operations = append(v.Operations, operation)
			}
		}

	}
	return v
}


func (v *astVisitor)extractGenDeclImports(node ast.Node) {

	gd, ok := node.(*ast.GenDecl)
	if ok {
		for _, spec := range gd.Specs {
				is, ok := spec.(*ast.ImportSpec)
				if ok {
					quotedImport := is.Path.Value
					unquotedImport := strings.Trim(quotedImport, "\"")
					first, last := filepath.Split(unquotedImport)
					if first == "" {
						last = first
					}
					v.Imports[last] = unquotedImport
					log.Printf( "Found import %s -> %s",  last, unquotedImport)
				}
		}
	}
}


func extractGenDeclForStruct(node ast.Node, imports map[string]string) (model.Struct, bool) {
	found := false
	var str model.Struct

	gd, ok := node.(*ast.GenDecl)
	if ok {
		// Continue parsing to see if it a struct
		str, found = extractSpecsForStruct(gd.Specs, imports)
		if ok {
			// Docline of struct (that could contain annotations) appear far before the details of the struct
			str.DocLines = extractDocLines(gd.Doc)
		}
	}

	return str, found
}

func extractGenDecForInterface(node ast.Node, imports map[string]string) (model.Interface, bool) {
	found := false
	var iface model.Interface

	gd, ok := node.(*ast.GenDecl)
	if ok {
		// Continue parsing to see if it an interface
		iface, found = extractSpecsForInterface(gd.Specs, imports)
		if ok {
			// Docline of interface (that could contain annotations) appear far before the details of the struct
			iface.DocLines = extractDocLines(gd.Doc)
		}
	}

	return iface, found
}

func extractSpecsForStruct(specs []ast.Spec, imports map[string]string) (model.Struct, bool) {
	found := false
	str := model.Struct{}

	if len(specs) >= 1 {
		ts, ok := specs[0].(*ast.TypeSpec)
		if ok {
			str.Name = ts.Name.Name

			ss, ok := ts.Type.(*ast.StructType)
			if ok {
				str.Fields = extractFieldList(ss.Fields, imports)
				found = true
			}
		}
	}

	return str, found
}

func extractSpecsForInterface(specs []ast.Spec, imports map[string]string) (model.Interface, bool) {
	found := false
	interf := model.Interface{}

	if len(specs) >= 1 {
		ts, ok := specs[0].(*ast.TypeSpec)
		if ok {
			interf.Name = ts.Name.Name

			it, ok := ts.Type.(*ast.InterfaceType)
			if ok {
				interf.Methods = extractInterfaceMethods(it.Methods, imports)
				found = true
			}
		}
	}

	return interf, found
}

func extractPackageName(node ast.Node) (string, bool) {
	found := false
	packageName := ""

	f, found := node.(*ast.File)
	if found {
		if f.Name != nil {
			packageName = f.Name.Name
		}
	}
	return packageName, found
}

func extractOperation(node ast.Node, imports map[string]string) (model.Operation, bool) {
	found := false
	oper := model.Operation{}

	fd, found := node.(*ast.FuncDecl)
	if found {
		oper.DocLines = extractDocLines(fd.Doc)

		if fd.Recv != nil {
			recvd := extractFieldList(fd.Recv, imports)
			if len(recvd) >= 1 {
				oper.RelatedStruct = &(recvd[0])
			}
		}

		if fd.Name != nil {
			oper.Name = fd.Name.Name
		}

		if fd.Type.Params != nil {
			oper.InputArgs = extractFieldList(fd.Type.Params, imports)
		}

		if fd.Type.Results != nil {
			oper.OutputArgs = extractFieldList(fd.Type.Results, imports)
		}
	}
	return oper, found
}

func extractDocLines(doc *ast.CommentGroup) []string {
	docLines := []string{}
	if doc != nil {
		for _, line := range doc.List {
			docLines = append(docLines, line.Text)
		}
	}
	return docLines
}

func extractComments(comment *ast.CommentGroup) []string {
	lines := []string{}
	if comment != nil {
		for _, c := range comment.List {
			lines = append(lines, c.Text)
		}
	}
	return lines
}

func extractTag(tag *ast.BasicLit) (string, bool) {
	if tag != nil {
		return tag.Value, true
	}
	return "", false
}

func extractFieldList(fl *ast.FieldList, imports map[string]string) []model.Field {
	fields := []model.Field{}
	if fl != nil {
		for _, p := range fl.List {
			flds := extractFields(p, imports)
			fields = append(fields, flds...)
		}
	}
	return fields
}

func extractInterfaceMethods(fl *ast.FieldList, imports map[string]string) []model.Operation {
	methods := []model.Operation{}

	for _, m := range fl.List {
		if len(m.Names) > 0 {
			oper := model.Operation{DocLines: extractDocLines(m.Doc)}

			oper.Name = m.Names[0].Name

			ft, found := m.Type.(*ast.FuncType)
			if found {
				if ft.Params != nil {
					oper.InputArgs = extractFieldList(ft.Params, imports)
				}

				if ft.Results != nil {
					oper.OutputArgs = extractFieldList(ft.Results, imports)
				}
				methods = append(methods, oper)
			}
		}
	}
	return methods
}

func extractFields(input *ast.Field, imports map[string]string) []model.Field {
	fields := []model.Field{}
	if input != nil {
		if len(input.Names) == 0 {
			field := _extractField(input,imports)
			fields = append(fields, field)
		} else {
			// A single field can refer to multiple: example: x,y int -> x int, y int
			for _, name := range input.Names {
				field := _extractField(input,imports)
				field.Name = name.Name
				fields = append(fields, field)
			}
		}
	}
	return fields
}

func _extractField(input *ast.Field, imports map[string]string) model.Field {
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
			field.IsSlice = true
			{
				ident, ok := arr.Elt.(*ast.Ident)
				if ok {
					field.TypeName = ident.Name
				}
				sel, ok := arr.Elt.(*ast.SelectorExpr)
				if ok {
					ident, ok = sel.X.(*ast.Ident)
					if ok {
						field.TypeName = fmt.Sprintf("%s.%s", ident.Name, sel.Sel.Name)
						field.PackageName = imports[ident.Name]
					}
				}
			}
			{
				star, ok := arr.Elt.(*ast.StarExpr)
				if ok {
					ident, ok := star.X.(*ast.Ident)
					if ok {
						field.TypeName = ident.Name
						field.IsPointer = true
					}
				}
			}
		}
	}

	{
		var mapKey string = ""
		var mapValue string = ""

		mapType, ok := input.Type.(*ast.MapType)
		if ok {
			{
				key, ok := mapType.Key.(*ast.Ident)
				if ok {
					mapKey = key.Name
				}
			}
			{
				value, ok := mapType.Value.(*ast.Ident)
				if ok {
					mapValue = value.Name
				}
			}
		}
		if mapKey != "" && mapValue != "" {
			field.TypeName = fmt.Sprintf("map[%s]%s", mapKey, mapValue)
		}

	}

	{
		star, ok := input.Type.(*ast.StarExpr)
		if ok {
			ident, ok := star.X.(*ast.Ident)
			if ok {
				//log.Printf("star ident: %+v", ident.Name)
				field.TypeName = ident.Name
				field.IsPointer = true
			}
			sel, ok := star.X.(*ast.SelectorExpr)
			if ok {
				ident, ok = sel.X.(*ast.Ident)
				if ok {
					field.TypeName = fmt.Sprintf( "%s.%s", ident.Name, sel.Sel.Name)
					field.IsPointer = true
					field.PackageName = imports[ident.Name]
				}
			}
		}
	}
	{
		ident, ok := input.Type.(*ast.Ident)
		if ok {
			field.TypeName = ident.Name
		}
	}
	{
		sel, ok := input.Type.(*ast.SelectorExpr)
		if ok {
			x, ok := sel.X.(*ast.Ident)
			if ok {
				field.Name = x.Name
				field.TypeName = fmt.Sprintf("%s.%s", x.Name, sel.Sel.Name)
				field.PackageName = imports[x.Name]
			}
		}
	}

	return field
}
