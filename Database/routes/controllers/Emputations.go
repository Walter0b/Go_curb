package controllers

import (
	"Go_curb/Database/components"
	"Go_curb/Database/initializers"
	"Go_curb/tableTypes"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllInvoicePayments(c *gin.Context) {

	id := c.Query("id")
	var invoices []tableTypes.InvoicePaymentReceived
	query := initializers.DB.Model(&tableTypes.InvoicePaymentReceived{}).Where("tag = '2'")
	embedField := c.Query("embed")
	components.Get(c, query, &invoices, embedField, id)
}

// CreateInvoiceImputations handles the creation of invoice imputations
func CreateInvoiceImputations(c *gin.Context) {
	var imputations []tableTypes.InvoicePaymentReceived

	if err := c.BindJSON(&imputations); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind JSON data", "message": err.Error()})
		return
	}

	// Start a database transaction
	tx := initializers.DB.Begin()

	var updatedImputations []tableTypes.InvoicePaymentReceived

	for _, imputationInput := range imputations {

		// Fetch details of the invoice and payment
		var invoiceDetails tableTypes.Invoice
		if err := tx.First(&invoiceDetails, imputationInput.IDInvoice).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to fetch invoice"})
			return
		}

		var paymentDetails tableTypes.PaymentReceived
		if err := tx.First(&paymentDetails, imputationInput.IDPaymentReceived).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to fetch payment received"})
			return
		}

		// Check if the amountApply exceeds the invoice balance
		amountApplyFloat := imputationInput.AmountApplyFloat()
		if amountApplyFloat > invoiceDetails.BalanceFloat() {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Amount apply exceeds the invoice balance"})
			return
		}

		// Calculate the net balance in payment
		netBalance := paymentDetails.AmountFloat() - paymentDetails.UsedAmountFloat()

		// Check if the amountApply exceeds the net balance in payment
		if amountApplyFloat > netBalance {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Amount apply exceeds the net balance in payment"})
			return
		}

		// Fetch existing imputation
		var existingImputation tableTypes.InvoicePaymentReceived
		if err := tx.Where("id_invoice = ? AND id_payment_received = ?", imputationInput.IDInvoice, imputationInput.IDPaymentReceived).First(&existingImputation).Error; err == nil {
			// If an existing imputation is found, update and append to the updatedImputations slice

			oldUsedAmount := paymentDetails.UsedAmountFloat() - existingImputation.AmountApplyFloat()

			oldPaymentBalance := paymentDetails.AmountFloat() + oldUsedAmount

			if amountApplyFloat > oldPaymentBalance {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": "Amount apply exceeds the payment balance"})
				return
			}
			newUsedAmount := oldUsedAmount + amountApplyFloat
			newPaymentBalance := paymentDetails.AmountFloat() - newUsedAmount

			// Update payment with new imputation values
			if err := tx.Model(&paymentDetails).Updates(map[string]interface{}{"used_amount": newUsedAmount, "balance": newPaymentBalance}).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment"})
				return
			}
			// AmountApplyFloat, _ := components.ConvertStringToFloat64(existingImputation.AmountApply)

			// Update invoice with new imputation values
			oldCreditApply := (invoiceDetails.CreditApplyFloat() - existingImputation.AmountApplyFloat())
			oldBalance := invoiceDetails.AmountFloat() - oldCreditApply

			// Check if the amountApply exceeds the net balance in payment
			if amountApplyFloat > oldBalance {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": "Amount apply exceeds the net balance in invoice"})
				return
			}
			newCreditApply := oldCreditApply + amountApplyFloat
			newBalance := oldBalance - amountApplyFloat

			if err := tx.Model(&invoiceDetails).Updates(map[string]interface{}{"balance": newBalance, "credit_apply": newCreditApply}).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update invoice"})
				return
			}
			// Update existing imputation with new amount apply
			if err := tx.Model(&existingImputation).Update("amount_apply", imputationInput.AmountApply).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update existing imputation"})
				return
			}

			// Append the updated imputation to the slice
			updatedImputations = append(updatedImputations, existingImputation)
		} else {
			// If no existing imputation, proceed with new imputation
			// Update payment with new imputation values
			newUsedAmount := paymentDetails.UsedAmountFloat() + amountApplyFloat
			newPaymentBalance := paymentDetails.AmountFloat() - newUsedAmount

			if err := tx.Model(&paymentDetails).Updates(map[string]interface{}{"used_amount": newUsedAmount, "balance": newPaymentBalance}).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment"})
				return
			}

			// Update invoice with new imputation values
			newCreditApply := invoiceDetails.CreditApplyFloat() + amountApplyFloat
			newBalance := invoiceDetails.AmountFloat() - newCreditApply

			if err := tx.Model(&invoiceDetails).Updates(map[string]interface{}{"balance": newBalance, "credit_apply": newCreditApply}).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update invoice"})
				return
			}

			// Insert a new imputation record
			newImputation := tableTypes.InvoicePaymentReceived{
				IDPaymentReceived: imputationInput.IDPaymentReceived,
				IDInvoice:         imputationInput.IDInvoice,
				AmountApply:       imputationInput.AmountApply,
			}

			if err := tx.Create(&newImputation).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invoice payment received"})
				return
			}

			// Append the newly created imputation to the slice
			updatedImputations = append(updatedImputations, newImputation)
		}
	}

	// Commit the transaction
	tx.Commit()

	// Return the updated object data for all the imputations
	c.JSON(http.StatusOK, gin.H{
		"imputations": updatedImputations,
	})
}
