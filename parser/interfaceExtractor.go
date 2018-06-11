package parser

import (
	"go/ast"

	"github.com/MarcGrol/golangAnnotations/model"
)

func extractGenDecForInterface(node ast.Node, imports map[string]string) *model.Interface {
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
