package model

import (
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
