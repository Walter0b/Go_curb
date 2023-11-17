package controllers

import (
	"Go_curb/Database/initializers"
	"Go_curb/tableTypes"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Retrieve air bookings with server-side pagination
func GetAllItems(c *gin.Context) {
	var airBookings []tableTypes.AirBooking

	// Retrieve page and pageSize from query parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
		return
	}

	// Calculate offset based on page and pageSize
	offset := (page - 1) * pageSize

	//retrieving air bookings here
	query := initializers.DB.Where("id_invoice IS NULL AND product_type = 'flight' AND transaction_type = 'sales' AND status = 'pending'")

	// Count total records (for pagination metadata)
	var totalRowCount int64
	if err := query.Model(&tableTypes.AirBooking{}).Count(&totalRowCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retrieve paginated air bookings
	if err := query.Limit(pageSize).Offset(offset).Find(&airBookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create a response object with paginated data and metadata
	response := gin.H{
		"data":          airBookings,                                                  // Data for the current page
		"totalRowCount": totalRowCount,                                                // Total count of records
		"currentPage":   page,                                                         // Current page
		"pageSize":      pageSize,                                                     // Page size
		"totalPages":    int((totalRowCount + int64(pageSize) - 1) / int64(pageSize)), // Total pages
	}

	c.JSON(http.StatusOK, response)
}
