package parser

import (
	"fmt"
	"go/ast"
	"log"
	"reflect"
	"strings"

	"github.com/MarcGrol/golangAnnotations/model"
)

func extractFieldList(fieldList *ast.FieldList, imports map[string]string) []model.Field {
	mFields := make([]model.Field, 0)
	if fieldList != nil {
		for _, field := range fieldList.List {
			mFields = append(mFields, extractFields(field, imports)...)
		}
	}
	return mFields
}

func extractFields(field *ast.Field, imports map[string]string) []model.Field {
	mFields := make([]model.Field, 0)
	if field != nil {
		if len(field.Names) == 0 {
			if f, ok := extractField(field, imports); ok {
				mFields = append(mFields, f)
			}
		} else {
			// A single field can refer to multiple: example: x,y int -> x int, y int
			for _, name := range field.Names {
				if field, ok := extractField(field, imports); ok {
					field.Name = name.Name
					mFields = append(mFields, field)
				}
			}
		}
	}
	return mFields
}

func extractField(field *ast.Field, imports map[string]string) (model.Field, bool) {
	mField := model.Field{
		DocLines:     extractComments(field.Doc),
		CommentLines: extractComments(field.Comment),
		Tag:          extractTag(field.Tag),
	}

	if extractEllipsisField(field.Type, &mField, imports) {
		return mField, true
	}

	if extractSliceField(field.Type, &mField, imports) {
		return mField, true
	}

	if extractElementType(field.Type, &mField, imports) {
		return mField, true
	}

	if extractMapField(field.Type, &mField, imports) {
		return mField, true
	}

	if extractFuncTypeField(field.Type, &mField, imports) {
		return mField, true
	}

	if extractInterfaceField(field.Type, &mField, imports) {
		return mField, true
	}

	log.Printf("*** Could not understand field %+v: '%+v (%s)'", field, field.Names, field.Type)

	return mField, false
}

func extractEllipsisField(expr ast.Expr, mField *model.Field, imports map[string]string) bool {
	if ellipsisType, ok := expr.(*ast.Ellipsis); ok {

		mField.IsEllipsis = true

		if extractElementType(ellipsisType.Elt, mField, imports) {
			return true
		}
	}
	return false
}

func extractSliceField(fieldType ast.Expr, mField *model.Field, imports map[string]string) bool {
	if arrayType, ok := fieldType.(*ast.ArrayType); ok {

		mField.IsSlice = true

		if extractElementType(arrayType.Elt, mField, imports) {
			return true
		}
	}
	return false
}

func extractElementType(elt ast.Expr, mField *model.Field, imports map[string]string) bool {

	if extractPointerField(elt, mField, imports) {
		return true
	}

	if extractIdentField(elt, mField, imports) {
		return true
	}

	if extractSelectorField(elt, mField, imports) {
		return true
	}

	return false
}

func extractPointerField(fieldType ast.Expr, mField *model.Field, imports map[string]string) bool {
	if starExpr, ok := fieldType.(*ast.StarExpr); ok {

		mField.IsPointer = true

		if extractIdentField(starExpr.X, mField, imports) {
			return true
		}

		if extractSelectorField(starExpr.X, mField, imports); ok {
			return true
		}
	}
	return false
}

func extractIdentField(fieldType ast.Expr, mField *model.Field, imports map[string]string) bool {
	if ident, ok := fieldType.(*ast.Ident); ok {
		mField.TypeName = ident.Name
		return true
	}
	return false
}

func extractSelectorField(fieldType ast.Expr, mField *model.Field, imports map[string]string) bool {
	if selectorExpr, ok := fieldType.(*ast.SelectorExpr); ok {
		if ident, ok := selectorExpr.X.(*ast.Ident); ok {
			mField.PackageName = imports[ident.Name]
			mField.TypeName = formatExpression(selectorExpr, imports)
			return true
		}
	}
	return false
}

func extractMapField(fieldType ast.Expr, mField *model.Field, imports map[string]string) bool {
	if mapType, ok := fieldType.(*ast.MapType); ok {

		mField.IsMap = true

		if mapKey := formatExpression(mapType.Key, imports); mapKey != "" {
			if mapValue := formatExpression(mapType.Value, imports); mapValue != "" {
				mField.TypeName = fmt.Sprintf("map[%s]%s", mapKey, mapValue)
				return true
			}
		}
	}
	return false
}

