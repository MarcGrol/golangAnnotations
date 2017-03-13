package errorh

import (
	"errors"
	"fmt"
	"net/http"
)

func NewInternalErrorf(code int, format string, args ...interface{}) *Error {
	return NewInternalError(code, fmt.Errorf(format, args...))
}

func NewInternalError(code int, err error) *Error {
	newError := new(Error)
	newError.ErrorCode = code
	newError.underlyingError = err
	newError.ErrorMessage = err.Error()
	newError.httpCode = http.StatusInternalServerError
	return newError
}

func NewNotImplementedErrorf(code int, format string, args ...interface{}) *Error {
	return NewNotImplementedError(code, fmt.Errorf(format, args...))
}

func NewNotImplementedError(code int, err error) *Error {
	newError := new(Error)
	newError.ErrorCode = code
	newError.underlyingError = err
	newError.ErrorMessage = err.Error()
	newError.httpCode = http.StatusNotImplemented
	return newError
}

func NewConflictErrorf(code int, format string, args ...interface{}) *Error {
	return NewConflictError(code, fmt.Errorf(format, args...))
}

func NewConflictError(code int, err error) *Error {
	newError := new(Error)
	newError.ErrorCode = code
	newError.underlyingError = err
	newError.ErrorMessage = err.Error()
	newError.httpCode = http.StatusConflict
	return newError
}

func NewInvalidInputErrorf(code int, format string, args ...interface{}) *Error {
	newError := new(Error)
	newError.ErrorCode = code
	newError.underlyingError = fmt.Errorf(format, args...)
	newError.ErrorMessage = newError.underlyingError.Error()
	newError.httpCode = http.StatusBadRequest
	return newError
}

func NewInvalidInputErrorSpecific(code int, fieldErrors []FieldError) *Error {
	newError := new(Error)
	newError.ErrorCode = code
	newError.underlyingError = errors.New("Input validation error")
	newError.ErrorMessage = newError.underlyingError.Error()
	newError.httpCode = http.StatusBadRequest
	newError.FieldErrors = fieldErrors
	return newError
}

func NewNotFoundErrorf(code int, format string, args ...interface{}) *Error {
	return NewNotFoundError(code, fmt.Errorf(format, args...))
}

func NewNotFoundError(code int, err error) *Error {
	newError := new(Error)
	newError.ErrorCode = code
	newError.underlyingError = err
	newError.ErrorMessage = err.Error()
	newError.httpCode = http.StatusNotFound
	return newError
}

func NewNotAuthorizedErrorf(code int, format string, args ...interface{}) *Error {
	return NewNotAuthorizedError(code, fmt.Errorf(format, args...))
}

func NewNotAuthorizedError(code int, err error) *Error {
	newError := new(Error)
	newError.ErrorCode = code
	newError.underlyingError = err
	newError.ErrorMessage = newError.underlyingError.Error()
	newError.httpCode = http.StatusForbidden
	return newError
}

func (error Error) Error() string {
	return error.ErrorMessage
}

func (error Error) IsInternalError() bool {
	return error.httpCode == http.StatusInternalServerError
}

func (error Error) IsNotImplementedError() bool {
	return error.httpCode == http.StatusNotImplemented
}

func (error Error) IsInvalidInputError() bool {
	return error.httpCode == http.StatusBadRequest
}

func (error Error) GetFieldErrors() []FieldError {
	return error.FieldErrors
}

func (error Error) IsNotFoundError() bool {
	return error.httpCode == http.StatusNotFound
}

func (error Error) IsConflictError() bool {
	return error.httpCode == http.StatusConflict
}

func (error Error) IsNotAuthorizedError() bool {
	return error.httpCode == http.StatusForbidden
}

func (error Error) GetHttpCode() int {
	return error.httpCode
}

func (error Error) GetErrorCode() int {
	return error.ErrorCode
}

func GetFieldErrors(err error) []FieldError {
	return getFieldErrors(err)
}

func getFieldErrors(err error) []FieldError {
	if IsInvalidInputError(err) {
		e, _ := err.(InvalidInput)
		return e.GetFieldErrors()
	} else {
		return []FieldError{}
	}
}
