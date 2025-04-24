package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hackshel/tracker-server/pkg/errs"
	"github.com/hackshel/tracker-server/utils"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error_code": errs.ERROR_AUTH_HEADER,
				"error_msg":  errs.GetMsg(errs.ERROR_AUTH_HEADER),
			})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error_code": errs.ERROR_AUTH_TOKEN,
				"error_msg":  errs.GetMsg(errs.ERROR_AUTH_TOKEN),
			})
			c.Abort()
			return
		}

		fmt.Printf("claims token UserID %v %T", claims.UserID, claims.UserID)
		c.Set("username", claims.Username)
		c.Set("expire_at", claims.ExpiresAt)
		c.Set("userRole", claims.UserRole)
		c.Set("userID", claims.UserID)
		c.Set("PassKey", claims.Passkey)

		c.Next()
	}
}
