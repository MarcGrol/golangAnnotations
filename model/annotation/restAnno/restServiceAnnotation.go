package restAnno

import "github.com/MarcGrol/astTools/model/annotation"

const (
	TypeRestService = "RestService"
)

// make annotation-registry aware of this annotation
func init() {
	annotation.RegisterAnnotation(TypeRestService, []string{ParamPath}, validateRestServiceAnnotation)
}

func validateRestServiceAnnotation(annot annotation.Annotation) bool {
	if annot.Annotation == TypeRestService {
		_, ok := annot.With[ParamPath]
		return ok
	}
	return false
}
