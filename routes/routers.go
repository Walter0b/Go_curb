package routes

import (
	"Go_curb/tableTypes"
	// "log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CustomerRoutes(r *gin.Engine, db *gorm.DB) {
	r.GET("/", func(c *gin.Context) {
		var customers []tableTypes.Customer
		if err := db.Find(&customers).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, customers)
	})

	r.POST("/customers", func(c *gin.Context) {
		var customer tableTypes.Customer
		if err := c.ShouldBindJSON(&customer); err != nil {
			// log.Println(customer)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Create(&customer).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, customer)
	})
	
	r.DELETE("/customers/:id", func(c *gin.Context) {
		customerID := c.Param("id")

		if err := db.Where("id = ?", customerID).Delete(&tableTypes.Customer{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	
		c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
	})
	
	

}
