package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/shashtag-ventures/go-common/netutil"
)

// Binder encapsulates request decoding and validation logic.
// It uses validator/v10 for tag-based validation and supports strict JSON decoding.
type Binder struct {
	validator *validator.Validate
}

// New creates a new Binder instance with default configurations and custom validations.
func New() *Binder {
	v := validator.New()

	// Register custom validation for 'slug'
	_ = v.RegisterValidation("slug", func(fl validator.FieldLevel) bool {
		s := fl.Field().String()
		if s == "" {
			return true // Use 'required' tag for mandatory fields
		}
		// Basic slug validation: lowercase alphanumeric and hyphens
		for _, r := range s {
			if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-') {
				return false
			}
		}
		return !strings.HasPrefix(s, "-") && !strings.HasSuffix(s, "-")
	})

	// Register custom validation for 'git-url' using netutil.ValidateGitURL
	_ = v.RegisterValidation("git-url", func(fl validator.FieldLevel) bool {
		url := fl.Field().String()
		if url == "" {
			return true // Use 'required' tag for mandatory fields
		}
		return netutil.ValidateGitURL(url) == nil
	})

	return &Binder{
		validator: v,
	}
}

// DefaultBinder is the globally shared binder instance.
var DefaultBinder = New()

// DecodeAndValidate decodes the request body into v and validates it.
// It uses the default binder.
func DecodeAndValidate(r *http.Request, v interface{}) error {
	return DefaultBinder.DecodeAndValidate(r, v)
}

// DecodeAndValidate decodes the request body into v based on Content-Type and performs validation.
func (b *Binder) DecodeAndValidate(r *http.Request, v interface{}) error {
	if r.ContentLength == 0 {
		return b.Validate(v)
	}

	contentType := r.Header.Get("Content-Type")
	switch {
	case strings.Contains(contentType, "application/json"):
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields() // Enable strict decoding
		if err := decoder.Decode(v); err != nil {
			return fmt.Errorf("invalid json: %w", err)
		}
	case strings.Contains(contentType, "application/x-www-form-urlencoded"):
		if err := r.ParseForm(); err != nil {
			return fmt.Errorf("failed to parse form: %w", err)
		}
		if err := b.bindForm(r.Form, v); err != nil {
			return err
		}
	default:
		// Fallback to JSON if no content type is specified
		if err := json.NewDecoder(r.Body).Decode(v); err != nil {
			return fmt.Errorf("failed to decode request body: %w", err)
		}
	}

	return b.Validate(v)
}

// Validate performs tag-based validation on the given struct.
func (b *Binder) Validate(v interface{}) error {
	return b.validator.Struct(v)
}

// bindForm is a simple helper to bind URL values to a struct.
// Note: This is a basic implementation. For production, consider using a library like 'gorilla/schema'.
func (b *Binder) bindForm(values map[string][]string, v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("v must be a pointer to a struct")
	}

	elem := val.Elem()
	typ := elem.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("json") // Reuse json tag for form mapping for simplicity, or use 'form'
		if tag == "" || tag == "-" {
			tag = strings.ToLower(field.Name)
		}
		// Strip omitempty etc
		tag = strings.Split(tag, ",")[0]

		if val, ok := values[tag]; ok && len(val) > 0 {
			f := elem.Field(i)
			if !f.CanSet() {
				continue
			}

			switch f.Kind() {
			case reflect.String:
				f.SetString(val[0])
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				var intVal int64
				fmt.Sscanf(val[0], "%d", &intVal)
				f.SetInt(intVal)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				var uintVal uint64
				fmt.Sscanf(val[0], "%u", &uintVal)
				f.SetUint(uintVal)
			case reflect.Float32, reflect.Float64:
				var floatVal float64
				fmt.Sscanf(val[0], "%f", &floatVal)
				f.SetFloat(floatVal)
			case reflect.Bool:
				f.SetBool(val[0] == "true" || val[0] == "1" || val[0] == "on")
			}
		}
	}
	return nil
}
