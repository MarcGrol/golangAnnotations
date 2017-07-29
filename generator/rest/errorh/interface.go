package errorh

//go:generate golangAnnotations -input-dir .

// @JsonStruct()
type Error struct {
	message      string       `json:"-"`
	httpCode     int          `json:"-"`
	ErrorMessage string       `json:"errorMessage"`
	ErrorCode    int          `json:"errorCode"`
	FieldErrors  []FieldError `json:"fieldErrors"` // only applicable for invalidInput
}

// @JsonStruct()
type FieldError struct {
	SubCode int      `json:"subCode"`
	Field   string   `json:"field"`
	Msg     string   `json:"msg"`
	Args    []string `json:"args"`
}

func MapToError(err error) Error {
	return Error{
		message:      GetMessage(err),
		httpCode:     GetHttpCode(err),
		ErrorMessage: err.Error(),
		ErrorCode:    GetErrorCode(err),
		FieldErrors:  GetFieldErrors(err),
	}
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
		if internal, ok := err.(Internal); ok {
			return internal.IsInternalError()
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
		if notImplemented, ok := err.(NotImplemented); ok {
			return notImplemented.IsNotImplementedError()
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
		if invalidInput, ok := err.(InvalidInput); ok {
			return invalidInput.IsInvalidInputError()
		}
	}
	return false
}

func GetFieldErrors(err error) []FieldError {
	if err != nil {
		if invalidInput, ok := err.(InvalidInput); ok {
			return invalidInput.GetFieldErrors()
		}
	}
	return []FieldError{}
}

type NotAuthorized interface {
	HttpError
	IsNotAuthorizedError() bool
	GetMessage() string
}

func IsNotAuthorizedError(err error) bool {
	if err != nil {
		if notAuthorized, ok := err.(NotAuthorized); ok {
			return notAuthorized.IsNotAuthorizedError()
		}
	}
	return false
}

func GetMessage(err error) string {
	if err != nil {
		if notAuthorized, ok := err.(NotAuthorized); ok {
			return notAuthorized.GetMessage()
		}
	}
	return ""
}

type NotFound interface {
	HttpError
	IsNotFoundError() bool
}

func IsNotFoundError(err error) bool {
	if err != nil {
		if notFound, ok := err.(NotFound); ok {
			return notFound.IsNotFoundError()
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
		if conflict, ok := err.(Conflict); ok {
			return conflict.IsConflictError()
		}
	}
	return false
}
