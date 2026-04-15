package auth

import (
	"net/http"
	"strings"

	"github.com/fixia/golang-api/utils"
	"github.com/gin-gonic/gin"
)

const (
	CtxUserID = "userID"
	CtxEmail  = "email"
	CtxRole   = "role"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "missing Authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			utils.ErrorResponse(c, http.StatusUnauthorized, "invalid Authorization header format")
			c.Abort()
			return
		}

		claims, err := ValidateToken(parts[1])
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxEmail, claims.Email)
		c.Set(CtxRole, claims.Role)
		c.Next()
	}
}

func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		r, exists := c.Get(CtxRole)
		if !exists || r != role {
			utils.ErrorResponse(c, http.StatusForbidden, "forbidden: insufficient role")
			c.Abort()
			return
		}
		c.Next()
	}
}
