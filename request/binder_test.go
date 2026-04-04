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
	Name    string `json:"name" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
	Age     int    `json:"age" validate:"gte=0,lte=130"`
	Slug    string `json:"slug" validate:"omitempty,slug"`
	RepoURL string `json:"repo_url" validate:"omitempty,git-url"`
}

func TestDecodeAndValidate(t *testing.T) {
	// Test Case 1: Valid request
	t.Run("Valid Request", func(t *testing.T) {
		user := TestUser{
			Name:    "John Doe",
			Email:   "john@example.com",
			Age:     30,
			Slug:    "john-doe-123",
			RepoURL: "https://github.com/user/repo",
		}
		body, _ := json.Marshal(user)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		var decodedUser TestUser
		err := request.DecodeAndValidate(req, &decodedUser)

		assert.NoError(t, err)
		assert.Equal(t, user, decodedUser)
	})

	// Test Case 2: Invalid JSON
	t.Run("Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("{invalid json}"))
		req.Header.Set("Content-Type", "application/json")

		var decodedUser TestUser
		err := request.DecodeAndValidate(req, &decodedUser)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid json")
	})

	// Test Case 3: Strict JSON (unknown fields)
	t.Run("Strict JSON - Unknown Fields", func(t *testing.T) {
		body := `{"name":"John","email":"john@example.com","unknown":"field"}`
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")

		var decodedUser TestUser
		err := request.DecodeAndValidate(req, &decodedUser)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "json: unknown field \"unknown\"")
	})

	// Test Case 4: Custom Validation - Slug
	t.Run("Validation Failure - Invalid Slug", func(t *testing.T) {
		user := TestUser{
			Name:  "John",
			Email: "john@example.com",
			Slug:  "Invalid_Slug!",
		}
		body, _ := json.Marshal(user)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		var decodedUser TestUser
		err := request.DecodeAndValidate(req, &decodedUser)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed on the 'slug' tag")
	})

	// Test Case 5: Custom Validation - Git URL
	t.Run("Validation Failure - Invalid Git URL", func(t *testing.T) {
		user := TestUser{
			Name:    "John",
			Email:   "john@example.com",
			RepoURL: "ftp://github.com/user/repo", // Only http/https allowed
		}
		body, _ := json.Marshal(user)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		var decodedUser TestUser
		err := request.DecodeAndValidate(req, &decodedUser)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed on the 'git-url' tag")
	})

	// Test Case 6: Form Data
	t.Run("Form Data", func(t *testing.T) {
		form := "name=John+Doe&email=john@example.com&age=30&slug=john-doe"
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(form))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		var decodedUser TestUser
		err := request.DecodeAndValidate(req, &decodedUser)

		assert.NoError(t, err)
		assert.Equal(t, "John Doe", decodedUser.Name)
		assert.Equal(t, "john@example.com", decodedUser.Email)
		assert.Equal(t, 30, decodedUser.Age)
		assert.Equal(t, "john-doe", decodedUser.Slug)
	})
}

func TestBinder_New(t *testing.T) {
	b := request.New()
	assert.NotNil(t, b)

	user := TestUser{
		Name:  "John",
		Email: "invalid",
	}
	err := b.Validate(user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed on the 'email' tag")
}
