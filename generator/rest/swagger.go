package rest

import (
	"fmt"
	"github.com/MarcGrol/golangAnnotations/model"
	"strings"
)

func GetSwagger2(o model.Operation, s model.Struct) string {
	lines := []string{}

	swaggerRoute := fmt.Sprintf("swagger:operation %s %s%s %s", GetRestOperationMethod(o), GetRestServicePath(s), GetRestOperationPath(o), o.Name)
	lines = append(lines, swaggerRoute)
	lines = append(lines, "")
	lines = append(lines, "Some description")
	lines = append(lines, "")
	lines = append(lines, "---")

	params := getParams(o)
	//headerParams := getHeaderParams(o)
	bodyParams := getBodyParams(o)

	if len(params) > 0 || len(bodyParams) > 0 {
		lines = append(lines, "parameters:")
		if len(params) > 0 {
			lines = append(lines, params...)
		}
		if len(bodyParams) > 0 {
			lines = append(lines, bodyParams...)
		}
	}

	responses := getResponses(o)
	errorResponses := getErrorResponses(o)

	if len(responses) > 0 || len(errorResponses) > 0 {
		lines = append(lines, "responses:")
		if len(responses) > 0 {
			lines = append(lines, "  '200':")
			lines = append(lines, responses...)
		}
		if len(errorResponses) > 0 {
			lines = append(lines, "  '400':")
			lines = append(lines, errorResponses...)
		}
	}

	for i := range lines {
		lines[i] = fmt.Sprintf("// %s", lines[i])
	}
	return strings.Join(lines, "\n")
}

func getParams(o model.Operation) []string {
	lines := []string{}

	for _, arg := range o.InputArgs {
		if !IsCustomArg(arg) {
			if strings.Contains(GetRestOperationPath(o), fmt.Sprintf("{%s}", Uncapitalized(arg.Name))) {
				lines = append(lines, "  - in: path")
			} else {
				lines = append(lines, "  - in: query")
			}
			lines = append(lines, fmt.Sprintf("    name: %s", Uncapitalized(arg.Name)))
			lines = append(lines, fmt.Sprintf("    type: %s", getSwaggerType(arg.TypeName)))
			lines = append(lines, fmt.Sprintf("    required: %t", IsInputArgMandatory(o, arg)))
		}
	}
	return lines
}

func getBodyParams(o model.Operation) []string {
	lines := []string{}
	if HasInput(o) {
		lines = append(lines, "  - in: body")
		lines = append(lines, fmt.Sprintf("    name: %s", GetInputArgType(o)))
		lines = append(lines, fmt.Sprintf("    description: %s description", GetInputArgType(o)))
		lines = append(lines, "    schema:")
		lines = append(lines, fmt.Sprintf(`      "$ref": "#/definitions/%s"`, GetInputArgType(o)))
	}
	return lines
}

func getResponses(o model.Operation) []string {
	lines := []string{}

	for _, arg := range o.OutputArgs {
		if !IsErrorArg(arg) && !IsMetaCallbackArg(arg) {
			lines = append(lines, fmt.Sprintf("    description: %s response", o.Name))
			lines = append(lines, "    schema:")
			lines = append(lines, fmt.Sprintf(`      "$ref": "#/definitions/%s"`, arg.DereferencedTypeName()))
		}
	}
	return lines
}

func getErrorResponses(o model.Operation) []string {
	lines := []string{}

	for _, arg := range o.OutputArgs {
		if IsErrorArg(arg) {
			lines = append(lines, fmt.Sprintf("    description: %s response", o.Name))
			lines = append(lines, "    schema:")
			lines = append(lines, `      "$ref": "#/definitions/Error"`)
		}
	}
	return lines
}

func getSwaggerType(gotype string) string {
	typeMap := map[string]string{
		"[]string": "array",
		"int": "integer",
		"string": "string",
		"bool": "boolean",
	}
	return typeMap[gotype]
}
