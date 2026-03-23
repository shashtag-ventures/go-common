package strutil

import (
	"regexp"
	"strings"

	"github.com/google/uuid"
)

var slugRegex = regexp.MustCompile(`[^a-z0-9]+`)

// Slugify converts a string into a URL-friendly slug.
// It lowercase the string, replaces non-alphanumeric characters with hyphens,
// and removes leading/trailing hyphens.
func Slugify(s string) string {
	s = strings.ToLower(s)
	s = slugRegex.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

// RandomSuffix returns a short (5-character) random string based on a UUID.
// This is useful for creating unique slugs or identifiers when collisions occur.
func RandomSuffix() string {
	return uuid.New().String()[:5]
}
