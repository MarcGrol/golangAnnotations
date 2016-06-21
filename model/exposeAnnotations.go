package model

import (
	"strings"

	"github.com/MarcGrol/astTools/model/annotation"
	"github.com/MarcGrol/astTools/model/annotation/eventAnno"
	"github.com/MarcGrol/astTools/model/annotation/restAnno"
)

func init() {
	eventAnno.Register()
	restAnno.Register()
}

func (s Struct) IsEvent() bool {
	_, ok := annotation.ResolveAnnotations(s.DocLines)
	return ok
}

func (s Struct) GetAggregateName() string {
	val, ok := annotation.ResolveAnnotations(s.DocLines)
	if ok {
		return val.With[eventAnno.ParamAggregate]
	}
	return ""
}

func (s Struct) IsRestService() bool {
	_, ok := annotation.ResolveAnnotations(s.DocLines)
	return ok
}

func (o Struct) GetRestServicePath() string {
	val, ok := annotation.ResolveAnnotations(o.DocLines)
	if ok {
		return val.With[restAnno.ParamPath]
	}
	return ""
}

func (o Operation) IsRestOperation() bool {
	_, ok := annotation.ResolveAnnotations(o.DocLines)
	return ok
}

func (o Operation) GetRestOperationPath() string {
	val, ok := annotation.ResolveAnnotations(o.DocLines)
	if ok {
		return val.With[restAnno.ParamPath]
	}
	return ""
}

func (o Operation) GetRestOperationMethod() string {
	val, ok := annotation.ResolveAnnotations(o.DocLines)
	if ok {
		return val.With[restAnno.ParamMethod]
	}
	return ""
}

func (o Operation) HasInput() bool {
	if o.GetRestOperationMethod() == "POST" || o.GetRestOperationMethod() == "PUT" {
		return true
	}
	return false
}

func (o Operation) GetInputArgType() string {
	for _, arg := range o.InputArgs {
		if arg.TypeName != "int" && arg.TypeName != "string" {
			return arg.TypeName
		}
	}
	return ""
}

func (o Operation) GetInputArgName() string {
	for _, arg := range o.InputArgs {
		if arg.TypeName != "int" && arg.TypeName != "string" {
			return arg.Name
		}
	}
	return ""
}

func (o Operation) GetInputParamString() string {
	args := []string{}
	for _, arg := range o.InputArgs {
		args = append(args, arg.Name)
	}
	return strings.Join(args, ",")
}

func (o Operation) HasOutput() bool {
	for _, arg := range o.OutputArgs {
		if arg.TypeName != "error" {
			return true
		}
	}
	return false
}

func (o Operation) GetOutputArgType() string {
	for _, arg := range o.OutputArgs {
		if arg.TypeName != "error" {
			return arg.TypeName
		}
	}
	return ""
}

func (f Field) IsPrimitive() bool {
	return f.TypeName == "int" || f.TypeName == "string"
}
