package controllers

import (
	"Go_curb/Database/components"
	"Go_curb/Database/initializers"
	"Go_curb/tableTypes"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Retrieve all invoices with pagination
func GetAllInvoices(c *gin.Context) {
	var invoices []tableTypes.Invoice

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
	if err := initializers.DB.Limit(pageSize).Offset(offset).Find(&invoices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert invoices to the response format
	// var invoiceResponses []InvoiceResponse
	// for _, invoice := range invoices {
	// 	invoiceResponses = append(invoiceResponses, convertInvoiceToResponse(invoice))
	// }

	// Create a response object with paginated data and metadata
	response := gin.H{
		"data":          invoices,                                  // Data for the current page
		"totalRowCount": len(invoices),                             // Total count of records
		"currentPage":   page,                                      // Current page
		"pageSize":      pageSize,                                  // Page size
		"totalPages":    (len(invoices) + pageSize - 1) / pageSize, // Total pages
	}

	c.JSON(http.StatusOK, response)
}

// Get Invoice by ID Handler
func GetSpecificInvoicesHandler(c *gin.Context) {
	id := c.Query("id")
	var invoices []tableTypes.Invoice

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
	db := initializers.DB.Where(map[string]interface{}{"id_customer": id, "balance": "0"}).Limit(pageSize).Offset(offset)
	if err := db.Find(&invoices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"data":          invoices,                                  // Data for the current page
		"totalRowCount": len(invoices),                             // Total count of records
		"currentPage":   page,                                      // Current page
		"pageSize":      pageSize,                                  // Page size
		"totalPages":    (len(invoices) + pageSize - 1) / pageSize, // Total pages
	}

	c.JSON(http.StatusOK, response)
}

func GetSpecificInvoice(c *gin.Context) {
	idParam := c.Query("id")
	invoices := tableTypes.Invoice{}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Check if the route pattern includes "/customers/"
	if strings.Contains(c.FullPath(), "/customers/") {

		if err := initializers.DB.Where("id = ?", id).Preload("Customer").Find(&invoices).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invoices not found for the given ID"})
			return
		}

		// Combine invoices and customer information in the response
		response := gin.H{"invoices": invoices}
		c.JSON(http.StatusOK, response)
		return
	}
	if err := initializers.DB.First(&invoices, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}
	c.JSON(http.StatusOK, invoices)
}

func CreateInvoice(c *gin.Context) {
	var requestPayload tableTypes.RequestPayload

	// Bind the JSON payload to the requestPayload struct
	if err := c.ShouldBindJSON(&requestPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i := range requestPayload.TravelItems {
		requestPayload.TravelItems[i].Status = "invoiced"
	}

	// Parse DueDate
	dueDate, err := time.Parse("2006-01-02", requestPayload.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid due date format"})
		return
	}

	// Start a database transaction
	tx := initializers.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return
	}

	// Defer a function to handle the transaction commit or rollback
	defer func() {
		if r := recover(); r != nil {
			// Rollback the transaction on panic
			tx.Rollback()
		}
	}()

	// Create an Invoice
	newInvoice := tableTypes.Invoice{
		CreationDate:  time.Now(),
		InvoiceNumber: components.GenerateUniqueInvoiceNumber(),
		IDCustomer:    requestPayload.IDCustomer,
		DueDate:       dueDate,
		Amount:        requestPayload.Amount,
		Status:        "unpaid",
		Tag:           requestPayload.Tag,
	}
	newInvoice.Balance = requestPayload.Amount

	// Save the invoice to the database
	if err := tx.Create(&newInvoice).Error; err != nil {
		// Rollback the transaction on error
		tx.Rollback()
		log.Printf("Error creating invoice: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Link the invoice ID to corresponding AirBooking records
	for _, item := range requestPayload.TravelItems {
		var airbooking tableTypes.AirBooking
		if err := tx.First(&airbooking, item.ID).Error; err != nil {
			// Rollback the transaction on error
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Airbooking not found"})
			return
		}

		airbooking.IDInvoice = &newInvoice.ID
		airbooking.Status = item.Status

		if err := tx.Save(&airbooking).Error; err != nil {
			// Rollback the transaction on error
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Commit the transaction
	tx.Commit()

	// Respond with the created invoice
	c.JSON(http.StatusCreated, newInvoice)
}

// r.PUT("/invoices/:id", func(c *gin.Context) {
// 	id := c.Query("id")

// 	var updatedInvoice tableTypes.Invoice
// 	if err := c.ShouldBindJSON(&updatedInvoice); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	var existingInvoice tableTypes.Invoice
// 	if err := initializers.DB.First(&existingInvoice, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
// 		return
// 	}

// 	// Update fields
// 	existingInvoice.DueDate = updatedInvoice.DueDate
// 	existingInvoice.Tag = updatedInvoice.Tag

// 	// Update related AirBookings
// 	var airBookingIDs []int
// 	for _, travelItem := range updatedInvoice.TravelItems {
// 		airBookingIDs = append(airBookingIDs, travelItem.ID)
// 	}

// 	if err := initializers.DB.Model(&tableTypes.AirBooking{}).
// 		Where("id IN (?)", airBookingIDs).
// 		Updates(map[string]interface{}{"Status": "invoiced", "Id_invoice": existingInvoice.ID}).
// 		Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if err := initializers.DB.Save(&existingInvoice).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, existingInvoice)
// })

func DeleteInvoices(c *gin.Context) {
	id := c.Query("id")

	// Convert string id to uint
	invoiceID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var invoice tableTypes.Invoice
	if err := initializers.DB.First(&invoice, invoiceID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}

	// Check if payment has been received for the invoice
	var paymentCount int64
	if err := initializers.DB.Model(&tableTypes.InvoicePaymentReceived{}).
		Where("id_invoice = ?", invoice.ID).
		Count(&paymentCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if credit_apply is 0 and balance equals amount
	CreditApply, err := strconv.Atoi(invoice.CreditApply)
	if CreditApply != 0 || invoice.Balance != invoice.Amount || paymentCount > 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invoice cannot be deleted due to existing payments or non-zero credit_apply/balance"})
		return
	}

	// Update associated AirBooking records to set IDInvoice to null
	if err := initializers.DB.Model(&tableTypes.AirBooking{}).
		Where("id_invoice = ?", invoice.ID).
		Updates(map[string]interface{}{"id_invoice": nil, "status": "void"}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete the invoice
	if err := initializers.DB.Delete(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invoice deleted successfully"})
}
