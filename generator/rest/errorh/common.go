package errorh

func FieldErrorForMissingParameter(parameter string) FieldError {
	return FieldError{
		SubCode: 1000,
		Field:   parameter,
		Msg:     "Missing value for mandatory parameter %s",
		Args:    []string{parameter},
	}
}

func FieldErrorForInvalidParameter(parameter string) FieldError {
	return FieldError{
		SubCode: 1001,
		Field:   parameter,
		Msg:     "Invalid value for mandatory parameter %s",
		Args:    []string{parameter},
	}
}
