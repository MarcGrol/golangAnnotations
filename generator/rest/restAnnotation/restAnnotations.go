package restAnnotation

import "github.com/MarcGrol/golangAnnotations/annotation"

const (
	TypeRestOperation   = "RestOperation"
	TypeRestService     = "RestService"
	ParamCredentials    = "credentials"
	ParamNoValidation   = "novalidation"
	ParamProtected      = "protected"
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

func Get() []annotation.AnnotationDescriptor {
	return []annotation.AnnotationDescriptor{
		{
			Name:       TypeRestService,
			ParamNames: []string{ParamCredentials, ParamNoValidation, ParamProtected, ParamNoTest, ParamPath},
			Validator:  validateRestServiceAnnotation,
		},
		{
			Name:       TypeRestOperation,
			ParamNames: []string{ParamNoWrap, ParamAfter, ParamPath, ParamMethod, ParamForm, ParamFormat, ParamFilename, ParamOptional, ParamRoles, ParamProducesEvents},
			Validator:  validateRestOperationAnnotation,
		}}
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
