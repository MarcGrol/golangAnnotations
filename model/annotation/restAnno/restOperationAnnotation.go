package restAnno

import "github.com/MarcGrol/astTools/model/annotation"

const (
	typeRestOperation = "RestOperation"
	typeRestService   = "RestService"
	ParamPath         = "path"
	ParamMethod       = "method"
)

// Register makes the annotation-registry aware of these annotation
func Register() {
	annotation.RegisterAnnotation(typeRestOperation, []string{ParamMethod, ParamPath}, validateRestOperationAnnotation)
	annotation.RegisterAnnotation(typeRestService, []string{ParamPath}, validateRestServiceAnnotation)
}

func validateRestOperationAnnotation(annot annotation.Annotation) bool {
	if annot.Name == typeRestOperation {
		path, hasPath := annot.Attributes[ParamPath]
		method, hasMethod := annot.Attributes[ParamMethod]
		return ((hasPath && path != "") && hasMethod && method != "")
	}
	return false
}

func validateRestServiceAnnotation(annot annotation.Annotation) bool {
	if annot.Name == typeRestService {
		_, ok := annot.Attributes[ParamPath]
		return ok
	}
	return false
}
