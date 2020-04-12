package middleware

import (
	"net/http"

	"github.com/Isterdam/hack-the-crisis-backend/src/auth"
	"github.com/gin-gonic/gin"
)

func AuthRequired(c *gin.Context) {

	token := c.Request.Header.Get("Authorization")
	var claim auth.Claims
	isAuth := auth.IsValidToken(token, &claim)

	if !isAuth {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
	} else {
		c.Set("id", claim.ID)
	}
}
