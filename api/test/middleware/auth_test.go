package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adityawiryaa/api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func TestAPIKeyAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		apiKey     string
		headerKey  string
		wantStatus int
	}{
		{
			name:       "valid key",
			apiKey:     "test-secret",
			headerKey:  "test-secret",
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing key",
			apiKey:     "test-secret",
			headerKey:  "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "invalid key",
			apiKey:     "test-secret",
			headerKey:  "wrong-key",
			wantStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(middleware.APIKeyAuth(tt.apiKey))
			r.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.headerKey != "" {
				req.Header.Set("X-API-Key", tt.headerKey)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}
