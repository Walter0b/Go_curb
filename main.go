package main

import (
	"Go_curb/dbConnect"
	"Go_curb/initializers"
	"Go_curb/routes"
	"Go_curb/tableTypes"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
}

func main() {
	r := gin.Default()

	// Initialize the GORM database connection using the dbConnect package
	db, err := dbConnect.InitDB()
	if err != nil {
		// Handle the error as needed
		// You can log the error and exit the application, for example
		return
	}

	// Migrate the GORM models
	db.AutoMigrate(&tableTypes.Customer{})

	// Pass the GORM DB instance to your routes
	routes.CustomerRoutes(r, db)

	r.Run(":8080")
}
