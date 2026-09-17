package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware(frontendURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Chỉ cho phép origin được chỉ định (hoặc origin của request nếu khớp frontendURL)
		if origin != "" && (origin == frontendURL || frontendURL == "*") {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		} else if frontendURL != "" && frontendURL != "*" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", frontendURL)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
