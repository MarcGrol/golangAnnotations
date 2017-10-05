package restAnnotation

import "github.com/MarcGrol/golangAnnotations/annotation"

const (
	TypeRestOperation   = "RestOperation"
	TypeRestService     = "RestService"
	ParamCredentials    = "credentials"
	ParamNoValidation   = "novalidation"
	ParamNoTest         = "notest"
	ParamNoWrap         = "nowrap"
	ParamAfter          = "after"
	ParamPath           = "path"
	ParamMethod         = "method"
	ParamForm           = "form"
	ParamFormat         = "format"
	ParamFilename       = "filename"
	ParamOptional       = "optionalargs"
	ParamRoles          = "roles"
	ParamProducesEvents = "producesevents"
)

// Register makes the annotation-registry aware of these annotation
func Register() {
	annotation.RegisterAnnotation(TypeRestService, []string{ParamCredentials, ParamNoValidation, ParamNoTest, ParamPath}, validateRestServiceAnnotation)
	annotation.RegisterAnnotation(TypeRestOperation, []string{ParamNoWrap, ParamAfter, ParamPath, ParamMethod, ParamForm, ParamFormat, ParamFilename, ParamOptional, ParamRoles, ParamProducesEvents}, validateRestOperationAnnotation)
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
