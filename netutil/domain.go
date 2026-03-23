package netutil

import "regexp"

// validSubdomainRe validates the MCP vanity subdomain prefix.
// Only lowercase alphanumeric characters and hyphens, no leading/trailing hyphens.
var validSubdomainRe = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)

// IsValidSubdomain checks if a string is a valid subdomain prefix.
// It allows lowercase alphanumeric characters and hyphens, but no leading or trailing hyphens.
func IsValidSubdomain(subdomain string) bool {
	return validSubdomainRe.MatchString(subdomain)
}
