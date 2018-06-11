package parser

import (
	"fmt"
	"go/ast"
	"log"

	"github.com/MarcGrol/golangAnnotations/model"
)

func extractFieldList(fieldList *ast.FieldList, imports map[string]string) []model.Field {
	mFields := []model.Field{}
	if fieldList != nil {
		for _, field := range fieldList.List {
			mFields = append(mFields, extractFields(field, imports)...)
		}
	}
	return mFields
}

func extractFields(field *ast.Field, imports map[string]string) []model.Field {
	mFields := []model.Field{}
	if field != nil {
		if len(field.Names) == 0 {
			f, ok := extractField(field, imports)
			if ok {
				mFields = append(mFields, f)
			}
		} else {
			// A single field can refer to multiple: example: x,y int -> x int, y int
			for _, name := range field.Names {
				field, ok := extractField(field, imports)
				if ok {
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

	if extractSliceField(field, &mField, imports) {
		return mField, true
	}

	if extractMapField(field, &mField, imports) {
		return mField, true
	}

	if extractPointerField(field, &mField, imports) {
		return mField, true
	}

	if extractIdentField(field, &mField, imports) {
		return mField, true
	}

	if extractSelectorField(field, &mField, imports) {
		return mField, true
	}

	log.Printf("*** Could not understand field %+v: '%+v (%s)'", field, field.Names, field.Type)

	return mField, false
}

func extractSliceField(field *ast.Field, mField *model.Field, imports map[string]string) bool {
	arrayType, ok := field.Type.(*ast.ArrayType)
	if ok {
		mField.IsSlice = true
		if extractSliceSelectorField(arrayType, mField, imports) {
			return true
		}
		if extractSlicePointerField(arrayType, mField, imports) {
			return true
		}
	}
	return false
}

func extractSliceSelectorField(arrayType *ast.ArrayType, mField *model.Field, imports map[string]string) bool {
	ident, ok := arrayType.Elt.(*ast.Ident)
	if ok {
		mField.TypeName = ident.Name
		return true
	}

	selectorExpr, ok := arrayType.Elt.(*ast.SelectorExpr)
	if ok {
		ident, ok = selectorExpr.X.(*ast.Ident)
		if ok {
			mField.TypeName = fmt.Sprintf("%s.%s", ident.Name, selectorExpr.Sel.Name)
			mField.PackageName = imports[ident.Name]
			return true
		}
	}
	return false
}

func extractSlicePointerField(arrayType *ast.ArrayType, mField *model.Field, imports map[string]string) bool {
	starExpr, ok := arrayType.Elt.(*ast.StarExpr)
	if ok {
		if ok {
			ident, ok := starExpr.X.(*ast.Ident)
			if ok {
				mField.TypeName = ident.Name
				mField.IsPointer = true
				return true
			}
		}

		selectorExpr, ok := starExpr.X.(*ast.SelectorExpr)
		if ok {
			ident, ok := selectorExpr.X.(*ast.Ident)
			if ok {
				mField.PackageName = imports[ident.Name]
				mField.IsPointer = true
				mField.TypeName = fmt.Sprintf("%s.%s", ident.Name, selectorExpr.Sel.Name)
				return true
			}
		}
	}
	return false
}

func extractMapField(field *ast.Field, mField *model.Field, imports map[string]string) bool {
	mapKey := ""
	mapValue := ""

	mapType, ok := field.Type.(*ast.MapType)
	if ok {
		key, ok := mapType.Key.(*ast.Ident)
		if ok {
			mapKey = key.Name
		}

		{
			value, ok := mapType.Value.(*ast.Ident)
			if ok {
				mapValue = value.Name
			}
		}
		{
			value, ok := mapType.Value.(*ast.StarExpr)
			if ok {
				ident, ok := value.X.(*ast.Ident)
				if ok {
					mapValue = fmt.Sprintf("*%s", ident.Name)
				}
			}
		}
		{
			value, ok := mapType.Value.(*ast.ArrayType)
			if ok {
				{
					ident, ok := value.Elt.(*ast.Ident)
					if ok {
						mapValue = fmt.Sprintf("[]%s", ident.Name)
					}
				}
				{
					selectorExpr, ok := value.Elt.(*ast.SelectorExpr)
					if ok {
						ident, ok := selectorExpr.X.(*ast.Ident)
						if ok {
							mapValue = fmt.Sprintf("%s.%s", ident.Name, selectorExpr.Sel.Name)
						}
					}
				}
				{
					starExpr, ok := value.Elt.(*ast.StarExpr)
					if ok {
						ident, ok := starExpr.X.(*ast.Ident)
						if ok {
							mapValue = fmt.Sprintf("[]*%s", ident.Name)
						}
					}
				}
			}
		}

	}
	if mapKey != "" && mapValue != "" {
		mField.TypeName = fmt.Sprintf("map[%s]%s", mapKey, mapValue)
		mField.IsMap = true
		return true
	}

	return false
}

func extractPointerField(field *ast.Field, mField *model.Field, imports map[string]string) bool {
	{
		starExpr, ok := field.Type.(*ast.StarExpr)
		if ok {
			ident, ok := starExpr.X.(*ast.Ident)
			if ok {
				mField.TypeName = ident.Name
				mField.IsPointer = true
				return true
			}

			selectorExpr, ok := starExpr.X.(*ast.SelectorExpr)
			if ok {
				ident, ok = selectorExpr.X.(*ast.Ident)
				if ok {
					mField.TypeName = fmt.Sprintf("%s.%s", ident.Name, selectorExpr.Sel.Name)
					mField.IsPointer = true
					mField.PackageName = imports[ident.Name]
					return true
				}
			}
		}
	}
	return false
}

func extractIdentField(field *ast.Field, mField *model.Field, imports map[string]string) bool {
	ident, ok := field.Type.(*ast.Ident)
	if ok {
		mField.TypeName = ident.Name
		return true
	}
	return false
}

func extractSelectorField(field *ast.Field, mField *model.Field, imports map[string]string) bool {
	selectorExpr, ok := field.Type.(*ast.SelectorExpr)
	if ok {
		ident, ok := selectorExpr.X.(*ast.Ident)
		if ok {
			mField.Name = ident.Name
			mField.TypeName = fmt.Sprintf("%s.%s", ident.Name, selectorExpr.Sel.Name)
			mField.PackageName = imports[ident.Name]
			return true
		}
	}
	return false
}
