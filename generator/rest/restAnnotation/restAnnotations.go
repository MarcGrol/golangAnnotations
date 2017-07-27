package restAnnotation

import "github.com/MarcGrol/golangAnnotations/annotation"

const (
	TypeRestOperation = "RestOperation"
	TypeRestService   = "RestService"
	ParamNoValidation = "novalidation"
	ParamNoTest       = "notest"
	ParamNoWrap       = "nowrap"
	ParamAfter        = "after"
	ParamPath         = "path"
	ParamMethod       = "method"
	ParamForm         = "form"
	ParamFormat       = "format"
	ParamFilename     = "filename"
	ParamOptional     = "optionalargs"
	ParamRoles        = "roles"
)

// Register makes the annotation-registry aware of these annotation
func Register() {
	annotation.RegisterAnnotation(TypeRestOperation, []string{ParamNoWrap, ParamAfter, ParamPath, ParamMethod, ParamForm, ParamFormat, ParamFilename, ParamOptional, ParamRoles}, validateRestOperationAnnotation)
	annotation.RegisterAnnotation(TypeRestService, []string{ParamNoValidation, ParamNoTest, ParamPath}, validateRestServiceAnnotation)
}

func validateRestOperationAnnotation(annot annotation.Annotation) bool {
	if annot.Name == TypeRestOperation {
		method, hasMethod := annot.Attributes[ParamMethod]
		return hasMethod && method != ""
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
