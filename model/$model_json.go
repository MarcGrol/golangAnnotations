// Generated automatically by golangAnnotations: do not edit manually

package model

import "encoding/json"

// Helpers for json-struct ParsedSources

// MarshalJSON prevents nil slices in json
func (data ParsedSources) MarshalJSON() ([]byte, error) {
	type alias ParsedSources
	var raw = alias(data)

	if raw.Structs == nil {
		raw.Structs = []Struct{}
	}

	if raw.Operations == nil {
		raw.Operations = []Operation{}
	}

	if raw.Interfaces == nil {
		raw.Interfaces = []Interface{}
	}

	if raw.Typedefs == nil {
		raw.Typedefs = []Typedef{}
	}

	if raw.Enums == nil {
		raw.Enums = []Enum{}
	}

	return json.Marshal(raw)
}

// UnmarshalJSON prevents nil slices from json
func (data *ParsedSources) UnmarshalJSON(b []byte) error {
	type alias ParsedSources
	var raw alias
	err := json.Unmarshal(b, &raw)

	if raw.Structs == nil {
		raw.Structs = []Struct{}
	}

	if raw.Operations == nil {
		raw.Operations = []Operation{}
	}

	if raw.Interfaces == nil {
		raw.Interfaces = []Interface{}
	}

	if raw.Typedefs == nil {
		raw.Typedefs = []Typedef{}
	}

	if raw.Enums == nil {
		raw.Enums = []Enum{}
	}

	*data = ParsedSources(raw)

	return err
}

// Helpers for json-struct Operation

// MarshalJSON prevents nil slices in json
func (data Operation) MarshalJSON() ([]byte, error) {
	type alias Operation
	var raw = alias(data)

	if raw.DocLines == nil {
		raw.DocLines = []string{}
	}

	if raw.InputArgs == nil {
		raw.InputArgs = []Field{}
	}

	if raw.OutputArgs == nil {
		raw.OutputArgs = []Field{}
	}

	if raw.CommentLines == nil {
		raw.CommentLines = []string{}
	}

	return json.Marshal(raw)
}

// UnmarshalJSON prevents nil slices from json
func (data *Operation) UnmarshalJSON(b []byte) error {
	type alias Operation
	var raw alias
	err := json.Unmarshal(b, &raw)

	if raw.DocLines == nil {
		raw.DocLines = []string{}
	}

	if raw.InputArgs == nil {
		raw.InputArgs = []Field{}
	}

	if raw.OutputArgs == nil {
		raw.OutputArgs = []Field{}
	}

	if raw.CommentLines == nil {
		raw.CommentLines = []string{}
	}

	*data = Operation(raw)

	return err
}

// Helpers for json-struct Struct

// MarshalJSON prevents nil slices in json
func (data Struct) MarshalJSON() ([]byte, error) {
	type alias Struct
	var raw = alias(data)

	if raw.DocLines == nil {
		raw.DocLines = []string{}
	}

	if raw.Fields == nil {
		raw.Fields = []Field{}
	}

	if raw.Operations == nil {
		raw.Operations = []*Operation{}
	}

	if raw.CommentLines == nil {
		raw.CommentLines = []string{}
	}

	return json.Marshal(raw)
}

// UnmarshalJSON prevents nil slices from json
func (data *Struct) UnmarshalJSON(b []byte) error {
	type alias Struct
	var raw alias
	err := json.Unmarshal(b, &raw)

	if raw.DocLines == nil {
		raw.DocLines = []string{}
	}

	if raw.Fields == nil {
		raw.Fields = []Field{}
	}

	if raw.Operations == nil {
		raw.Operations = []*Operation{}
	}

	if raw.CommentLines == nil {
		raw.CommentLines = []string{}
	}

	*data = Struct(raw)

	return err
}

// Helpers for json-struct Interface

// MarshalJSON prevents nil slices in json
func (data Interface) MarshalJSON() ([]byte, error) {
	type alias Interface
	var raw = alias(data)

	if raw.DocLines == nil {
		raw.DocLines = []string{}
	}

	if raw.Methods == nil {
		raw.Methods = []Operation{}
	}

	if raw.CommentLines == nil {
		raw.CommentLines = []string{}
	}

	return json.Marshal(raw)
}

// UnmarshalJSON prevents nil slices from json
func (data *Interface) UnmarshalJSON(b []byte) error {
	type alias Interface
	var raw alias
	err := json.Unmarshal(b, &raw)

	if raw.DocLines == nil {
		raw.DocLines = []string{}
	}

	if raw.Methods == nil {
		raw.Methods = []Operation{}
	}

	if raw.CommentLines == nil {
		raw.CommentLines = []string{}
	}

	*data = Interface(raw)

	return err
}

// Helpers for json-struct Field

// MarshalJSON prevents nil slices in json
func (data Field) MarshalJSON() ([]byte, error) {
	type alias Field
	var raw = alias(data)

	if raw.DocLines == nil {
		raw.DocLines = []string{}
	}

	if raw.CommentLines == nil {
		raw.CommentLines = []string{}
	}

	return json.Marshal(raw)
}

// UnmarshalJSON prevents nil slices from json
func (data *Field) UnmarshalJSON(b []byte) error {
	type alias Field
	var raw alias
	err := json.Unmarshal(b, &raw)

	if raw.DocLines == nil {
		raw.DocLines = []string{}
	}

	if raw.CommentLines == nil {
		raw.CommentLines = []string{}
	}

	*data = Field(raw)

	return err
}

// Helpers for json-struct Typedef

// MarshalJSON prevents nil slices in json
func (data Typedef) MarshalJSON() ([]byte, error) {
	type alias Typedef
	var raw = alias(data)

	if raw.DocLines == nil {
		raw.DocLines = []string{}
	}

	return json.Marshal(raw)
}

// UnmarshalJSON prevents nil slices from json
func (data *Typedef) UnmarshalJSON(b []byte) error {
	type alias Typedef
	var raw alias
	err := json.Unmarshal(b, &raw)

	if raw.DocLines == nil {
		raw.DocLines = []string{}
	}

	*data = Typedef(raw)

	return err
}

// Helpers for json-struct Enum

// MarshalJSON prevents nil slices in json
func (data Enum) MarshalJSON() ([]byte, error) {
	type alias Enum
	var raw = alias(data)

	if raw.DocLines == nil {
		raw.DocLines = []string{}
	}

	if raw.EnumLiterals == nil {
		raw.EnumLiterals = []EnumLiteral{}
	}

	if raw.CommentLines == nil {
		raw.CommentLines = []string{}
	}

	return json.Marshal(raw)
}

// UnmarshalJSON prevents nil slices from json
func (data *Enum) UnmarshalJSON(b []byte) error {
	type alias Enum
	var raw alias
	err := json.Unmarshal(b, &raw)

	if raw.DocLines == nil {
		raw.DocLines = []string{}
	}

	if raw.EnumLiterals == nil {
		raw.EnumLiterals = []EnumLiteral{}
	}

	if raw.CommentLines == nil {
		raw.CommentLines = []string{}
	}

	*data = Enum(raw)

	return err
}

// Helpers for json-struct EnumLiteral
