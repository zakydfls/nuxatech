package utils

import "github.com/gin-gonic/gin"

func GetUserID(c *gin.Context) (userID string) {
	return c.MustGet("user_id").(string)
}

func GetEmail(c *gin.Context) (email string) {
	return c.MustGet("email").(string)
}
