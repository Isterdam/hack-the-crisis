package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func PrivateAuthenticationHandler(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token != os.Getenv("PRIVATE_API_KEY") {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "Unauthorized",
		})
	}
}
