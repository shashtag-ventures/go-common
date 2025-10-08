package jsonResponse_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	customErrors "github.com/shashtag-ventures/go-common/errors" // Import customErrors
	"github.com/shashtag-ventures/go-common/jsonResponse"
	"github.com/stretchr/testify/assert"
)

// Mock validator.FieldError for testing
type mockFieldError struct {
	fieldName string
	tagName   string
}

func (m mockFieldError) Tag() string             { return m.tagName }
func (m mockFieldError) ActualTag() string       { return m.tagName }
func (m mockFieldError) Namespace() string       { return "" }
func (m mockFieldError) StructNamespace() string { return "" }
func (m mockFieldError) Field() string           { return m.fieldName }
func (m mockFieldError) StructField() string     { return "" }
func (m mockFieldError) Value() interface{}      { return nil }
func (m mockFieldError) Param() string           { return "" }
func (m mockFieldError) Kind() reflect.Kind      { return reflect.Invalid }
func (m mockFieldError) Type() reflect.Type      { return nil }
func (m mockFieldError) Error() string {
	return fmt.Sprintf("Key: '%s' Error:Field validation for '%s' failed on the '%s' tag", m.fieldName, m.fieldName, m.tagName)
}
func (m mockFieldError) Translate(trans ut.Translator) string { return "" }

func TestJsonResponse(t *testing.T) {
	// Test Case 1: Successful JSON response with map data
	t.Run("Map Data", func(t *testing.T) {
		rr := httptest.NewRecorder()
		data := map[string]string{"message": "success"}
		statusCode := http.StatusOK

		err := jsonResponse.JsonResponse(rr, statusCode, data)

		assert.NoError(t, err)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		assert.Equal(t, statusCode, rr.Code)

		expectedBody, _ := json.Marshal(data)
		assert.Equal(t, string(expectedBody)+"\n", rr.Body.String())
	})

	// Test Case 2: Successful JSON response with struct data
	t.Run("Struct Data", func(t *testing.T) {
		rr := httptest.NewRecorder()
		type Response struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		data := Response{Name: "Test", Age: 30}
		statusCode := http.StatusCreated

		err := jsonResponse.JsonResponse(rr, statusCode, data)

		assert.NoError(t, err)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		assert.Equal(t, statusCode, rr.Code)

		expectedBody, _ := json.Marshal(data)
		assert.Equal(t, string(expectedBody)+"\n", rr.Body.String())
	})
}

