package errorh

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	{
		err := makeNil()
		assert.Nil(t, err)
		assert.False(t, IsInternalError(err))
		assert.False(t, IsInvalidInputError(err))
		assert.False(t, IsNotFoundError(err))
		assert.Equal(t, 500, GetHttpCode(err))
		assert.Equal(t, 0, GetErrorCode(err))
	}
	{
		err := errors.New("my unclassified error")
		assert.False(t, IsInternalError(err))
		assert.False(t, IsInvalidInputError(err))
		assert.False(t, IsNotFoundError(err))
		assert.False(t, IsNotAuthorizedError(err))
		assert.Equal(t, "my unclassified error", err.Error())
		assert.Equal(t, 500, GetHttpCode(err))
		assert.Equal(t, 0, GetErrorCode(err))
	}
	{
		err := NewInternalErrorf(1, "my %s error", "internal")
		assert.True(t, IsInternalError(err))
		assert.False(t, IsInvalidInputError(err))
		assert.False(t, IsNotFoundError(err))
		assert.False(t, IsNotAuthorizedError(err))
		assert.Equal(t, "my internal error", err.Error())
		assert.Equal(t, 500, GetHttpCode(err))
		assert.Equal(t, 1, GetErrorCode(err))
	}
	{
		err := NewInvalidInputErrorSpecific(2, []FieldError{
			{
				SubCode: 112,
				Field:   "email",
				Msg:     "Invalid value %s for parameter 'email'",
				Args:    []string{"test@home.nl"},
			},
		})
		assert.False(t, IsInternalError(err))
		assert.False(t, IsNotFoundError(err))
		assert.False(t, IsNotAuthorizedError(err))
		assert.True(t, IsInvalidInputError(err))
		assert.Equal(t, "Input validation error", err.Error())
		assert.Equal(t, 400, GetHttpCode(err))
		assert.Equal(t, 2, GetErrorCode(err))
	}
	{
		err := NewNotFoundErrorf(3, "my %s error", "not found")
		assert.False(t, IsInternalError(err))
		assert.False(t, IsInvalidInputError(err))
		assert.True(t, IsNotFoundError(err))
		assert.False(t, IsNotAuthorizedError(err))
		assert.Equal(t, "my not found error", err.Error())
		assert.Equal(t, 404, GetHttpCode(err))
		assert.Equal(t, 3, GetErrorCode(err))
	}
	{
		err := NewNotAuthorizedErrorf(4, "my %s error", "not authorized")
		assert.False(t, IsInternalError(err))
		assert.False(t, IsInvalidInputError(err))
		assert.False(t, IsNotFoundError(err))
		assert.True(t, IsNotAuthorizedError(err))
		assert.Equal(t, "my not authorized error", err.Error())
		assert.Equal(t, 403, GetHttpCode(err))
		assert.Equal(t, 4, GetErrorCode(err))
	}
}

func makeNil() error {
	return nil
}
