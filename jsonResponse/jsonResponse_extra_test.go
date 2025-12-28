package jsonResponse

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	customErrors "github.com/shashtag-ventures/go-common/errors"
	"github.com/stretchr/testify/assert"
)

func TestSendAutoErrorResponse(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{
			name:       "Not Found Error",
			err:        customErrors.ErrNotFound,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Unauthorized Error",
			err:        customErrors.ErrUnauthorized,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Internal Error",
			err:        customErrors.ErrInternal,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "Unmapped Generic Error",
			err:        errors.New("something went wrong"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			SendAutoErrorResponse(w, tt.err)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
