package rest

import (
	"net/http"
	"strings"

	"slices"

	"github.com/charmingruby/doris/lib/security"
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	jwt *security.JWT
}

func NewMiddleware(jwt *security.JWT) *Middleware {
	return &Middleware{
		jwt: jwt,
	}
}

func (m *Middleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})

			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format",
			})

			return
		}

		token := parts[1]

		sub, payload, err := m.jwt.Validate(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		c.Set("api-key-id", sub)
		c.Set("tier", payload.Tier)

		c.Next()
	}
}

func (m *Middleware) RBAC(allowedTiers ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})

			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format",
			})

			return
		}

		token := parts[1]

		sub, payload, err := m.jwt.Validate(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})

			return
		}

		hasAccess := slices.Contains(allowedTiers, payload.Tier)

		if !hasAccess {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "insufficient permissions",
			})

			return
		}

		c.Set("api-key-id", sub)
		c.Set("tier", payload.Tier)

		c.Next()
	}
}
