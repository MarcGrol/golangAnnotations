package parser

import (
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

var debugAstOfSources = false

type myParser struct {
}

func New() Parser {
	return &myParser{}
}

func (p *myParser) ParseSourceDir(dirName string, includeRegex string, excludeRegex string) (model.ParsedSources, error) {
	if debugAstOfSources {
		dumpFilesInDir(dirName)
	}
	packages, err := parseDir(dirName, includeRegex, excludeRegex)
	if err != nil {
		log.Printf("error parsing dir %s: %s", dirName, err.Error())
		return model.ParsedSources{}, err
	}

	v := &astVisitor{
		Imports: map[string]string{},
	}
	for _, aPackage := range packages {
		parsePackage(aPackage, v)
	}

	embedOperationsInStructs(v)

	embedTypedefDocLinesInEnum(v)

	return model.ParsedSources{
		Structs:    v.Structs,
		Operations: v.Operations,
		Interfaces: v.Interfaces,
		Typedefs:   v.Typedefs,
		Enums:      v.Enums,
	}, nil
}

func parsePackage(aPackage *ast.Package, v *astVisitor) {
	for _, fileEntry := range sortedFileEntries(aPackage.Files) {
		v.CurrentFilename = fileEntry.key

		appEngineOnly := true
		for _, commentGroup := range fileEntry.file.Comments {
			if commentGroup != nil {
				for _, comment := range commentGroup.List {
					if comment != nil && comment.Text == "// +build !appengine" {
						appEngineOnly = false
					}
				}
			}
		}
		if appEngineOnly {
			ast.Walk(v, &fileEntry.file)
		}
	}
}

func parseSourceFile(srcFilename string) (model.ParsedSources, error) {
	if debugAstOfSources {
		dumpFile(srcFilename)
	}

	v, err := doParseFile(srcFilename)
	if err != nil {
		log.Printf("error parsing src %s: %s", srcFilename, err.Error())
		return model.ParsedSources{}, err
	}

	embedOperationsInStructs(v)
	embedTypedefDocLinesInEnum(v)

	return model.ParsedSources{
		Structs:    v.Structs,
		Operations: v.Operations,
		Interfaces: v.Interfaces,
		Typedefs:   v.Typedefs,
		Enums:      v.Enums,
	}, nil
}

func doParseFile(srcFilename string) (*astVisitor, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, srcFilename, nil, parser.ParseComments)
	if err != nil {
		log.Printf("error parsing src-file %s: %s", srcFilename, err.Error())
		return nil, err
	}
	v := &astVisitor{
		Imports: map[string]string{},
	}
	v.CurrentFilename = srcFilename
	ast.Walk(v, file)

	return v, nil
}

type fileEntry struct {
	key  string
	file ast.File
}

type fileEntries []fileEntry

func (list fileEntries) Len() int {
	return len(list)
}

func (list fileEntries) Less(i, j int) bool {
	return list[i].key < list[j].key
}

func (list fileEntries) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func sortedFileEntries(fileMap map[string]*ast.File) fileEntries {
	var fileEntries fileEntries = make([]fileEntry, 0, len(fileMap))
	for key, file := range fileMap {
		if file != nil {
			fileEntries = append(fileEntries, fileEntry{
				key:  key,
				file: *file,
			})
		}
	}
	sort.Sort(fileEntries)
	return fileEntries
}

func embedOperationsInStructs(visitor *astVisitor) {
	mStructMap := make(map[string]*model.Struct)
	for idx := range visitor.Structs {
		mStructMap[(&visitor.Structs[idx]).Name] = &visitor.Structs[idx]
	}
	for idx := range visitor.Operations {
		mOperation := visitor.Operations[idx]
		if mOperation.RelatedStruct != nil {
			if mStruct, ok := mStructMap[mOperation.RelatedStruct.DereferencedTypeName()]; ok {
				mStruct.Operations = append(mStruct.Operations, &mOperation)
			}
		}
	}

}

func embedTypedefDocLinesInEnum(visitor *astVisitor) {
	for idx, mEnum := range visitor.Enums {
		for _, typedef := range visitor.Typedefs {
			if typedef.Name == mEnum.Name {
				visitor.Enums[idx].DocLines = typedef.DocLines
				break
			}
		}
	}
}

func parseDir(dirName string, includeRegex string, excludeRegex string) (map[string]*ast.Package, error) {
	var includePattern = regexp.MustCompile(includeRegex)
	var excludePattern = regexp.MustCompile(excludeRegex)

	fileSet := token.NewFileSet()
	packageMap, err := parser.ParseDir(fileSet, dirName, func(fi os.FileInfo) bool {
		if excludePattern.MatchString(fi.Name()) {
			return false
		}
		return includePattern.MatchString(fi.Name())
	}, parser.ParseComments)
	if err != nil {
		log.Printf("error parsing dir %s: %s", dirName, err.Error())
		return packageMap, err
	}

	return packageMap, nil
}

