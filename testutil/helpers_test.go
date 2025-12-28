package testutil

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindCookie(t *testing.T) {
	cookies := []*http.Cookie{
		{Name: "session", Value: "val1"},
		{Name: "jwt_token", Value: "val2"},
	}

	// Found
	cookie := FindCookie(cookies, "jwt_token")
	assert.NotNil(t, cookie)
	assert.Equal(t, "val2", cookie.Value)

	// Not found
	assert.Nil(t, FindCookie(cookies, "missing"))
}

func TestCreateRequestWithJWT(t *testing.T) {
	req := CreateRequestWithJWT("test-token")
	cookie, err := req.Cookie("jwt_token")
	assert.NoError(t, err)
	assert.Equal(t, "test-token", cookie.Value)
}

func TestCreateTempCSV(t *testing.T) {
	dir := t.TempDir()
	filename := "test.csv"
	content := "header1,header2\nval1,val2"

	path := CreateTempCSV(t, dir, filename, content)
	
	// Verify file exists and has content
	data, err := os.ReadFile(path)
	assert.NoError(t, err)
	assert.Equal(t, content, string(data))
}
