package Middlewares

import (
	"SM/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const RolesKey = "allowed_roles"
const ClaimsKey = "claims"

func AuthorizeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedRolesString, exists := c.Get(RolesKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, services.ForbiddenResult("Unauthorized"))
			return

		}
		allowedRoles := allowedRolesString.([]string)

		claimsString, exists := c.Get(ClaimsKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, services.ForbiddenResult("Unauthorized"))
			return
		}

		claims := claimsString.(jwt.MapClaims)

		userRole, ok := claims["roles"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, services.ForbiddenResult("Unauthorized"))
			return
		}

		for _, role := range allowedRoles {
			if role == userRole {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, services.ForbiddenResult("Unauthorized"))
	}
}
