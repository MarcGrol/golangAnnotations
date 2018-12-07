package model

import (
	"fmt"
	"regexp"
	"strings"
)

var splittableTypeName = regexp.MustCompile(`\*?((\w+)\.)?(\w+)`)

func (f Field) SplitTypeName() (string, string) {
	submatch := splittableTypeName.FindStringSubmatch(f.TypeName)
	if len(submatch) == 2 {
		return "", submatch[1]
	} else if len(submatch) == 4 {
		return submatch[2], submatch[3]
	}
	return "", ""
}

func (f Field) EmptyInstance() string {
	if f.IsPointer() {
		return fmt.Sprintf("&%s{}", f.DereferencedTypeName())
	}
	return fmt.Sprintf("%s{}", f.TypeName)
}

func (f Field) DereferencedTypeName() string {
	return strings.TrimPrefix(f.TypeName, "*")
}

func (f Field) IsPointer() bool {
	return strings.HasPrefix(f.TypeName, "*")
}

func (f Field) IsSlice() bool {
	return strings.HasPrefix(f.TypeName, "[]")
}

func (f Field) IsPrimitive() bool {
	return f.IsBool() || f.IsInt() || f.IsString()
}

func (f Field) IsPrimitiveSlice() bool {
	return f.IsBoolSlice() || f.IsIntSlice() || f.IsStringSlice()
}

const (
	type_bool   = "bool"
	type_int    = "int"
	type_string = "string"
	type_date   = "mydate.MyDate"
)

func (f Field) IsBool() bool {
	return f.DereferencedTypeName() == type_bool
}

func (f Field) IsBoolSlice() bool {
	return f.TypeName == "[]"+type_bool
}

func (f Field) IsInt() bool {
	return f.DereferencedTypeName() == type_int
}

func (f Field) IsIntSlice() bool {
	return f.TypeName == "[]"+type_int
}

func (f Field) IsString() bool {
	return f.DereferencedTypeName() == type_string
}

func (f Field) IsStringSlice() bool {
	return f.TypeName == "[]"+type_string
}

func (f Field) IsDate() bool {
	return f.DereferencedTypeName() == type_date
}

func (f Field) IsDateSlice() bool {
	return f.TypeName == "[]"+type_date
}

func (f Field) IsCustom() bool {
	return !f.IsPrimitive() && !f.IsPrimitiveSlice() && !f.IsDate() && !f.IsDateSlice()
}
