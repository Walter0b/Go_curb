package controllers

import (
	"Go_curb/Database/components"
	"Go_curb/Database/initializers"
	"Go_curb/tableTypes"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func GetPayments(c *gin.Context) {
	var payments []tableTypes.PaymentReceived
	var paymentsEmbedded []tableTypes.PaymentCustomer
	query := initializers.DB.Model(&tableTypes.PaymentReceived{}).Where("tag = '2'").Order("ID DESC")
	embedType := reflect.TypeOf(tableTypes.PaymentCustomer{})
	embedField := c.Query("embed")
	id := c.Query("id")
	components.Get(c, query, &payments, &paymentsEmbedded, embedType, embedField, id)
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
	paymentReceived.Slug = components.GenerateRandomSlug()

	// Automatically generate a unique number with the format "PMR-{dynamic_number}"

	nextNumber := components.GenerateRandomSlug()
	paymentReceived.Number = fmt.Sprintf("PMR-%d-2", nextNumber)

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

// DELETE /payments/:id
func DeletePayments(c *gin.Context) {
	id := c.Query("id")

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

	// Consider additional checks and logic here based on requirements

	// Delete the payment
	if err := initializers.DB.Delete(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment deleted successfully"})
}
