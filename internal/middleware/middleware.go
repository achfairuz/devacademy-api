package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/iyuz/devacademy-api/internal/models"
	"github.com/iyuz/devacademy-api/internal/utils"
	"github.com/iyuz/devacademy-api/pkg/response"
)

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := parseClaims(c, secret)
		if claims == nil {
			response.Error(c, http.StatusUnauthorized, "unauthorized", "missing or invalid authorization header")
			c.Abort()
			return
		}
		setUserContext(c, claims)
		c.Next()
	}
}

func OptionalAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if claims := parseClaims(c, secret); claims != nil {
			setUserContext(c, claims)
		}
		c.Next()
	}
}

func parseClaims(c *gin.Context, secret string) *utils.Claims {
	header := c.GetHeader("Authorization")
	if header == "" {
		return nil
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return nil
	}

	claims, err := utils.ParseToken(secret, parts[1])
	if err != nil {
		return nil
	}
	return claims
}

func setUserContext(c *gin.Context, claims *utils.Claims) {
	c.Set("user_id", claims.UserID)
	c.Set("email", claims.Email)
	c.Set("role", claims.Role)
}

func RequireRole(roles ...models.Role) gin.HandlerFunc {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r.String()] = true
	}
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		roleStr, _ := role.(string)
		if !allowed[roleStr] {
			response.Error(c, http.StatusForbidden, "forbidden", "insufficient role permission")
			c.Abort()
			return
		}
		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return RequireRole(models.RoleAdmin)
}

func MentorOnly() gin.HandlerFunc {
	return RequireRole(models.RoleMentor)
}

func StudentOnly() gin.HandlerFunc {
	return RequireRole(models.RoleStudent)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
