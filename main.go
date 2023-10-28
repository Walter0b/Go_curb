package main

import (
    "github.com/gin-gonic/gin"
    "gocrub/routes"
    "gocrub/initializers"
    "gocrub/dbConnect"
    "gocrub/tableTypes"
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
