package restAnno

import "github.com/MarcGrol/astTools/model/annotation"

const (
	typeRestOperation = "RestOperation"
	typeRestService   = "RestService"
	ParamPath         = "Path"
	ParamMethod       = "Method"
)

// Register makes the annotation-registry aware of these annotation
func Register() {
	annotation.RegisterAnnotation(typeRestOperation, []string{ParamMethod, ParamPath}, validateRestOperationAnnotation)
	annotation.RegisterAnnotation(typeRestService, []string{ParamPath}, validateRestServiceAnnotation)
}

func validateRestOperationAnnotation(annot annotation.Annotation) bool {
	if annot.Annotation == typeRestOperation {
		path, hasPath := annot.With[ParamPath]
		method, hasMethod := annot.With[ParamMethod]
		return ((hasPath && path != "") && hasMethod && method != "")
	}
	return false
}

func validateRestServiceAnnotation(annot annotation.Annotation) bool {
	if annot.Annotation == typeRestService {
		_, ok := annot.With[ParamPath]
		return ok
	}
	return false
}
