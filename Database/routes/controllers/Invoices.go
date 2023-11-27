package controllers

import (
	"Go_curb/Database/components"
	"Go_curb/Database/initializers"
	"Go_curb/tableTypes"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Retrieve all invoices with pagination
func GetAllInvoices(c *gin.Context) {

	id := c.Query("id")
	var invoices []tableTypes.Invoice
	var invoicesCustomerType []tableTypes.InvoiceCustomer
	query := initializers.DB.Model(&tableTypes.Invoice{}).Where("tag = '2'")
	embedType := reflect.TypeOf(tableTypes.InvoiceCustomer{})
	embedField := c.Query("embed")
	components.PaginateWithEmbed(c, query, &invoices, &invoicesCustomerType, embedType, embedField, id)
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
		CustomerID:    requestPayload.CustomerID,
		DueDate:       dueDate,
		Amount:        requestPayload.Amount,
		Status:        "unpaid",
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
