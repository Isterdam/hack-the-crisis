package main

import (
	"fmt"

	"github.com/Isterdam/hack-the-crisis-backend/src/api"
	"github.com/Isterdam/hack-the-crisis-backend/src/db"
	"github.com/Isterdam/hack-the-crisis-backend/src/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	r := gin.Default()
	
	sql, err := db.InitDB()

	if err != nil {
		fmt.Printf("%s", err)
	}

	r.Use(cors.Default())

	r.Use(func(c *gin.Context) {
		c.Set("db", sql)
	})

	api.Initialize_constants()
	handlers.Init_public_routes(r)
	handlers.Init_company_routes(r)

	r.GET("/", root)

	r.Run(":8080")
}

func root(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello world!",
	})
}
