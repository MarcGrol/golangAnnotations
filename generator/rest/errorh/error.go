package errorh

import (
	"fmt"
	"net/http"
)

func NewInternalErrorf(code int, format string, args ...interface{}) *Error {
	return NewInternalError(code, fmt.Errorf(format, args...))
}

func NewInternalError(code int, err error) *Error {
	return &Error{
		httpCode:     http.StatusInternalServerError,
		ErrorMessage: err.Error(),
		ErrorCode:    code,
		FieldErrors:  []FieldError{},
	}
}

func NewNotImplementedErrorf(code int, format string, args ...interface{}) *Error {
	return NewNotImplementedError(code, fmt.Errorf(format, args...))
}

func NewNotImplementedError(code int, err error) *Error {
	return &Error{
		httpCode:     http.StatusNotImplemented,
		ErrorMessage: err.Error(),
		ErrorCode:    code,
		FieldErrors:  []FieldError{},
	}
}

func NewInvalidInputErrorf(code int, format string, args ...interface{}) *Error {
	return &Error{
		httpCode:     http.StatusBadRequest,
		ErrorMessage: fmt.Sprintf(format, args...),
		ErrorCode:    code,
		FieldErrors:  []FieldError{},
	}
}

func NewInvalidInputErrorSpecific(code int, fieldErrors []FieldError) *Error {
	return &Error{
		httpCode:     http.StatusBadRequest,
		ErrorMessage: "Input validation error",
		ErrorCode:    code,
		FieldErrors:  fieldErrors,
	}
}

func NewNotAuthorizedErrorf(code int, format string, args ...interface{}) *Error {
	return NewNotAuthorizedErrorfWithMessage(fmt.Errorf(format, args...), code, format, args...)
}

func NewNotAuthorizedError(code int, err error) *Error {
	return NewNotAuthorizedErrorfWithMessage(err, code, err.Error())
}

func NewNotAuthorizedErrorfWithMessage(err error, code int, format string, args ...interface{}) *Error {
	return &Error{
		message:      fmt.Sprintf(format, args...),
		httpCode:     http.StatusForbidden,
		ErrorMessage: err.Error(),
		ErrorCode:    code,
		FieldErrors:  []FieldError{},
	}
}
func NewNotFoundErrorf(code int, format string, args ...interface{}) *Error {
	return NewNotFoundError(code, fmt.Errorf(format, args...))
}

func NewNotFoundError(code int, err error) *Error {
	return &Error{
		httpCode:     http.StatusNotFound,
		ErrorMessage: err.Error(),
		ErrorCode:    code,
		FieldErrors:  []FieldError{},
	}
}

func NewConflictErrorf(code int, format string, args ...interface{}) *Error {
	return NewConflictError(code, fmt.Errorf(format, args...))
}

func NewConflictError(code int, err error) *Error {
	return &Error{
		httpCode:     http.StatusConflict,
		ErrorMessage: err.Error(),
		ErrorCode:    code,
		FieldErrors:  []FieldError{},
	}
}

func (error Error) Error() string {
	return error.ErrorMessage
}

func (error Error) GetMessage() string {
	return error.message
}

func (error Error) GetHttpCode() int {
	return error.httpCode
}

func (error Error) GetErrorCode() int {
	return error.ErrorCode
}

func (error Error) GetFieldErrors() []FieldError {
	return error.FieldErrors
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

func (error Error) IsNotAuthorizedError() bool {
	return error.httpCode == http.StatusForbidden
}

func (error Error) IsNotFoundError() bool {
	return error.httpCode == http.StatusNotFound
}

func (error Error) IsConflictError() bool {
	return error.httpCode == http.StatusConflict
}
