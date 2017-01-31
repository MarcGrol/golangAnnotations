package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/MarcGrol/golangAnnotations/model"
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
	v := &astVisitor{
		Imports: map[string]string{},
	}
	v.CurrentFilename = srcFilename
	ast.Walk(v, f)

	embedOperationsInStructs(v)

	embedTypedefDocLinesInEnum(v)

	result := model.ParsedSources{
		Structs:    v.Structs,
		Operations: v.Operations,
		Interfaces: v.Interfaces,
		Typedefs:   v.Typedefs,
		Enums:      v.Enums,
	}
	return result, nil
}

type FileEntry struct {
	key  string
	file ast.File
}

type FileEntries []FileEntry

func (list FileEntries) Len() int {
	return len(list)
}

func (list FileEntries) Less(i, j int) bool {
	return list[i].key < list[j].key
}

func (list FileEntries) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func SortedFileEntries(fileMap map[string]*ast.File) FileEntries {
	var fileEntries FileEntries = make([]FileEntry, 0, len(fileMap))
	for key, file := range fileMap {
		if file != nil {
			fileEntries = append(fileEntries, FileEntry{
				key:  key,
				file: *file,
			})
		}
	}
	sort.Sort(fileEntries)
	return fileEntries
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

	v := &astVisitor{
		Imports: map[string]string{},
	}
	for _, p := range packages {
		for _, entry := range SortedFileEntries(p.Files) {
			v.CurrentFilename = entry.key
			ast.Walk(v, &entry.file)
		}
	}

	embedOperationsInStructs(v)

	embedTypedefDocLinesInEnum(v)

	result := model.ParsedSources{
		Structs:    v.Structs,
		Operations: v.Operations,
		Interfaces: v.Interfaces,
		Typedefs:   v.Typedefs,
		Enums:      v.Enums,
	}

	return result, nil
}

func embedOperationsInStructs(visitor *astVisitor) {
	allStructs := make(map[string]*model.Struct)
	for idx := range visitor.Structs {
		allStructs[(&visitor.Structs[idx]).Name] = &visitor.Structs[idx]
	}
	for idx := range visitor.Operations {
		oper := visitor.Operations[idx]
		if oper.RelatedStruct != nil {
			found, exists := allStructs[(*oper.RelatedStruct).TypeName]
			if exists {
				found.Operations = append(found.Operations, &oper)
			}
		}
	}

}

func embedTypedefDocLinesInEnum(visitor *astVisitor) {
	for idx, e := range visitor.Enums {
		for _, td := range visitor.Typedefs {
			if td.Name == e.Name {
				visitor.Enums[idx].DocLines = td.DocLines
				break
			}
		}
	}
}

func parseDir(dirName string, filenameRegex string) (map[string]*ast.Package, error) {
	var pattern = regexp.MustCompile(filenameRegex)

	packages := make(map[string]*ast.Package)
	var err error

	fileSet := token.NewFileSet()
	packages, err = parser.ParseDir(
		fileSet,
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
	CurrentFilename string
	PackageName     string
	Filename        string
	Imports         map[string]string
	Structs         []model.Struct
	Operations      []model.Operation
	Interfaces      []model.Interface
	Typedefs        []model.Typedef
	Enums           []model.Enum
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
			str, found := extractGenDeclForStruct(node, v.Imports)
			if found {
				str.PackageName = v.PackageName
				str.Filename = v.CurrentFilename
				v.Structs = append(v.Structs, str)
			}
		}

		{
			// if struct, get its fields
			td, found := extractGenDeclForTypedef(node, v.Imports)
			if found {
				td.PackageName = v.PackageName
				td.Filename = v.CurrentFilename
				v.Typedefs = append(v.Typedefs, td)
			}
		}
		{
			// if struct, get its fields
			e, found := extractGenDeclForEnum(node, v.Imports)
			if found {
				e.PackageName = v.PackageName
				e.Filename = v.CurrentFilename
				v.Enums = append(v.Enums, e)
			}
		}
		{
			// if interfaces, get its methods
			iface, found := extractGenDecForInterface(node, v.Imports)
			if found {
				iface.PackageName = v.PackageName
				iface.Filename = v.CurrentFilename
				v.Interfaces = append(v.Interfaces, iface)
			}
		}

		{
			// if operation, get its signature
			operation, ok := extractOperation(node, v.Imports)
			if ok {
				operation.PackageName = v.PackageName
				operation.Filename = v.CurrentFilename
				v.Operations = append(v.Operations, operation)
			}
		}

	}
	return v
}

