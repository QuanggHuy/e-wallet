package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nguyenhuy260301/e-wallet/internal/auth"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "thiếu Authorization header"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header phải có dạng: Bearer <token>"})
			c.Abort()
			return
		}

		claims, err := auth.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token không hợp lệ hoặc đã hết hạn"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("account_id", claims.AccountID)
		c.Next()
	}
}
