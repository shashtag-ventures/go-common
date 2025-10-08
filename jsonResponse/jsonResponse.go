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

// SendErrorResponse sends a consistent JSON error response.
// It handles validation errors specifically and provides a generic error response otherwise.
func SendErrorResponse(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")

	var errorResponse ErrorResponse
	errorResponse.Status = statusCode
	errorResponse.Error = http.StatusText(statusCode)
	errorResponse.Message = err.Error() // Default message

	var validationErrs validator.ValidationErrors
	// Check if the error is a validation error.
	if errors.As(err, &validationErrs) {
		var errMSGS []string
		// Iterate over validation errors and create human-readable messages.
		for _, e := range validationErrs {
			switch e.ActualTag() {
			case "required":
				errMSGS = append(errMSGS, e.Field()+" is required")
			case "email":
				errMSGS = append(errMSGS, e.Field()+" should be a valid email address")
			default:
				errMSGS = append(errMSGS, e.Field()+" is invalid")
			}
		}
		errorResponse.Status = http.StatusBadRequest
		errorResponse.Error = "Validation Error"
		errorResponse.Message = strings.Join(errMSGS, ", ")
	} else {
		// Check for custom errors and map them to appropriate status codes.
		for customErr, mappedStatusCode := range errorStatusCodeMap {
			if errors.Is(err, customErr) {
				errorResponse.Status = mappedStatusCode
				errorResponse.Error = http.StatusText(mappedStatusCode)
				if customErr == customErrors.ErrInternal {
					errorResponse.Message = "An unexpected internal server error occurred."
				} else {
					errorResponse.Message = err.Error()
				}
				break
			}
		}
	}

	w.WriteHeader(errorResponse.Status)
	json.NewEncoder(w).Encode(errorResponse)
}

// JsonError creates a generic error for internal use.
// It wraps the original error with a custom message.
func JsonError(statusCode int, err error, message string) error {
	return fmt.Errorf("%s: %w", message, err)
}
