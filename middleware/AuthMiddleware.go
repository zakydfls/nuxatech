package middleware

import (
	"nuxatech-nextmedis/service"
	"strings"

	"github.com/gin-gonic/gin"
)

var authService service.AuthService

func SetAuthService(service service.AuthService) {
	authService = service
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header is required"})
			return
		}

		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) != 2 {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token format"})
			return
		}

		tokenString := splitToken[1]
		payload, err := authService.ValidateAccessToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}

		c.Set("user_id", payload.UserID)
		c.Set("username", payload.Username)
		c.Set("email", payload.Email)

		c.Next()
	}
}
