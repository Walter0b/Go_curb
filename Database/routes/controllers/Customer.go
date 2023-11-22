package controllers

import (
	"Go_curb/Database/initializers"
	"Go_curb/tableTypes"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Retrieve all customers with pagination
func GetAllCustomer(c *gin.Context) {
	var customers []tableTypes.Customer
	embed := c.Param("embed")
	if embed != "" {
		// Use reflection to check if the specified association exists in Customer model
		if field, found := reflect.TypeOf(tableTypes.Customer{}).FieldByName(embed); found {
			// Check if the field is a struct (assumes it's an association)
			if field.Type.Kind() == reflect.Struct {
				if err := initializers.DB.Preload(embed).Find(&customers).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("%s not found for the given ID", embed)})
					return
				}

				// Combine association and customer information in the response
				response := gin.H{embed: customers}
				c.JSON(http.StatusOK, response)
				return
			}
		}

		// Handle the case where the specified association doesn't exist
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid association: %s", embed)})
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if err != nil || pageSize < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
		return
	}

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
//
//	func GetSpecificCustomer(c *gin.Context) {
//		id := c.Query("id")
//		customerID := tableTypes.Customer{}
//		if err := initializers.DB.First(&customerID, id).Error; err != nil {
//			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
//			return
//		}
//		c.JSON(http.StatusOK, customerID)
//	}
func GetSpecificCustomer(c *gin.Context) {
	id := c.Query("id")
	embedParam := c.Param("embed")
	customer := tableTypes.Customer{}

	// Check if the route pattern includes "/customers/"
	if embedParam != "" {

		if err := initializers.DB.Where("id = ?", id).Preload("Invoice").Find(&customer).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invoices not found for the given ID"})
			return
		}

		// Combine invoices and customer information in the response
		response := gin.H{"invoices": customer}
		c.JSON(http.StatusOK, response)
		return
	}
	if err := initializers.DB.First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}
	c.JSON(http.StatusOK, customer)
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
