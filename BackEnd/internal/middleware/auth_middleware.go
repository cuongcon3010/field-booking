package middleware

import (
	"net/http"
	"strings"

	"field-booking/backend/pkg/jwt"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserIDKey  = "user_id"
	ContextEmailKey   = "email"
	AccessTokenCookie = "access_token"
)

func AuthMiddleware(tokenService jwt.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// 1. Ưu tiên lấy token từ HttpOnly cookie
		if cookie, err := c.Cookie(AccessTokenCookie); err == nil && cookie != "" {
			tokenString = cookie
		}

		// 2. Dự phòng: lấy từ header Authorization: Bearer <token>
		if tokenString == "" {
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Authentication required",
			})
			c.Abort()
			return
		}

		claims, err := tokenService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextEmailKey, claims.Email)
		c.Next()
	}
}
