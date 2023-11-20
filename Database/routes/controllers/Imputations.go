package controllers

import (
	"Go_curb/Database/components"
	"Go_curb/Database/initializers"
	"Go_curb/tableTypes"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllInvoicePayments(c *gin.Context) {
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

	var Imputations []tableTypes.InvoicePaymentReceived
	var totalRowCount int64 // Total count of records

	// Count total records
	if err := initializers.DB.Model(&tableTypes.InvoicePaymentReceived{}).Count(&totalRowCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	offset := (page - 1) * pageSize

	if err := initializers.DB.Limit(pageSize).Offset(offset).Find(&Imputations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"data":          Imputations,   // Data for the current page
		"totalRowCount": totalRowCount, // Total count of records
	}

	c.JSON(http.StatusOK, response)
}

// CreateInvoiceImputations handles the creation of invoice imputations
func CreateInvoiceImputations(c *gin.Context) {
	var invoice_payment_received []tableTypes.InvoicePaymentReceived
	var invoice tableTypes.Invoice
	var payment_received tableTypes.PaymentReceived
	tabReplace := map[string]string{
		"$": "",
		",": "",
	}

	// Start a database transaction
	tx := initializers.DB.Begin()

	if err := c.BindJSON(&invoice_payment_received); err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind JSON data"})
		return
	}

	for i := range invoice_payment_received {
		if err := tx.First(&invoice, invoice_payment_received[i].IDInvoice).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to fetch invoice"})
			return
		}

		if err := tx.First(&payment_received, invoice_payment_received[i].IDPaymentReceived).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to fetch payment received"})
			return
		}

		amountApplyFloat := invoice_payment_received[i].AmountApply

		if amountApplyFloat < 0 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid negative amount apply"})
			return
		}

		balanceInvoiceFloat, err := strconv.ParseFloat(components.ReplaceAllMultiple(invoice.Balance, tabReplace), 64)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse balance of invoice"})
			return
		}

		creditApplyFloat, err := strconv.ParseFloat(components.ReplaceAllMultiple(invoice.CreditApply, tabReplace), 64)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse credit apply of invoice"})
			return
		}

		balancePaymentReceivedFloat, err := strconv.ParseFloat(components.ReplaceAllMultiple(payment_received.Balance, tabReplace), 64)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse balance of payment received"})
			return
		}

		usedAmountPaymentReceivedFloat, err := strconv.ParseFloat(components.ReplaceAllMultiple(payment_received.UsedAmount, tabReplace), 64)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse used amount of payment received"})
			return
		}

		usedAmountPaymentReceivedFloat += amountApplyFloat
		balanceInvoiceFloat -= amountApplyFloat

		if balanceInvoiceFloat < 0 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Amount apply exceeds the balance of the invoice"})
			return
		}

		creditApplyFloat += amountApplyFloat
		balancePaymentReceivedFloat -= amountApplyFloat

		balanceInvoiceStr := strconv.FormatFloat(balanceInvoiceFloat, 'f', 2, 64)
		creditApplyStr := strconv.FormatFloat(creditApplyFloat, 'f', 2, 64)
		balancePaymentReceivedStr := strconv.FormatFloat(balancePaymentReceivedFloat, 'f', 2, 64)
		usedAmountPaymentReceivedStr := strconv.FormatFloat(usedAmountPaymentReceivedFloat, 'f', 2, 64)

		if invoice.Balance == "0" {
			invoice.Status = "paid"
		}

		invoice.Balance = balanceInvoiceStr
		invoice.CreditApply = creditApplyStr
		payment_received.Balance = balancePaymentReceivedStr
		payment_received.UsedAmount = usedAmountPaymentReceivedStr

		if err := tx.Save(&payment_received).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save payment received"})
			return
		}

		invoice_payment_received[i].InvoiceAmount = components.ConvertStringToFloat64(invoice.Amount)

		if err := tx.Create(&invoice_payment_received[i]).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invoice payment received"})
			return
		}

		if err := tx.Save(&invoice).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save invoice"})
			return
		}
	}

	// Commit the transaction if all operations were successful
	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{"message": "Invoice imputations created successfully"})
}
