package errorh

type ErrorWithCodes interface {
	GetHttpCode() int
	GetErrorCode() int
}

type InternalError interface {
	error
	ErrorWithCodes
	IsInternalError() bool
}

type NotFound interface {
	error
	ErrorWithCodes
	IsNotFoundError() bool
}

type NotAuthorized interface {
	error
	ErrorWithCodes
	IsNotAuthorizedError() bool
}

type FieldError struct {
	SubCode int  		`json:"subCode"`
	Field   string  	`json:"field"`
	Msg     string  	`json:"msg"`
	Args    []string 	`json:"args"`
}

type InvalidInput interface {
	error
	ErrorWithCodes
	IsInvalidInputError() bool
	GetFieldErrors() []FieldError
}

func IsErrorWithCodes(err error) bool {
	if err != nil {
		if _, ok := err.(ErrorWithCodes); ok {
			return true
		}
	}
	return false
}

func GetErrorCode(err error) int {
	if err != nil {
		if specificError, ok := err.(ErrorWithCodes); ok {
			return specificError.GetErrorCode()
		}
	}
	return 0
}

func GetHttpCode(err error) int {
	if err != nil {
		if specificError, ok := err.(ErrorWithCodes); ok {
			return specificError.GetHttpCode()
		}
	}
	return 500
}

func IsInternalError(err error) bool {
	if err != nil {
		if specificError, ok := err.(InternalError); ok {
			return specificError.IsInternalError()
		}
	}
	return false
}

func IsInvalidInputError(err error) bool {
	if err != nil {
		if specificError, ok := err.(InvalidInput); ok {
			return specificError.IsInvalidInputError()
		}
	}
	return false
}

func IsNotFoundError(err error) bool {
	if err != nil {
		if specificError, ok := err.(NotFound); ok {
			return specificError.IsNotFoundError()
		}
	}
	return false
}

func IsNotAuthorizedError(err error) bool {
	if err != nil {
		if specificError, ok := err.(NotAuthorized); ok {
			return specificError.IsNotAuthorizedError()
		}
	}
	return false
}