func dumpFile(srcFilename string) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, srcFilename, nil, parser.ParseComments)
	if err != nil {
		log.Printf("error parsing src %s: %s", srcFilename, err.Error())
		return
	}
	ast.Print(fileSet, file)
}

func dumpFilesInDir(dirName string) {
	fileSet := token.NewFileSet()
	packageMap, err := parser.ParseDir(
		fileSet,
		dirName,
		nil,
		parser.ParseComments)
	if err != nil {
		log.Printf("error parsing dir %s: %s", dirName, err.Error())
	}
	for _, aPackage := range packageMap {
		for _, file := range aPackage.Files {
			ast.Print(fileSet, file)
		}
	}
}

// =====================================================================================================================

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
		if packageName, ok := extractPackageName(node); ok {
			v.PackageName = packageName
		}

		// extract all imports into a map
		v.extractGenDeclImports(node)

		v.parseAsStruct(node)
		v.parseAsTypedef(node)
		v.parseAsEnum(node)
		v.parseAsInterFace(node)
		v.parseAsOperation(node)

	}
	return v
}

func (v *astVisitor) extractGenDeclImports(node ast.Node) {
	if genDecl, ok := node.(*ast.GenDecl); ok {
		for _, spec := range genDecl.Specs {
			if importSpec, ok := spec.(*ast.ImportSpec); ok {
				quotedImport := importSpec.Path.Value
				unquotedImport := strings.Trim(quotedImport, "\"")
				init, last := filepath.Split(unquotedImport)
				if init == "" {
					last = init
				}
				v.Imports[last] = unquotedImport
			}
		}
	}
}

func (v *astVisitor) parseAsStruct(node ast.Node) {
	if mStruct := extractGenDeclForStruct(node, v.Imports); mStruct != nil {
		mStruct.PackageName = v.PackageName
		mStruct.Filename = v.CurrentFilename
		v.Structs = append(v.Structs, *mStruct)
	}
}

func (v *astVisitor) parseAsTypedef(node ast.Node) {
	if mTypedef := extractGenDeclForTypedef(node); mTypedef != nil {
		mTypedef.PackageName = v.PackageName
		mTypedef.Filename = v.CurrentFilename
		v.Typedefs = append(v.Typedefs, *mTypedef)
	}
}

func (v *astVisitor) parseAsEnum(node ast.Node) {
	if mEnum := extractGenDeclForEnum(node); mEnum != nil {
		mEnum.PackageName = v.PackageName
		mEnum.Filename = v.CurrentFilename
		v.Enums = append(v.Enums, *mEnum)
	}
}

func (v *astVisitor) parseAsInterFace(node ast.Node) {
	// if interfaces, get its methods
	if mInterface := extractInterface(node, v.Imports); mInterface != nil {
		mInterface.PackageName = v.PackageName
		mInterface.Filename = v.CurrentFilename
		v.Interfaces = append(v.Interfaces, *mInterface)
	}
}

func (v *astVisitor) parseAsOperation(node ast.Node) {
	// if mOperation, get its signature
	if mOperation := extractOperation(node, v.Imports); mOperation != nil {
		mOperation.PackageName = v.PackageName
		mOperation.Filename = v.CurrentFilename
		v.Operations = append(v.Operations, *mOperation)
	}
}

// =====================================================================================================================

func extractPackageName(node ast.Node) (string, bool) {
	if file, ok := node.(*ast.File); ok {
		if file.Name != nil {
			return file.Name.Name, true
		}
		return "", true
	}
	return "", false
}

// ------------------------------------------------------ STRUCT -------------------------------------------------------

func extractGenDeclForStruct(node ast.Node, imports map[string]string) *model.Struct {
	if genDecl, ok := node.(*ast.GenDecl); ok {
		// Continue parsing to see if it is a struct
		if mStruct := extractSpecsForStruct(genDecl.Specs, imports); mStruct != nil {
			// Docline of struct (that could contain annotations) appear far before the details of the struct
			mStruct.DocLines = extractComments(genDecl.Doc)
			return mStruct
		}
	}
	return nil
}

func extractSpecsForStruct(specs []ast.Spec, imports map[string]string) *model.Struct {
	if len(specs) >= 1 {
		if typeSpec, ok := specs[0].(*ast.TypeSpec); ok {
			if structType, ok := typeSpec.Type.(*ast.StructType); ok {
				return &model.Struct{
					Name:   typeSpec.Name.Name,
					Fields: extractFieldList(structType.Fields, imports),
				}
			}
		}
	}
	return nil
}

// ------------------------------------------------------ TYPEDEF ------------------------------------------------------

