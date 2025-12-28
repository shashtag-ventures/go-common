package strutil_test

import (
	"testing"

	"github.com/shashtag-ventures/go-common/strutil"
	"github.com/stretchr/testify/assert"
)

func TestSlugify(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello-world"},
		{"  Hello   World  ", "hello-world"},
		{"Hello@World!", "hello-world"},
		{"My Awesome Project 123", "my-awesome-project-123"},
		{"---Leading and Trailing---", "leading-and-trailing"},
		{"UPPER CASE", "upper-case"},
		{"Complex !@#$%^&*() Characters", "complex-characters"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, strutil.Slugify(tt.input))
		})
	}
}
