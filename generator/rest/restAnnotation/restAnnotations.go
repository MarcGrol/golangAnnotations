package restAnnotation

import "github.com/MarcGrol/golangAnnotations/annotation"

const (
	TypeRestOperation = "RestOperation"
	TypeRestService   = "RestService"
	ParamPath         = "path"
	ParamMethod       = "method"
	ParamOptional     = "optionalargs"
)

// Register makes the annotation-registry aware of these annotation
func Register() {
	annotation.RegisterAnnotation(TypeRestOperation, []string{ParamMethod, ParamPath, ParamOptional}, validateRestOperationAnnotation)
	annotation.RegisterAnnotation(TypeRestService, []string{ParamPath}, validateRestServiceAnnotation)
}

func validateRestOperationAnnotation(annot annotation.Annotation) bool {
	if annot.Name == TypeRestOperation {
		method, hasMethod := annot.Attributes[ParamMethod]
		return (hasMethod && method != "")
	}
	return false
}

func validateRestServiceAnnotation(annot annotation.Annotation) bool {
	if annot.Name == TypeRestService {
		_, ok := annot.Attributes[ParamPath]
		return ok
	}
	return false
}
