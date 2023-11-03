package main

import (
	"Go_curb/Database/dbConnect"
	"Go_curb/Database/initializers"
	"Go_curb/Database/routes"
	"Go_curb/tableTypes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
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

	// Initialize the GORM database connection using the dbConnect package
	db, err := dbConnect.InitDB()
	if err != nil {
		// Handle the error as needed
		return
	}

	// Migrate the GORM models
	db.AutoMigrate(&tableTypes.Customer{})

	// Pass the GORM DB instance to the routes
	routes.CustomerRoutes(r, db)

	r.Run(":8080")
}
