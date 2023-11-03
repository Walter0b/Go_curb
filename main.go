package main

import (
	"Go_curb/Database/dbConnect"
	"Go_curb/Database/initializers"
	"Go_curb/Database/routes"
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
		return
	}

	// Migrate the GORM models
	db.AutoMigrate(&tableTypes.Customer{})

	// Pass the GORM DB instance to the routes
	routes.CustomerRoutes(r, db)

	r.Run(":8080")
}
