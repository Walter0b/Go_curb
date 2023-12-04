package main

import (
	"Go_curb/Database/initializers"
	"Go_curb/Database/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDatabase()

}

func main() {
	r := gin.Default()

	// Use the CORS middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true

	r.Use(cors.New(config))

	routes.Api(r)
	r.Run()
}
