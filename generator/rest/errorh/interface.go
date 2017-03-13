package errorh

type Error struct {
	httpCode        int
	underlyingError error
	ErrorMessage    string       `json:"errorMessage"`
	ErrorCode       int          `json:"errorCode"`
	FieldErrors     []FieldError `json:"fieldErrors"` // only applicable for invalidinput
}

type FieldError struct {
	SubCode int      `json:"subCode"`
	Field   string   `json:"field"`
	Msg     string   `json:"msg"`
	Args    []string `json:"args"`
}

type HttpCode interface {
	GetHttpCode() int
}

func GetErrorCode(err error) int {
	if err != nil {
		if specificError, ok := err.(ErrorCode); ok {
			return specificError.GetErrorCode()
		}
	}
	return 0
}

type ErrorCode interface {
	GetErrorCode() int
}

func GetHttpCode(err error) int {
	if err != nil {
		if specificError, ok := err.(HttpCode); ok {
			return specificError.GetHttpCode()
		}
	}
	return 500
}

type Internal interface {
	error
	IsInternalError() bool
}

func IsInternalError(err error) bool {
	if err != nil {
		if specificError, ok := err.(Internal); ok {
			return specificError.IsInternalError()
		}
	}
	return false
}

type NotImplemented interface {
	error
	IsNotImplementedError() bool
}

func IsNotImplementedError(err error) bool {
	if err != nil {
		if specificError, ok := err.(NotImplemented); ok {
			return specificError.IsNotImplementedError()
		}
	}
	return false
}

type InvalidInput interface {
	error
	IsInvalidInputError() bool
	GetFieldErrors() []FieldError
}

func IsInvalidInputError(err error) bool {
	if err != nil {
		if specificError, ok := err.(InvalidInput); ok {
			return specificError.IsInvalidInputError()
		}
	}
	return false
}

func GetFieldErrors(err error) []FieldError {
	if err != nil {
		if specificError, ok := err.(InvalidInput); ok {
			return specificError.GetFieldErrors()
		}
	}
	return []FieldError{}
}

type NotAuthorized interface {
	error
	IsNotAuthorizedError() bool
}

func IsNotAuthorizedError(err error) bool {
	if err != nil {
		if specificError, ok := err.(NotAuthorized); ok {
			return specificError.IsNotAuthorizedError()
		}
	}
	return false
}

type NotFound interface {
	error
	IsNotFoundError() bool
}

func IsNotFoundError(err error) bool {
	if err != nil {
		if specificError, ok := err.(NotFound); ok {
			return specificError.IsNotFoundError()
		}
	}
	return false
}

type Conflict interface {
	error
	IsConflictError() bool
}

func IsConflictError(err error) bool {
	if err != nil {
		if specificError, ok := err.(Conflict); ok {
			return specificError.IsConflictError()
		}
	}
	return false
}
