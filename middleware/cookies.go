package middleware

import (
	"net/http"
	"time"
)

const (
	JWTCookieName           = "jwt_token"
	IsAuthenticatedCookie   = "is_authenticated"
	defaultCookieDuration   = 24 * time.Hour
)

// SetAuthCookies sets the JWT and authentication status cookies.
// It sets "jwt_token" as HttpOnly and "is_authenticated" as accessible by JS.
func SetAuthCookies(w http.ResponseWriter, token string, isSecure bool, domain string) {
	expiration := time.Now().Add(defaultCookieDuration)
	
	http.SetCookie(w, &http.Cookie{
		Name:     JWTCookieName,
		Value:    token,
		Path:     "/",
		Domain:   domain,
		Expires:  expiration,
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     IsAuthenticatedCookie,
		Value:    "true",
		Path:     "/",
		Domain:   domain,
		Expires:  expiration,
		HttpOnly: false,
		Secure:   isSecure,
		SameSite: http.SameSiteLaxMode,
	})
}

// ClearAuthCookies clears the JWT and authentication status cookies.
func ClearAuthCookies(w http.ResponseWriter, isSecure bool, domain string) {
	http.SetCookie(w, &http.Cookie{
		Name:     JWTCookieName,
		Value:    "",
		Path:     "/",
		Domain:   domain,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     IsAuthenticatedCookie,
		Value:    "",
		Path:     "/",
		Domain:   domain,
		Expires:  time.Unix(0, 0),
		HttpOnly: false,
		Secure:   isSecure,
		SameSite: http.SameSiteLaxMode,
	})
}
