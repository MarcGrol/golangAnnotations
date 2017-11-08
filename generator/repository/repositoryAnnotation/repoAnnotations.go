package repositoryAnnotation

import "github.com/MarcGrol/golangAnnotations/annotation"

const (
	TypeRepository = "Repository"
	ParamAggregate = "aggregate"
	ParamPackage   = "package"
	ParamModel     = "model"
	ParamMethods   = "methods"
)

// Register makes the annotation-registry aware of this annotation
func Get() []annotation.AnnotationDescriptor {
	return []annotation.AnnotationDescriptor{
		{
			Name:       TypeRepository,
			ParamNames: []string{ParamAggregate, ParamPackage, ParamModel, ParamMethods},
			Validator:  validateRepositoryAnnotation,
		},
	}
}

func validateRepositoryAnnotation(annot annotation.Annotation) bool {
	if annot.Name == TypeRepository {
		aggregate, hasAggregate := annot.Attributes[ParamAggregate]
		if !hasAggregate || aggregate == "" {
			return false
		}
		methods, hasMethods := annot.Attributes[ParamMethods]
		if !hasMethods || methods == "" {
			return false
		}
		return true
	}
	return false
}
