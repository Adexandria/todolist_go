package Middlewares

import (
	"SM/repositories/Utilities"
	"SM/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthenticateMiddleware(tokenManager Utilities.ITokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, services.NotAuthorizedResult("Unauthorized"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, services.NotAuthorizedResult("Unauthorized"))
			return
		}

		tokenString := parts[1]

		claims, err := tokenManager.DecodeToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, services.NotAuthorizedResult("Unauthorized"))
			return
		}

		c.Set(ClaimsKey, claims)
		c.Next()
	}
}
