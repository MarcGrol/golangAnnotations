package parser

import (
	"fmt"
	"go/ast"
	"log"
	"reflect"
	"strings"

	"github.com/f0rt/golangAnnotations/model"
)

func extractFieldList(fieldList *ast.FieldList, imports map[string]string, commentMap ast.CommentMap) []model.Field {
	mFields := make([]model.Field, 0)
	if fieldList != nil {
		for _, field := range fieldList.List {
			mFields = append(mFields, extractFields(field, imports, commentMap)...)
		}
	}
	return mFields
}

func extractFields(field *ast.Field, imports map[string]string, commentMap ast.CommentMap) []model.Field {
	mFields := make([]model.Field, 0)
	if field != nil {
		if mField := extractField(field, imports, commentMap); mField != nil {
			if len(field.Names) == 0 {
				mFields = append(mFields, *mField)
			} else {
				// A single field can refer to multiple: example: x,y int -> x int, y int
				for _, name := range field.Names {
					mField.Name = name.Name
					mFields = append(mFields, *mField)
				}
			}
		}
	}
	return mFields
}

func extractField(field *ast.Field, imports map[string]string, commentMap ast.CommentMap) *model.Field {
	if fieldType := processExpression(field.Type, imports); fieldType != nil {
		comments := extractComments(field.Comment)

		// Fallback to comments extracted by CommentMap
		if len(comments) == 0{
			if commentMap[field] != nil && commentMap != nil{
				for _, commentGroup := range commentMap[field]{
					comments = append(comments, extractComments(commentGroup)...)
				}
			}
		}

		return &model.Field{
			PackageName:  fieldType.PackageName,
			DocLines:     extractComments(field.Doc),
			Name:         fieldType.Name,
			TypeName:     fieldType.TypeName,
			Tag:          extractTag(field.Tag),
			CommentLines: comments,
		}
	}
	return nil
}

func processExpression(expr ast.Expr, imports map[string]string) *Expression {

	if mExpr := processEllipsis(expr, imports); mExpr != nil {
		return mExpr
	}
	if mExpr := processArrayType(expr, imports); mExpr != nil {
		return mExpr
	}
	if mExpr := processStarExpr(expr, imports); mExpr != nil {
		return mExpr
	}
	if mExpr := processIdent(expr, imports); mExpr != nil {
		return mExpr
	}
	if mExpr := processSelectorExpr(expr, imports); mExpr != nil {
		return mExpr
	}
	if mExpr := processMapType(expr, imports); mExpr != nil {
		return mExpr
	}
	if mExpr := processFuncType(expr, imports); mExpr != nil {
		return mExpr
	}
	if mExpr := processInterfaceType(expr, imports); mExpr != nil {
		return mExpr
	}

	log.Printf("*** Could not understand expression %+v", reflect.TypeOf(expr))
	return nil
}

func processEllipsis(expr ast.Expr, imports map[string]string) *Expression {
	if ellipsisType, ok := expr.(*ast.Ellipsis); ok {
		mExpr := &Expression{
			TypeName: "...",
		}
		if ellipsisType.Elt != nil {
			if elt := processExpression(ellipsisType.Elt, imports); elt != nil {
				mExpr.PackageName = elt.PackageName
				mExpr.TypeName = fmt.Sprintf("...%s", elt.TypeName)
			}
		}
		return mExpr
	}
	return nil
}

func processArrayType(fieldType ast.Expr, imports map[string]string) *Expression {
	if arrayType, ok := fieldType.(*ast.ArrayType); ok {
		if elt := processExpression(arrayType.Elt, imports); elt != nil {
			typeName := fmt.Sprintf("[]%s", elt.TypeName)
			return &Expression{
				PackageName: elt.PackageName,
				TypeName:    typeName,
			}
		}
	}
	return nil
}

func processStarExpr(fieldType ast.Expr, imports map[string]string) *Expression {
	if starExpr, ok := fieldType.(*ast.StarExpr); ok {
		if x := processExpression(starExpr.X, imports); x != nil {
			typeName := fmt.Sprintf("*%s", x.TypeName)
			return &Expression{
				PackageName: x.PackageName,
				TypeName:    typeName,
			}
		}
	}
	return nil
}

func processIdent(fieldType ast.Expr, imports map[string]string) *Expression {
	if ident, ok := fieldType.(*ast.Ident); ok {
		return &Expression{
			TypeName: ident.Name,
		}
	}
	return nil
}

func processSelectorExpr(fieldType ast.Expr, imports map[string]string) *Expression {
	if selectorExpr, ok := fieldType.(*ast.SelectorExpr); ok {
		if ident, ok := selectorExpr.X.(*ast.Ident); ok {
			typeName := fmt.Sprintf("%s.%s", ident.Name, selectorExpr.Sel.Name)
			return &Expression{
				PackageName: imports[ident.Name],
				TypeName:    typeName,
			}
		}
	}
	return nil
}

func processMapType(fieldType ast.Expr, imports map[string]string) *Expression {
	if mapType, ok := fieldType.(*ast.MapType); ok {
		if key := processExpression(mapType.Key, imports); key != nil {
			if value := processExpression(mapType.Value, imports); value != nil {
				typeName := fmt.Sprintf("map[%s]%s", key.TypeName, value.TypeName)
				return &Expression{
					TypeName: typeName,
				}
			}
		}
	}
	return nil
}

func processFuncType(fieldType ast.Expr, imports map[string]string) *Expression {
	if funcType, ok := fieldType.(*ast.FuncType); ok {
		params := make([]string, 0)
		for _, param := range funcType.Params.List {
			if paramField := extractField(param, imports, nil); paramField != nil {
				formattedParam := paramField.TypeName
				if paramField.Name != "" {
					formattedParam = fmt.Sprintf("%s %s", paramField.Name, paramField.TypeName)
				}
				params = append(params, formattedParam)
			}
		}
		results := make([]string, 0)
		if funcType.Results != nil {
			for _, result := range funcType.Results.List {
				if resultType := processExpression(result.Type, imports); resultType != nil {
					results = append(results, resultType.TypeName)
				}
			}
		}
		typeName := fmt.Sprintf("(%s)%s", strings.Join(params, ","), strings.Join(results, ","))
		return &Expression{
			TypeName: typeName,
		}
	}
	return nil
}

func processInterfaceType(fieldType ast.Expr, imports map[string]string) *Expression {
	if interfaceType, ok := fieldType.(*ast.InterfaceType); ok {
		methods := make([]string, 0)
		for _, method := range extractFieldList(interfaceType.Methods, imports, nil) {
			methods = append(methods, fmt.Sprintf("%s%s", method.Name, method.TypeName))
		}
		typeName := fmt.Sprintf("interface{%s}", strings.Join(methods, ","))
		return &Expression{
			TypeName: typeName,
		}
	}
	return nil
}

type Expression struct {
	PackageName string
	Name        string
	TypeName    string
}
