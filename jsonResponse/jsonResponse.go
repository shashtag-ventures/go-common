package jsonResponse

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	customErrors "github.com/shashtag-ventures/go-common/errors"
)

// ErrorResponse defines the structure for standardized API error responses.
type ErrorResponse struct {
	Status  int    `json:"status"`  // HTTP status code
	Error   string `json:"error"`   // A short, human-readable summary of the error
	Message string `json:"message"` // A more detailed, human-readable message about the error
}

// errorStatusCodeMap maps custom error types to their corresponding HTTP status codes.
var errorStatusCodeMap = map[error]int{
	customErrors.ErrNotFound:      http.StatusNotFound,
	customErrors.ErrInvalidInput:  http.StatusBadRequest,
	customErrors.ErrUnauthorized:  http.StatusUnauthorized,
	customErrors.ErrForbidden:     http.StatusForbidden,
	customErrors.ErrAlreadyExists: http.StatusConflict,
	customErrors.ErrInternal:      http.StatusInternalServerError,
}

// JsonResponse sends a JSON response with the given status code and data.
func JsonResponse(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

// SendErrorResponse sends a consistent JSON error response using the given statusCode.
// It formats validation errors into human-readable messages but does not override the status code.
// For automatic status code detection from error types, use SendAutoErrorResponse instead.
func SendErrorResponse(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")

	var errorResponse ErrorResponse
	errorResponse.Status = statusCode
	errorResponse.Error = http.StatusText(statusCode)
	if err != nil {
		errorResponse.Message = err.Error()
	} else {
		errorResponse.Message = errorResponse.Error
	}

	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		var errMSGS []string
		for _, e := range validationErrs {
			switch e.ActualTag() {
			case "required":
				errMSGS = append(errMSGS, e.Field()+" is required")
			case "email":
				errMSGS = append(errMSGS, e.Field()+" should be a valid email address")
			case "slug":
				errMSGS = append(errMSGS, e.Field()+" must be a valid slug (lowercase alphanumeric and hyphens)")
			case "git-url":
				errMSGS = append(errMSGS, e.Field()+" must be a valid git repository URL (http/https)")
			default:
				errMSGS = append(errMSGS, e.Field()+" is invalid")
			}
		}
		errorResponse.Error = "Validation Error"
		errorResponse.Message = strings.Join(errMSGS, ", ")
	} else if errors.Is(err, customErrors.ErrInternal) {
		errorResponse.Message = "An unexpected internal server error occurred."
	}

	w.WriteHeader(errorResponse.Status)
	json.NewEncoder(w).Encode(errorResponse)
}

// SendAutoErrorResponse automatically determines the HTTP status code based on the error type
// and sends a consistent JSON error response. If no mapping is found, it defaults to 500.
func SendAutoErrorResponse(w http.ResponseWriter, err error) {
	statusCode := http.StatusInternalServerError

	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		statusCode = http.StatusBadRequest
	} else {
		for customErr, mappedStatusCode := range errorStatusCodeMap {
			if errors.Is(err, customErr) {
				statusCode = mappedStatusCode
				break
			}
		}
	}

	SendErrorResponse(w, err, statusCode)
}

// JsonError creates a generic error for internal use.
// It wraps the original error with a custom message.
func JsonError(statusCode int, err error, message string) error {
	return fmt.Errorf("%s: %w", message, err)
}
