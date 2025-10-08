package testutil

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// FindCookie is a helper function to find a cookie by name from a slice of cookies.
func FindCookie(cookies []*http.Cookie, name string) *http.Cookie {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}

// CreateRequestWithJWT is a helper function to create an HTTP request with a JWT cookie.
func CreateRequestWithJWT(token string) *http.Request {
	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "jwt_token",
		Value: token,
	})
	return req
}

// CreateTempCSV is a helper function to create a temporary CSV file for testing.
func CreateTempCSV(t *testing.T, dir, filename, content string) string {
	filePath := filepath.Join(dir, filename)
	err := os.WriteFile(filePath, []byte(content), 0644)
	assert.NoError(t, err)
	return filePath
}
