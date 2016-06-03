package model

import "github.com/MarcGrol/astTools/model/annotation"

// make annotation info available to template of generator
func (s Struct) IsEvent() bool {
	_, ok := annotation.ResolveAnnotations(s.DocLines)
	return ok
}

func (s Struct) GetAggregateName() string {
	val, ok := annotation.ResolveAnnotations(s.DocLines)
	if ok {
		return val.With["Aggregate"]
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
		return val.With["Path"]
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
		return val.With["Path"]
	}
	return ""
}

func (o Operation) GetRestOperationMethod() string {
	val, ok := annotation.ResolveAnnotations(o.DocLines)
	if ok {
		return val.With["Method"]
	}
	return ""
}
