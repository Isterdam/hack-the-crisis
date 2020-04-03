package api

import (
	"github.com/gin-gonic/gin"
)

func get_stores(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello world!",
	})
}