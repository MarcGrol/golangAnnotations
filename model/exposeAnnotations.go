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
	annotation, ok := annotation.ResolveAnnotations(s.DocLines)
	if !ok || annotation.Name != "Event" {
		return false
	}
	return ok
}

func (s Struct) GetAggregateName() string {
	val, ok := annotation.ResolveAnnotations(s.DocLines)
	if ok {
		return val.Attributes[eventAnno.ParamAggregate]
	}
	return ""
}

func (s Struct) IsRestService() bool {
	annotation, ok := annotation.ResolveAnnotations(s.DocLines)
	if !ok || annotation.Name != "RestService" {
		return false
	}
	return ok
}

func (o Struct) GetRestServicePath() string {
	val, ok := annotation.ResolveAnnotations(o.DocLines)
	if ok {
		return val.Attributes[restAnno.ParamPath]
	}
	return ""
}

func (o Operation) IsRestOperation() bool {
	annotation, ok := annotation.ResolveAnnotations(o.DocLines)
	if !ok || annotation.Name != "RestOperation" {
		return false
	}
	return ok
}

func (o Operation) GetRestOperationPath() string {
	val, ok := annotation.ResolveAnnotations(o.DocLines)
	if ok {
		return val.Attributes[restAnno.ParamPath]
	}
	return ""
}

func (o Operation) GetRestOperationMethod() string {
	val, ok := annotation.ResolveAnnotations(o.DocLines)
	if ok {
		return val.Attributes[restAnno.ParamMethod]
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

func (f Field) IsNumber() bool {
	return f.TypeName == "int"
}
