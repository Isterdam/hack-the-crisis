package main

import (
	"fmt"
	"log"

	"github.com/Isterdam/hack-the-crisis-backend/src/api"
	"github.com/Isterdam/hack-the-crisis-backend/src/db"
	_ "github.com/Isterdam/hack-the-crisis-backend/src/docs"
	"github.com/Isterdam/hack-the-crisis-backend/src/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title ShopAlone API
// @version 1.0
// @description Swagger API for ShopAlone API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email ad5880mu-s@student.lu.se

// @BasePath /
func main() {
	r := gin.Default()

	sql, err := db.InitDB()

	if err != nil {
		fmt.Printf("%s", err)
		log.Fatal(err)
	}
	config := cors.DefaultConfig()
	config.AllowOriginFunc = func(origin string) bool {
		return true
	}
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))
	//r.Use(cors.Default())

	r.Use(func(c *gin.Context) {
		c.Set("db", sql)
	})

	// localhost:8080/swagger/index.html to access documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api.InitializeCache()
	handlers.InitPublicRoutes(r)
	handlers.InitCompanyRoutes(r)

	r.Run(":8080")
}
