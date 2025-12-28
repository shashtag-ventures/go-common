package strutil

import (
	"regexp"
	"strings"
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
