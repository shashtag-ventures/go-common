package netutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSafeRedirectURL(t *testing.T) {
	allowedBase := "https://example.com"

	tests := []struct {
		name        string
		redirect    string
		allowedBase string
		want        bool
		wantErr     bool
	}{
		{
			name:        "Relative path",
			redirect:    "/dashboard",
			allowedBase: allowedBase,
			want:        true,
			wantErr:     false,
		},
		{
			name:        "Same domain absolute",
			redirect:    "https://example.com/dashboard",
			allowedBase: allowedBase,
			want:        true,
			wantErr:     false,
		},
		{
			name:        "Subdomain",
			redirect:    "https://app.example.com/dashboard",
			allowedBase: allowedBase,
			want:        true,
			wantErr:     false,
		},
		{
			name:        "Different domain",
			redirect:    "https://evil.com",
			allowedBase: allowedBase,
			want:        false,
			wantErr:     false,
		},
		{
			name:        "Malicious subdomain attempt",
			redirect:    "https://example.com.evil.com",
			allowedBase: allowedBase,
			want:        false,
			wantErr:     false,
		},
		{
			name:        "Invalid redirect URL",
			redirect:    "://invalid",
			allowedBase: allowedBase,
			want:        false,
			wantErr:     true,
		},
		{
			name:        "Invalid base URL",
			redirect:    "/dashboard",
			allowedBase: "://invalid",
			want:        false,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsSafeRedirectURL(tt.redirect, tt.allowedBase)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
