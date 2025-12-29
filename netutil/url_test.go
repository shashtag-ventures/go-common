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

func TestGetCookieDomain(t *testing.T) {
	tests := []struct {
		name        string
		frontendURL string
		want        string
	}{
		{
			name:        "Localhost with port",
			frontendURL: "http://localhost:3000",
			want:        "",
		},
		{
			name:        "Localhost IP",
			frontendURL: "http://127.0.0.1:3000",
			want:        "",
		},
		{
			name:        "Production Domain",
			frontendURL: "https://karada.ai",
			want:        ".karada.ai",
		},
		{
			name:        "Production Domain with www",
			frontendURL: "https://www.karada.ai",
			want:        ".www.karada.ai",
		},
		{
			name:        "Subdomain",
			frontendURL: "https://app.example.com",
			want:        ".app.example.com",
		},
		{
			name:        "Invalid URL",
			frontendURL: ":/invalid",
			want:        "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCookieDomain(tt.frontendURL); got != tt.want {
				t.Errorf("GetCookieDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}
