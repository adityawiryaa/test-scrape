package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/adityawiryaa/api/pkg/response"
)

func APIKeyAuth(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")
		if key == "" {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "missing API key")
			c.Abort()
			return
		}
		if key != apiKey {
			response.Error(c, http.StatusForbidden, "FORBIDDEN", "invalid API key")
			c.Abort()
			return
		}
		c.Next()
	}
}
