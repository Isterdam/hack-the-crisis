package main

import (
	"github.com/Isterdam/hack-the-crisis-backend/src/handlers/public"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	public.init_public_handlers(r)

	// r.GET("/", root)
	// r.POST("/register", register_company)

	r.Run(":8080")
}

/*
func root(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello world!",
	})
}
*/
