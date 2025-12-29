package netutil

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// IsSafeRedirectURL checks if a redirect URL is safe to use.
// It allows relative paths or absolute URLs that match the allowed base URL's hostname (or subdomains).
func IsSafeRedirectURL(redirectURL, allowedBaseURL string) (bool, error) {
	parsedRedirect, err := url.Parse(redirectURL)
	if err != nil {
		return false, fmt.Errorf("invalid redirect URL: %w", err)
	}

	parsedBase, err := url.Parse(allowedBaseURL)
	if err != nil {
		return false, fmt.Errorf("invalid allowed base URL: %w", err)
	}

	// Allow relative paths (e.g., "/dashboard")
	if !parsedRedirect.IsAbs() {
		return true, nil
	}

	// Allow absolute URLs to the trusted domain (including subdomains)
	// e.g. if allowedBaseURL is "example.com", then "app.example.com" is allowed.
	// We verify that the redirect hostname ends with the allowed hostname.
	// This simple check works for most cases but be aware of "evil.com/example.com" if not parsing correctly.
	// url.Parse handles the host extraction reliably.
	allowedHost := parsedBase.Hostname()
	redirectHost := parsedRedirect.Hostname()

	if redirectHost == allowedHost || strings.HasSuffix(redirectHost, "."+allowedHost) {
		return true, nil
	}

	return false, nil
}

// GetCookieDomain extracts the domain for cookie setting from a URL string.
// It handles localhost correctly (returning empty string) and prefixes root domains with a dot.
func GetCookieDomain(frontendURL string) string {
	if parsedURL, err := url.Parse(frontendURL); err == nil {
		hostname := parsedURL.Hostname()

		// If hostname is an IP address, return empty to let browser handle it
		if net.ParseIP(hostname) != nil {
			return ""
		}

		// Avoid setting domain for localhost to allow browser default behavior
		if strings.Contains(hostname, ".") && !strings.HasSuffix(hostname, "localhost") {
			if !strings.HasPrefix(hostname, ".") {
				return "." + hostname
			}
			return hostname
		}
	}
	return ""
}