func TestSendErrorResponse(t *testing.T) {
	// Test Case 1: Validation Error (required field)
	t.Run("Validation Error - Required", func(t *testing.T) {
		rr := httptest.NewRecorder()
		valErrors := validator.ValidationErrors([]validator.FieldError{
			mockFieldError{fieldName: "Name", tagName: "required"},
		})
		statusCode := http.StatusBadRequest

		jsonResponse.SendErrorResponse(rr, valErrors, statusCode)

		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		assert.Equal(t, statusCode, rr.Code)

		var response jsonResponse.ErrorResponse
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, statusCode, response.Status)
		assert.Equal(t, "Validation Error", response.Error)
		assert.Equal(t, "Name is required", response.Message)
	})

	// Test Case 2: Validation Error (email format)
	t.Run("Validation Error - Email", func(t *testing.T) {
		rr := httptest.NewRecorder()
		valErrors := validator.ValidationErrors([]validator.FieldError{
			mockFieldError{fieldName: "Email", tagName: "email"},
		})
		statusCode := http.StatusBadRequest

		jsonResponse.SendErrorResponse(rr, valErrors, statusCode)

		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		assert.Equal(t, statusCode, rr.Code)

		var response jsonResponse.ErrorResponse
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, statusCode, response.Status)
		assert.Equal(t, "Validation Error", response.Error)
		assert.Equal(t, "Email should be a valid email address", response.Message)
	})

	// Test Case 3: Generic Error
	t.Run("Generic Error", func(t *testing.T) {
		rr := httptest.NewRecorder()
		genericErr := errors.New("something went wrong")
		statusCode := http.StatusInternalServerError

		jsonResponse.SendErrorResponse(rr, genericErr, statusCode)

		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		assert.Equal(t, statusCode, rr.Code)

		var response jsonResponse.ErrorResponse
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, statusCode, response.Status)
		assert.Equal(t, http.StatusText(statusCode), response.Error)
		assert.Equal(t, genericErr.Error(), response.Message)
	})

	// New Test Cases for customErrors
	t.Run("Custom Error - Not Found", func(t *testing.T) {
		rr := httptest.NewRecorder()
		err := customErrors.New("user not found", customErrors.ErrNotFound)
		statusCode := http.StatusNotFound

		jsonResponse.SendErrorResponse(rr, err, statusCode)

		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		assert.Equal(t, statusCode, rr.Code)

		var response jsonResponse.ErrorResponse
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.Equal(t, statusCode, response.Status)
		assert.Equal(t, http.StatusText(statusCode), response.Error)
		assert.Equal(t, err.Error(), response.Message)
	})

	t.Run("Custom Error - Invalid Input", func(t *testing.T) {
		rr := httptest.NewRecorder()
		err := customErrors.New("invalid email format", customErrors.ErrInvalidInput)
		statusCode := http.StatusBadRequest

		jsonResponse.SendErrorResponse(rr, err, statusCode)

		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		assert.Equal(t, statusCode, rr.Code)

		var response jsonResponse.ErrorResponse
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.Equal(t, statusCode, response.Status)
		assert.Equal(t, http.StatusText(statusCode), response.Error)
		assert.Equal(t, err.Error(), response.Message)
	})

	t.Run("Custom Error - Unauthorized", func(t *testing.T) {
		rr := httptest.NewRecorder()
		err := customErrors.New("invalid credentials", customErrors.ErrUnauthorized)
		statusCode := http.StatusUnauthorized

		jsonResponse.SendErrorResponse(rr, err, statusCode)

		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		assert.Equal(t, statusCode, rr.Code)

		var response jsonResponse.ErrorResponse
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.Equal(t, statusCode, response.Status)
		assert.Equal(t, http.StatusText(statusCode), response.Error)
		assert.Equal(t, err.Error(), response.Message)
	})

	t.Run("Custom Error - Forbidden", func(t *testing.T) {
		rr := httptest.NewRecorder()
		err := customErrors.New("access denied", customErrors.ErrForbidden)
		statusCode := http.StatusForbidden

		jsonResponse.SendErrorResponse(rr, err, statusCode)

		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		assert.Equal(t, statusCode, rr.Code)

		var response jsonResponse.ErrorResponse
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.Equal(t, statusCode, response.Status)
		assert.Equal(t, http.StatusText(statusCode), response.Error)
		assert.Equal(t, err.Error(), response.Message)
	})

	t.Run("Custom Error - Already Exists", func(t *testing.T) {
		rr := httptest.NewRecorder()
		err := customErrors.New("user already exists", customErrors.ErrAlreadyExists)
		statusCode := http.StatusConflict

		jsonResponse.SendErrorResponse(rr, err, statusCode)

		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		assert.Equal(t, statusCode, rr.Code)

		var response jsonResponse.ErrorResponse
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.Equal(t, statusCode, response.Status)
		assert.Equal(t, http.StatusText(statusCode), response.Error)
		assert.Equal(t, err.Error(), response.Message)
	})

	t.Run("Custom Error - Internal Error", func(t *testing.T) {
		rr := httptest.NewRecorder()
		err := customErrors.New("database connection failed", customErrors.ErrInternal)
		statusCode := http.StatusInternalServerError

		jsonResponse.SendErrorResponse(rr, err, statusCode)

		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		assert.Equal(t, statusCode, rr.Code)

		var response jsonResponse.ErrorResponse
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.Equal(t, statusCode, response.Status)
		assert.Equal(t, http.StatusText(statusCode), response.Error)
		assert.Equal(t, "An unexpected internal server error occurred.", response.Message) // Generic message
	})

	t.Run("Custom Error - Wrapped Custom Error", func(t *testing.T) {
		rr := httptest.NewRecorder()
		originalErr := errors.New("underlying db error")
		wrappedErr := customErrors.New("failed to fetch user", originalErr) // This is not a customErrors.Err* type
		statusCode := http.StatusInternalServerError

		jsonResponse.SendErrorResponse(rr, wrappedErr, statusCode)

		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		assert.Equal(t, statusCode, rr.Code)

		var response jsonResponse.ErrorResponse
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.Equal(t, statusCode, response.Status)
		assert.Equal(t, http.StatusText(statusCode), response.Error)
		assert.Equal(t, wrappedErr.Error(), response.Message)
	})
}

func TestJsonError(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := jsonResponse.JsonError(http.StatusBadRequest, originalErr, "custom message")

	assert.Error(t, wrappedErr)
	assert.Contains(t, wrappedErr.Error(), "custom message")
	assert.True(t, errors.Is(wrappedErr, originalErr))
}
