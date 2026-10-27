package Middlewares

import (
	"SM/repositories/Utilities"
	"SM/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

const RolesKey = "allowed_roles"
const ClaimsKey = "claims"
const IsAnonymousRouteKey = "is_anonymous_route"

func AttachRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(RolesKey, roles)
		c.Next()
	}
}

func AllowAnonymous() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(IsAnonymousRouteKey, true)
		c.Next()
	}
}
func AuthorizeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAnonymousRoute := c.GetBool(IsAnonymousRouteKey)

		if isAnonymousRoute {
			c.Next()
			return
		}

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

		claims := claimsString.(*Utilities.IdentityClaim)
		if claims.Role == "" || claims.Role == "null" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, services.ForbiddenResult("Unauthorized"))
			return
		}

		for _, role := range allowedRoles {
			if role == claims.Role {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, services.ForbiddenResult("Unauthorized"))
	}
}
