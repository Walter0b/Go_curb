package controllers

import (
	"Go_curb/Database/components"
	"Go_curb/Database/initializers"
	"Go_curb/tableTypes"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func GetAllPayments(c *gin.Context) {
	var payments []tableTypes.PaymentReceived

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

	// Fetch invoices from the database
	if err := initializers.DB.Limit(pageSize).Offset(offset).Find(&payments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create a response object with paginated data and metadata
	response := gin.H{
		"data":          payments,                                  // Data for the current page
		"totalRowCount": len(payments),                             // Total count of records
		"currentPage":   page,                                      // Current page
		"pageSize":      pageSize,                                  // Page size
		"totalPages":    (len(payments) + pageSize - 1) / pageSize, // Total pages
	}

	c.JSON(http.StatusOK, response)
}

func CreatePayments(c *gin.Context) {
	// Bind the JSON payload to the PaymentReceived struct
	var paymentReceived tableTypes.PaymentReceived
	if err := c.ShouldBindJSON(&paymentReceived); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default values
	paymentReceived.Date = time.Now()
	paymentReceived.Status = "open"
	paymentReceived.Balance = paymentReceived.Amount
	paymentReceived.Tag = "2"
	paymentReceived.Slug = components.GenerateRandomSlug()
	// Validate and save to the database
	if err := initializers.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "slug"}},
		DoUpdates: clause.AssignmentColumns(
			[]string{"number", "date", "balance", "amount", "currency_rate", "fop", "reference", "deducted_tax", "note",
				"used_amount", "status", "base_amount", "is_reconciled", "id_customer", "id_chart_of_accounts_from", "type",
				"transfert_type", "already_used", "receipiant_name", "tag"},
		),
	}).Create(&paymentReceived).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with the created payment
	c.JSON(http.StatusCreated, paymentReceived)
}

// PUT /payments/:id

// DELETE /payments/:id
func DeletePayments(c *gin.Context) {
	id := c.Param("id")

	// Retrieve the payment
	var payment tableTypes.PaymentReceived
	if err := initializers.DB.First(&payment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	// Check if the payment has been used
	if payment.Status == "used" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete a used payment"})
		return
	}

	// Consider additional checks and logic here based on your requirements

	// Delete the payment
	if err := initializers.DB.Delete(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment deleted successfully"})
}

// GET /payments/:id - Retrieve a specific customer by ID
func GetSpecifiPayments(c *gin.Context) {
	id := c.Param("id")
	customerID := tableTypes.Customer{}
	if err := initializers.DB.First(&customerID, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}
	c.JSON(http.StatusOK, customerID)
}