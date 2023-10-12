package middleware

import (
	"grandatma_api/utility"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Validate(c *gin.Context) {
	reqHeader := c.GetHeader("Authorization")

	if reqHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Tidak ada token",
		})
		return
	}

	tokSlice := strings.Split(reqHeader, "Bearer ")

	if len(tokSlice) < 2 {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Invalid token",
		})
		return
	}

	tokString := tokSlice[1]

	token, err := utility.VerifyToken(tokString)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	userID := token.Id
	role := token.Role

	c.Set("userID", userID)
	c.Set("role", role)
	c.Next()

}
