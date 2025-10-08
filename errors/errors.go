package errors

import (
	"errors"
	"fmt"
)

// Common custom error types used throughout the application.
// These provide semantic meaning to errors, allowing for more precise error handling.
var (
	ErrNotFound      = errors.New("resource not found")      // Indicates that a requested resource could not be found.
	ErrInvalidInput  = errors.New("invalid input")           // Indicates that the provided input is invalid.
	ErrUnauthorized  = errors.New("unauthorized")          // Indicates that the request lacks valid authentication credentials.
	ErrForbidden     = errors.New("forbidden")             // Indicates that the server understood the request but refuses to authorize it.
	ErrInternal      = errors.New("internal server error") // Indicates an unexpected internal server error.
	ErrAlreadyExists = errors.New("resource already exists") // Indicates that a resource with the same identifier already exists.
)

// New wraps an error with a message, preserving the original error.
// This allows for adding context to an error while still being able to check its underlying type.
func New(msg string, err error) error {
	if err == nil {
		return errors.New(msg)
	}
	return fmt.Errorf("%s: %w", msg, err)
}
