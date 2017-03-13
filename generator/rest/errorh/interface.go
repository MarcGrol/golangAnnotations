package errorh

//go:generate golangAnnotations -input-dir .

// @JsonStruct()
type Error struct {
	httpCode        int          `json:"-"`
	underlyingError error        `json:"-"`
	ErrorMessage    string       `json:"errorMessage"`
	ErrorCode       int          `json:"errorCode"`
	FieldErrors     []FieldError `json:"fieldErrors"` // only applicable for invalidinput
}

// @JsonStruct()
type FieldError struct {
	SubCode int      `json:"subCode"`
	Field   string   `json:"field"`
	Msg     string   `json:"msg"`
	Args    []string `json:"args"`
}

type HttpError interface {
	error
	GetHttpCode() int
	GetErrorCode() int
}

func GetHttpCode(err error) int {
	if err != nil {
		if httpError, ok := err.(HttpError); ok {
			return httpError.GetHttpCode()
		}
	}
	return 500
}

func GetErrorCode(err error) int {
	if err != nil {
		if httpError, ok := err.(HttpError); ok {
			return httpError.GetErrorCode()
		}
	}
	return 0
}

type Internal interface {
	HttpError
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
	HttpError
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
	HttpError
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
	HttpError
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
	HttpError
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
	HttpError
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
