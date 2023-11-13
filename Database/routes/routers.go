package routes

import (
	"Go_curb/tableTypes"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CustomerRoutes(r *gin.Engine, db *gorm.DB) {

	// GET /customers?page=1&pageSize=10 - Retrieve all customers with pagination
	r.GET("/customers", func(c *gin.Context) {
		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		page += page
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
		if err := db.Model(&tableTypes.Customer{}).Count(&totalRowCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		offset := (page - 1) * pageSize

		if err := db.Limit(pageSize).Offset(offset).Find(&customers).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := gin.H{
			"data":          customers,     // Data for the current page
			"totalRowCount": totalRowCount, // Total count of records
		}

		c.JSON(http.StatusOK, response)
	})

	// GET /customers/:id - Retrieve a specific customer by ID
	r.GET("/customers/:id", func(c *gin.Context) {
		id := c.Param("id")
		customerID := tableTypes.Customer{}
		if err := db.First(&customerID, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			return
		}
		c.JSON(http.StatusOK, customerID)
	})

	// PUT /customers/:id - Update a specific customer by ID
	r.PUT("/customers/:id", func(c *gin.Context) {
		id := c.Param("id")
		customerID := tableTypes.Customer{}
		if err := db.First(&customerID, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			return
		}

		if err := c.ShouldBindJSON(&customerID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Save(&customerID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, customerID)
	})

	// POST /customers - Create a new customer
	r.POST("/customers", func(c *gin.Context) {
		var customer tableTypes.Customer
		if err := c.ShouldBindJSON(&customer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Create(&customer).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, customer)
	})

	// DELETE /customers/:id - Delete a specific customer by ID
	r.DELETE("/customers/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := db.Where("id = ?", id).Delete(&tableTypes.Customer{}).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
	})
	r.GET("/currencies", func(c *gin.Context) {
		var currencies []tableTypes.Currency
		if err := db.Find(&currencies).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, currencies)
	})

	// GET / - Retrieve all Id from countries
	r.GET("/countries", func(c *gin.Context) {
		var Countries []tableTypes.Country
		if err := db.Find(&Countries).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, Countries)
	})
}
