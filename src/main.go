package main

import (
	"github.com/Isterdam/hack-the-crisis-backend/src/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// handlers.Init_public_routes(r)
	// handlers.Init_company_routes(r)

	r.GET("/", root)

	r.Run(":8080")
}

/*
func root(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello world!",
	})
}
*/