func extractFuncTypeField(fieldType ast.Expr, mField *model.Field, imports map[string]string) bool {
	if funcType, ok := fieldType.(*ast.FuncType); ok {
		params := make([]string, 0)
		for _, param := range funcType.Params.List {
			if paramField, ok := extractField(param, imports); ok {
				formattedParam := paramField.TypeName
				if paramField.Name != "" {
					formattedParam = fmt.Sprintf("%s %s", paramField.Name, paramField.TypeName)
				}
				params = append(params, formattedParam)
			} else {
				log.Printf("Skipping unrecognized funcType.Param: %+v\n", param)
			}
		}
		results := make([]string, 0)
		if funcType.Results != nil {
			for _, result := range funcType.Results.List {
				formattedResult := formatExpression(result.Type, imports)
				if formattedResult != "" {
					results = append(results, formattedResult)
				} else {
					log.Printf("Skipping unrecognized functType.Result: %+v\n", result)
				}
			}
		}
		mField.TypeName = fmt.Sprintf("(%s)%s", strings.Join(params, ","), strings.Join(results, ","))
		return true
	}
	return false
}

func extractInterfaceField(fieldType ast.Expr, mField *model.Field, imports map[string]string) bool {
	if interfaceType, ok := fieldType.(*ast.InterfaceType); ok {
		methods := make([]string, 0)
		for _, method := range extractFieldList(interfaceType.Methods, imports) {
			methods = append(methods, fmt.Sprintf("%s%s", method.Name, method.TypeName))
		}
		mField.TypeName = fmt.Sprintf("interface{%s}", strings.Join(methods, ","))
		return true
	}
	return false
}

func formatExpression(expr ast.Expr, imports map[string]string) string {

	if arrayType, ok := expr.(*ast.ArrayType); ok {
		if arrayElt := formatExpression(arrayType.Elt, imports); arrayElt != "" {
			return fmt.Sprintf("[]%s", arrayElt)
		}
	}

	if starExpr, ok := expr.(*ast.StarExpr); ok {
		if starX := formatExpression(starExpr.X, imports); starX != "" {
			return fmt.Sprintf("*%s", starX)
		}
	}

	if ident, ok := expr.(*ast.Ident); ok {
		return fmt.Sprintf("%s", ident.Name)
	}

	if selectorExpr, ok := expr.(*ast.SelectorExpr); ok {
		if ident, ok := selectorExpr.X.(*ast.Ident); ok {
			return fmt.Sprintf("%s.%s", ident.Name, selectorExpr.Sel.Name)
		}
	}

	if mapType, ok := expr.(*ast.MapType); ok {
		if mapKey := formatExpression(mapType.Key, imports); mapKey != "" {
			if mapValue := formatExpression(mapType.Value, imports); mapValue != "" {
				return fmt.Sprintf("map[%s]%s", mapKey, mapValue)
			}
		}
	}

	if funcType, ok := expr.(*ast.FuncType); ok {
		params := make([]string, 0)
		for _, param := range funcType.Params.List {
			if formattedParam := formatExpression(param.Type, imports); formattedParam != "" {
				params = append(params, formattedParam)
			} else {
				log.Printf("Skipping unrecognized funcType.Param: %+v\n", param)
			}
		}
		results := make([]string, 0)
		if funcType.Results != nil {
			for _, result := range funcType.Results.List {
				formattedResult := formatExpression(result.Type, imports)
				if formattedResult != "" {
					results = append(results, formattedResult)
				} else {
					log.Printf("Skipping unrecognized functType.Result: %+v\n", result)
				}
			}
		}
		return fmt.Sprintf("(%s)%s", strings.Join(params, ","), strings.Join(results, ","))
	}

	if interfaceType, ok := expr.(*ast.InterfaceType); ok {
		methods := make([]string, 0)
		for _, method := range extractFieldList(interfaceType.Methods, imports) {
			methods = append(methods, fmt.Sprintf("%s%s", method.Name, method.TypeName))
		}
		return fmt.Sprintf("interface{%s}", strings.Join(methods, ","))
	}

	log.Printf("Unrecognized expression: %+v\n", reflect.TypeOf(expr))

	return ""
}
