package request

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// DecodeAndValidate decodes the request body into v and validates it.
func DecodeAndValidate(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		// Consider returning a more specific decoding error
		return err
	}
	defer r.Body.Close()

	if err := validate.Struct(v); err != nil {
		return err
	}

	return nil
}
