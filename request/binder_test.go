package request_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shashtag-ventures/go-common/request"
	"github.com/stretchr/testify/assert"
)

type TestUser struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Age   int    `json:"age" validate:"gte=0,lte=130"`
}

func TestDecodeAndValidate(t *testing.T) {
	// Test Case 1: Valid request
	t.Run("Valid Request", func(t *testing.T) {
		user := TestUser{
			Name:  "John Doe",
			Email: "john@example.com",
			Age:   30,
		}
		body, _ := json.Marshal(user)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		var decodedUser TestUser
		err := request.DecodeAndValidate(req, &decodedUser)

		assert.NoError(t, err)
		assert.Equal(t, user, decodedUser)
	})

	// Test Case 2: Invalid JSON
	t.Run("Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("{invalid json}"))

		var decodedUser TestUser
		err := request.DecodeAndValidate(req, &decodedUser)

		assert.Error(t, err)
	})

	// Test Case 3: Validation Failure (missing required field)
	t.Run("Validation Failure - Missing Field", func(t *testing.T) {
		user := map[string]interface{}{
			"email": "john@example.com",
			"age":   30,
		}
		body, _ := json.Marshal(user)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		var decodedUser TestUser
		err := request.DecodeAndValidate(req, &decodedUser)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Key: 'TestUser.Name' Error:Field validation for 'Name' failed on the 'required' tag")
	})

	// Test Case 4: Validation Failure (invalid email)
	t.Run("Validation Failure - Invalid Email", func(t *testing.T) {
		user := TestUser{
			Name:  "John Doe",
			Email: "invalid-email",
			Age:   30,
		}
		body, _ := json.Marshal(user)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		var decodedUser TestUser
		err := request.DecodeAndValidate(req, &decodedUser)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Key: 'TestUser.Email' Error:Field validation for 'Email' failed on the 'email' tag")
	})
}
