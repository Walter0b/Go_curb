package controllers

import (
	"Go_curb/Database/initializers"
	"Go_curb/tableTypes"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Retrieve all customers with pagination
func GetAllCustomer(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
		return
	}

	var customers []tableTypes.Customer
	var totalRowCount int64 // Total count of records

	// Count total records
	if err := initializers.DB.Model(&tableTypes.Customer{}).Count(&totalRowCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	offset := (page - 1) * pageSize

	if err := initializers.DB.Limit(pageSize).Offset(offset).Find(&customers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"data":          customers,     // Data for the current page
		"totalRowCount": totalRowCount, // Total count of records
	}

	c.JSON(http.StatusOK, response)
}

// Retrieve a specific customer by ID
func GetSpecificCustomer(c *gin.Context) {
	id := c.Param("id")
	customerID := tableTypes.Customer{}
	if err := initializers.DB.First(&customerID, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}
	c.JSON(http.StatusOK, customerID)
}

// Update a specific customer by ID
func UpdateCustomer(c *gin.Context) {
	id := c.Param("id")
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

// Create a new customer
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

// Delete a specific customer by ID
func DeleteCutomer(c *gin.Context) {

	id := c.Param("id")
	if err := initializers.DB.Where("id = ?", id).Delete(&tableTypes.Customer{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}
