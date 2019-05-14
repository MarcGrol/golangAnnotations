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
	//bodyparams := getBodyParams(o)

	if len(params) > 0 {
		lines = append(lines, "parameters:")
		lines = append(lines, params...)
	}

	responses := getResponses(o)

	// responses:
	//   '200':
	//     description: {{$oper.Name}} response
	// 	   schema:
	// 	     type: object
	//       items:
	// 		   "$ref": "#/definitions/{{GetOutputArgType .}}"

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

func getResponses(o model.Operation) []string {
	lines := []string{}

	lines = append(lines, "responses:")
	lines = append(lines, "  '200':")
	lines = append(lines, fmt.Sprintf("    description: %s response", o.Name))
	lines = append(lines, "    schema:")
	lines = append(lines, "      type: object")
	lines = append(lines, "      properties:")
	lines = append(lines, "")
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