func extractGenDeclForTypedef(node ast.Node) *model.Typedef {
	if genDecl, ok := node.(*ast.GenDecl); ok {
		// Continue parsing to see if it a struct
		if mTypedef := extractSpecsForTypedef(genDecl.Specs); mTypedef != nil {
			mTypedef.DocLines = extractComments(genDecl.Doc)
			return mTypedef
		}
	}
	return nil
}

func extractSpecsForTypedef(specs []ast.Spec) *model.Typedef {
	if len(specs) >= 1 {
		if typeSpec, ok := specs[0].(*ast.TypeSpec); ok {
			mTypedef := model.Typedef{
				Name: typeSpec.Name.Name,
			}
			if ident, ok := typeSpec.Type.(*ast.Ident); ok {
				mTypedef.Type = ident.Name
			}
			return &mTypedef
		}
	}
	return nil
}

// ------------------------------------------------------- ENUM --------------------------------------------------------

func extractGenDeclForEnum(node ast.Node) *model.Enum {
	if genDecl, ok := node.(*ast.GenDecl); ok {
		// Continue parsing to see if it is an enum
		// Docs live in the related typedef
		return extractSpecsForEnum(genDecl.Specs)
	}
	return nil
}

func extractSpecsForEnum(specs []ast.Spec) *model.Enum {
	if typeName, ok := extractEnumTypeName(specs); ok {
		mEnum := model.Enum{
			Name:         typeName,
			EnumLiterals: []model.EnumLiteral{},
		}
		for _, spec := range specs {
			if valueSpec, ok := spec.(*ast.ValueSpec); ok {
				enumLiteral := model.EnumLiteral{
					Name: valueSpec.Names[0].Name,
				}
				for _, value := range valueSpec.Values {
					if basicLit, ok := value.(*ast.BasicLit); ok {
						enumLiteral.Value = strings.Trim(basicLit.Value, "\"")
						break
					}
				}
				mEnum.EnumLiterals = append(mEnum.EnumLiterals, enumLiteral)
			}
		}
		return &mEnum
	}
	return nil
}

func extractEnumTypeName(specs []ast.Spec) (string, bool) {
	for _, spec := range specs {
		if valueSpec, ok := spec.(*ast.ValueSpec); ok {
			if valueSpec.Type != nil {
				for _, name := range valueSpec.Names {
					if ident, ok := valueSpec.Type.(*ast.Ident); ok {
						if name.Obj.Kind == ast.Con {
							return ident.Name, true
						}
					}
				}
			}
		}
	}
	return "", false
}

// ----------------------------------------------------- INTERFACE -----------------------------------------------------

func extractInterface(node ast.Node, imports map[string]string) *model.Interface {
	if genDecl, ok := node.(*ast.GenDecl); ok {
		// Continue parsing to see if it an interface
		if mInterface := extractSpecsForInterface(genDecl.Specs, imports); mInterface != nil {
			// Docline of interface (that could contain annotations) appear far before the details of the struct
			mInterface.DocLines = extractComments(genDecl.Doc)
			return mInterface
		}
	}
	return nil
}

func extractSpecsForInterface(specs []ast.Spec, imports map[string]string) *model.Interface {
	if len(specs) >= 1 {
		if typeSpec, ok := specs[0].(*ast.TypeSpec); ok {
			if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
				return &model.Interface{
					Name:    typeSpec.Name.Name,
					Methods: extractInterfaceMethods(interfaceType.Methods, imports),
				}
			}
		}
	}
	return nil
}

func extractInterfaceMethods(fieldList *ast.FieldList, imports map[string]string) []model.Operation {
	methods := make([]model.Operation, 0)
	for _, field := range fieldList.List {
		if len(field.Names) > 0 {
			if funcType, ok := field.Type.(*ast.FuncType); ok {
				methods = append(methods, model.Operation{
					DocLines:   extractComments(field.Doc),
					Name:       field.Names[0].Name,
					InputArgs:  extractFieldList(funcType.Params, imports),
					OutputArgs: extractFieldList(funcType.Results, imports),
				})
			}
		}
	}
	return methods
}

// ----------------------------------------------------- OPERATION -----------------------------------------------------

func extractOperation(node ast.Node, imports map[string]string) *model.Operation {
	if funcDecl, ok := node.(*ast.FuncDecl); ok {
		mOperation := model.Operation{
			DocLines: extractComments(funcDecl.Doc),
		}

		if funcDecl.Recv != nil {
			fields := extractFieldList(funcDecl.Recv, imports)
			if len(fields) >= 1 {
				mOperation.RelatedStruct = &(fields[0])
			}
		}

		if funcDecl.Name != nil {
			mOperation.Name = funcDecl.Name.Name
		}

		if funcDecl.Type.Params != nil {
			mOperation.InputArgs = extractFieldList(funcDecl.Type.Params, imports)
		}

		if funcDecl.Type.Results != nil {
			mOperation.OutputArgs = extractFieldList(funcDecl.Type.Results, imports)
		}
		return &mOperation
	}
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
