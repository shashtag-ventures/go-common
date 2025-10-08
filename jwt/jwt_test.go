package jwt_test

import (
	"testing"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/shashtag-ventures/go-common/jwt"
	"github.com/stretchr/testify/assert"
)

const (
	testSecret    = "supersecretjwtkey"
	wrongSecret = "wrongsecret"
)

func TestCreateToken(t *testing.T) {
	userID := uint(123)
	role := "user"

	tokenString, err := jwt.CreateToken(userID, role, testSecret)

	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Parse the token to verify claims
	claims, err := jwt.ParseToken(tokenString, testSecret)
	assert.NoError(t, err)
	assert.NotNil(t, claims)

	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, role, claims.Role)
	assert.True(t, claims.ExpiresAt.Time.After(time.Now()))
}

func TestParseToken(t *testing.T) {
	userID := uint(456)
	role := "admin"

	// Test Case 1: Successful parsing of a valid token
	t.Run("Valid Token", func(t *testing.T) {
		tokenString, err := jwt.CreateToken(userID, role, testSecret)
		assert.NoError(t, err)

		claims, err := jwt.ParseToken(tokenString, testSecret)
		assert.NoError(t, err)
		assert.NotNil(t, claims)

		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, role, claims.Role)
		assert.True(t, claims.ExpiresAt.Time.After(time.Now()))
	})

	// Test Case 2: Invalid token (wrong secret)
	t.Run("Invalid Secret", func(t *testing.T) {
		tokenString, err := jwt.CreateToken(userID, role, testSecret)
		assert.NoError(t, err)

		claims, err := jwt.ParseToken(tokenString, wrongSecret)
		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Contains(t, err.Error(), "signature is invalid")
	})

	// Test Case 3: Expired token
	t.Run("Expired Token", func(t *testing.T) {
		// Create a token that expires in the past
		oldClaims := &jwt.Claims{
			UserID: userID,
			Role:   role,
			RegisteredClaims: gojwt.RegisteredClaims{
				ExpiresAt: gojwt.NewNumericDate(time.Now().Add(-time.Hour)), // Expired an hour ago
			},
		}
		token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, oldClaims)
		tokenString, err := token.SignedString([]byte(testSecret))
		assert.NoError(t, err)

		claims, err := jwt.ParseToken(tokenString, testSecret)
		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Contains(t, err.Error(), "token is expired")
	})

	// Test Case 4: Malformed token
	t.Run("Malformed Token", func(t *testing.T) {
		malformedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
		claims, err := jwt.ParseToken(malformedToken, testSecret)
		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Contains(t, err.Error(), "token signature is invalid")
	})
}
