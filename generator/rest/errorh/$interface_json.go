// Generated automatically by golangAnnotations: do not edit manually

package errorh

import "encoding/json"

// Helpers for json-struct Error

// MarshalJSON prevents nil slices in json
func (data Error) MarshalJSON() ([]byte, error) {
	type alias Error
	var raw = alias(data)

	if raw.FieldErrors == nil {
		raw.FieldErrors = []FieldError{}
	}

	return json.Marshal(raw)
}

// UnmarshalJSON prevents nil slices from json
func (data *Error) UnmarshalJSON(b []byte) error {
	type alias Error
	var raw alias
	err := json.Unmarshal(b, &raw)

	if raw.FieldErrors == nil {
		raw.FieldErrors = []FieldError{}
	}

	*data = Error(raw)

	return err
}

// Helpers for json-struct FieldError

// MarshalJSON prevents nil slices in json
func (data FieldError) MarshalJSON() ([]byte, error) {
	type alias FieldError
	var raw = alias(data)

	if raw.Args == nil {
		raw.Args = []string{}
	}

	return json.Marshal(raw)
}

// UnmarshalJSON prevents nil slices from json
func (data *FieldError) UnmarshalJSON(b []byte) error {
	type alias FieldError
	var raw alias
	err := json.Unmarshal(b, &raw)

	if raw.Args == nil {
		raw.Args = []string{}
	}

	*data = FieldError(raw)

	return err
}
