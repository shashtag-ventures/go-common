package errors_test

import (
	"errors"
	"testing"

	customErrors "github.com/shashtag-ventures/go-common/errors"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("New with nil error", func(t *testing.T) {
		msg := "test message"
		err := customErrors.New(msg, nil)
		assert.Error(t, err)
		assert.Equal(t, msg, err.Error())
	})

	t.Run("New with non-nil error", func(t *testing.T) {
		originalErr := errors.New("original error")
		msg := "wrapped message"
		wrappedErr := customErrors.New(msg, originalErr)
		assert.Error(t, wrappedErr)
		assert.Contains(t, wrappedErr.Error(), msg)
		assert.True(t, errors.Is(wrappedErr, originalErr))
	})

	t.Run("New with custom error", func(t *testing.T) {
		msg := "failed operation"
		wrappedErr := customErrors.New(msg, customErrors.ErrNotFound)
		assert.Error(t, wrappedErr)
		assert.Contains(t, wrappedErr.Error(), msg)
		assert.True(t, errors.Is(wrappedErr, customErrors.ErrNotFound))
	})
}

func TestCustomErrorVariables(t *testing.T) {
	t.Run("ErrNotFound is distinct", func(t *testing.T) {
		assert.Equal(t, "resource not found", customErrors.ErrNotFound.Error())
		assert.False(t, errors.Is(customErrors.ErrNotFound, errors.New("some other error")))
	})

	t.Run("ErrInvalidInput is distinct", func(t *testing.T) {
		assert.Equal(t, "invalid input", customErrors.ErrInvalidInput.Error())
	})

	t.Run("ErrUnauthorized is distinct", func(t *testing.T) {
		assert.Equal(t, "unauthorized", customErrors.ErrUnauthorized.Error())
	})

	t.Run("ErrForbidden is distinct", func(t *testing.T) {
		assert.Equal(t, "forbidden", customErrors.ErrForbidden.Error())
	})

	t.Run("ErrInternal is distinct", func(t *testing.T) {
		assert.Equal(t, "internal server error", customErrors.ErrInternal.Error())
	})

	t.Run("ErrAlreadyExists is distinct", func(t *testing.T) {
		assert.Equal(t, "resource already exists", customErrors.ErrAlreadyExists.Error())
	})
}