func (v *astVisitor) extractGenDeclImports(node ast.Node) {

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
				//log.Printf( "Found import %s -> %s",  last, unquotedImport)
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

func extractGenDeclForTypedef(node ast.Node, imports map[string]string) (model.Typedef, bool) {
	found := false
	var td model.Typedef

	gd, ok := node.(*ast.GenDecl)
	if ok {
		// Continue parsing to see if it a struct
		td, found = extractSpecsForTypedef(gd.Specs, imports)
		if found {
			td.DocLines = extractDocLines(gd.Doc)
		}
	}

	return td, found
}

func extractGenDeclForEnum(node ast.Node, imports map[string]string) (model.Enum, bool) {
	found := false
	var e model.Enum

	gd, ok := node.(*ast.GenDecl)
	if ok {
		// Continue parsing to see if it an enum
		e, found = extractSpecsForEnum(gd.Specs, imports)
		// Docs live in the related typdef
	}

	return e, found
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

func extractSpecsForEnum(specs []ast.Spec, imports map[string]string) (model.Enum, bool) {
	found := false
	enumeration := model.Enum{}

	// parse type part

	// parse const part
	if len(specs) >= 1 {
		isEnumConstant := false
		typeName := ""
		for _, vs := range specs {
			s, ok := vs.(*ast.ValueSpec)
			if ok {
				if s.Type != nil {
					for _, n := range s.Names {
						i, ok := s.Type.(*ast.Ident)
						if ok {
							typeName = i.Name
						}
						if n.Obj.Kind == ast.Con {
							isEnumConstant = true
							break
						}
					}
				}
			}
		}

		if isEnumConstant {

			enumeration.Name = typeName
			enumeration.EnumLiterals = []model.EnumLiteral{}
			for _, vs := range specs {
				s, ok := vs.(*ast.ValueSpec)
				if ok {
					var data *int = nil
					if s.Names[0].Obj != nil {
						i, ok := s.Names[0].Obj.Data.(int)
						if ok {
							data = &i
						}
					}

					literal := model.EnumLiteral{
						Name: s.Names[0].Name,
						Data: data,
					}

					for _, v := range s.Values {

						b, ok := v.(*ast.BasicLit)
						if ok {
							literal.Value = strings.Trim(b.Value, "\"")
							break
						}
					}
					enumeration.EnumLiterals = append(enumeration.EnumLiterals, literal)
				}
			}
			found = true
		}
	}

	return enumeration, found
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

func extractSpecsForTypedef(specs []ast.Spec, imports map[string]string) (model.Typedef, bool) {
	found := false
	td := model.Typedef{}

	if len(specs) >= 1 {
		ts, ok := specs[0].(*ast.TypeSpec)
		if ok {
			td.Name = ts.Name.Name
			rt, ok := ts.Type.(*ast.Ident)
			if ok {
				td.Type = rt.Name
			}
			found = true
		}
	}

	return td, found
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
			field := _extractField(input, imports)
			fields = append(fields, field)
		} else {
			// A single field can refer to multiple: example: x,y int -> x int, y int
			for _, name := range input.Names {
				field := _extractField(input, imports)
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
				elt, ok := arr.Elt.(*ast.StarExpr)
				if ok {
					if ok {
						ident, ok := elt.X.(*ast.Ident)
						if ok {
							field.TypeName = ident.Name
							field.IsPointer = true
						}
					}

					x, ok := elt.X.(*ast.SelectorExpr)
					if ok {
						xx, ok := x.X.(*ast.Ident)
						if ok {
							field.PackageName = imports[xx.Name]
							field.IsPointer = true
							field.TypeName = fmt.Sprintf("%s.%s", xx.Name, x.Sel.Name)
						}
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
					field.TypeName = fmt.Sprintf("%s.%s", ident.Name, sel.Sel.Name)
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
