package errorh

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Error struct {
	underlyingError error
	errorType       int
	errorCode       int
	fieldErrors     []FieldError // only applicable for invalidinput
}

const (
	errorTypeInternal      = 500
	errorTypeInvalidInput  = 400
	errorTypeNotFound      = 404
	errorTypeNotAuthorized = 403
)

func NewInternalErrorf(code int, format string, args ...interface{}) *Error {
	return NewInternalError(code, fmt.Errorf(format, args...))
}

func NewInternalError(code int, err error) *Error {
	newError := new(Error)
	newError.errorCode = code
	newError.underlyingError = err
	newError.errorType = errorTypeInternal
	return newError
}

func NewInvalidInputErrorf(code int, format string, args ...interface{}) *Error {
	newError := new(Error)
	newError.errorCode = code
	newError.underlyingError = fmt.Errorf(format, args...)
	newError.errorType = errorTypeInvalidInput
	return newError
}

func NewInvalidInputErrorSpecific(code int, fieldErrors []FieldError) *Error {
	newError := new(Error)
	newError.errorCode = code
	newError.underlyingError = errors.New("Input validation error")
	newError.errorType = errorTypeInvalidInput
	newError.fieldErrors = fieldErrors
	return newError
}

func NewNotFoundErrorf(code int, format string, args ...interface{}) *Error {
	return NewNotFoundError(code, fmt.Errorf(format, args...))
}

func NewNotFoundError(code int, err error) *Error {
	newError := new(Error)
	newError.errorCode = code
	newError.underlyingError = err
	newError.errorType = errorTypeNotFound
	return newError
}

func NewNotAuthorizedErrorf(code int, format string, args ...interface{}) *Error {
	return NewNotAuthorizedError(code, fmt.Errorf(format, args...))
}

func NewNotAuthorizedError(code int, err error) *Error {
	newError := new(Error)
	newError.errorCode = code
	newError.underlyingError = err
	newError.errorType = errorTypeNotAuthorized
	return newError
}

func (err Error) Error() string {
	return err.underlyingError.Error()
}

func (err Error) IsInternalError() bool {
	return err.errorType == errorTypeInternal
}

func (err Error) IsInvalidInputError() bool {
	return err.errorType == errorTypeInvalidInput
}

func (err Error) GetFieldErrors() []FieldError {
	return err.fieldErrors
}

func (err Error) IsNotFoundError() bool {
	return err.errorType == errorTypeNotFound
}

func (err Error) IsNotAuthorizedError() bool {
	return err.errorType == errorTypeNotAuthorized
}

func (err Error) GetHttpCode() int {
	return err.errorType
}

func (err Error) GetErrorCode() int {
	return err.errorCode
}

func HandleHttpError(err error, w http.ResponseWriter) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(GetHttpCode(err))

	errorBody := struct {
		ErrorCode    int          `json:"errorCode"`
		ErrorMessage string       `json:"errorMessage"`
		FieldErrors  []FieldError `json:"fieldErrors"`
	}{
		ErrorCode:    geErrorCode(err),
		ErrorMessage: err.Error(),
		FieldErrors:  getFieldErrors(err),
	}
	json.NewEncoder(w).Encode(errorBody)
}

func getFieldErrors(err error) []FieldError {
	if IsInvalidInputError(err) {
		e, _ := err.(InvalidInput)
		return e.GetFieldErrors()
	} else {
		return []FieldError{}
	}
}

func geErrorCode(err error) int {
	if IsErrorWithCodes(err) {
		e, _ := err.(InvalidInput)
		return e.GetErrorCode()
	} else {
		return 0
	}
}
