package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"lab2/internal/auth"
)

const AccountIDKey = "account_id"

// RequireJWT validates the Bearer token and checks the expected purpose claim.
func RequireJWT(secret []byte, purpose string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")
		if tokenString == header {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			return
		}

		claims, err := auth.ParseToken(secret, tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		if claims.Purpose != purpose {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "token purpose mismatch"})
			return
		}

		c.Set(AccountIDKey, claims.AccountID)
		c.Next()
	}
}
