package controllers

import (
	"Go_curb/Database/components"
	"Go_curb/Database/initializers"
	"Go_curb/tableTypes"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Retrieve all customers with pagination
func GetAllCustomer(c *gin.Context) {

	id := c.Query("id")
	var customer []tableTypes.Customer

	query := initializers.DB.Model(&tableTypes.Customer{}).Order("Customer_name ASC")
	embedField := c.Query("embed")
	components.Get(c, query, &customer, embedField, id)
}

func CreateCustomer(c *gin.Context) {
	var customer tableTypes.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := initializers.DB.Create(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, customer)
}

// Update a specific customer by ID
func UpdateCustomer(c *gin.Context) {
	id := c.Query("id")
	customerID := tableTypes.Customer{}
	if err := initializers.DB.First(&customerID, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	if err := c.ShouldBindJSON(&customerID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := initializers.DB.Save(&customerID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, customerID)
}

// Delete a specific customer by ID
func DeleteCutomer(c *gin.Context) {

	id := c.Query("id")
	if err := initializers.DB.Where("id = ?", id).Delete(&tableTypes.Customer{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}
