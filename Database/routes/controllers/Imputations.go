package controllers

import (
	"Go_curb/Database/components"
	"Go_curb/Database/initializers"
	"Go_curb/tableTypes"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func ReplaceAllMultiple(chaine string, tabReplace map[string]string) string {
	result := chaine

	for old, new := range tabReplace {
		result = strings.ReplaceAll(result, old, new)
	}

	return result
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

	if err := c.BindJSON(&invoice_payment_received); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind JSON data"})
		return
	}

	for i := range invoice_payment_received {
		if err := initializers.DB.First(&invoice, invoice_payment_received[i].IDInvoice).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to fetch invoice"})
			return
		}

		if err := initializers.DB.First(&payment_received, invoice_payment_received[i].IDPaymentReceived).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to fetch payment received"})
			return
		}

		amountApplyFloat := invoice_payment_received[i].AmountApply

		if amountApplyFloat < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid negative amount apply"})
			return
		}

		balanceInvoiceFloat, err := strconv.ParseFloat(ReplaceAllMultiple(invoice.Balance, tabReplace), 64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse balance of invoice"})
			return
		}

		creditApplyFloat, err := strconv.ParseFloat(ReplaceAllMultiple(invoice.CreditApply, tabReplace), 64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse credit apply of invoice"})
			return
		}

		balancePaymentReceivedFloat, err := strconv.ParseFloat(ReplaceAllMultiple(payment_received.Balance, tabReplace), 64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse balance of payment received"})
			return
		}

		usedAmountPaymentReceivedFloat, err := strconv.ParseFloat(ReplaceAllMultiple(payment_received.UsedAmount, tabReplace), 64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse used amount of payment received"})
			return
		}

		usedAmountPaymentReceivedFloat += amountApplyFloat
		balanceInvoiceFloat -= amountApplyFloat

		if balanceInvoiceFloat < 0 {
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

		if err := initializers.DB.Save(&payment_received).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save payment received"})
			return
		}

		invoice_payment_received[i].InvoiceAmount = components.ConvertStringToFloat64(invoice.Amount)

		if err := initializers.DB.Create(&invoice_payment_received[i]).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invoice payment received"})
			return
		}

		if err := initializers.DB.Save(&invoice).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save invoice"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Invoice imputations created successfully"})
	return
}
