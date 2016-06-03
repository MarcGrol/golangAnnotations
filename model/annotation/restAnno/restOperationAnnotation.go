package restAnno

import "github.com/MarcGrol/astTools/model/annotation"

const (
	TypeRestOperation = "RestOperation"
	ParamPath         = "Path"
	ParamMethod       = "Method"
)

// make annotation-registry aware of this annotation
func init() {
	annotation.RegisterAnnotation(TypeRestOperation, []string{ParamMethod, ParamPath}, validateRestOperationAnnotation)
}

func validateRestOperationAnnotation(annot annotation.Annotation) bool {
	if annot.Annotation == TypeRestOperation {
		_, hasPath := annot.With[ParamPath]
		_, hasMethod := annot.With[ParamMethod]
		return (hasPath && hasMethod)
	}
	return false
}
