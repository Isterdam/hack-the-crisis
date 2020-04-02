package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", root)
	// r.POST("/register", register_company)

	r.Run(":8080")
}

func root(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello world!",
	})
}
