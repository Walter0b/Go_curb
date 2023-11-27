package controllers

import (
	"Go_curb/Database/components"
	"Go_curb/Database/initializers"
	"Go_curb/tableTypes"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
)

// func GetAllInvoicePayments(c *gin.Context) {
// 	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
// 		return
// 	}

// 	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
// 	if err != nil || pageSize < 1 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
// 		return
// 	}

// 	var Imputations []tableTypes.InvoicePaymentReceived
// 	var totalRowCount int64 // Total count of records

// 	// Count total records
// 	if err := initializers.DB.Model(&tableTypes.InvoicePaymentReceived{}).Count(&totalRowCount).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	offset := (page - 1) * pageSize

// 	if err := initializers.DB.Limit(pageSize).Offset(offset).Find(&Imputations).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	response := gin.H{
// 		"data":          Imputations,   // Data for the current page
// 		"totalRowCount": totalRowCount, // Total count of records
// 	}

// 	c.JSON(http.StatusOK, response)
// }

func GetAllInvoicePayments(c *gin.Context) {
	var invoices []tableTypes.Invoice
	var invoicesCustomerType []tableTypes.InvoicePaymentReceived
	query := initializers.DB.Model(&tableTypes.Invoice{}).Where("tag = '2'")
	embedType := reflect.TypeOf(tableTypes.InvoicePaymentReceived{})
	embedField := c.Query("embed")
	components.PaginateWithEmbed(c, query, &invoices, &invoicesCustomerType, embedType, embedField)
}

// CreateInvoiceImputations handles the creation of invoice imputations
func CreateInvoiceImputations(c *gin.Context) {

	var InvoicePaymentReceived []tableTypes.InvoicePaymentReceived

	// Start a database transaction
	tx := initializers.DB.Begin()

	if err := c.BindJSON(&InvoicePaymentReceived); err != nil {
		// fmt.Printf("Parsed JSON data: %+v\n", InvoicePaymentReceived)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind JSON data", "message": err.Error()})
		return
	}

	for i := range InvoicePaymentReceived {

		tx := initializers.DB.Begin()
		var invoice tableTypes.Invoice
		if err := tx.Where("id = ?", InvoicePaymentReceived[i].IDInvoice).First(&invoice).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to fetch invoice"})
			return
		}

		// Fetch the payment received by primary key
		var PaymentReceived tableTypes.PaymentReceived
		if err := tx.Where("id = ?", InvoicePaymentReceived[i].IDPaymentReceived).First(&PaymentReceived).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to fetch payment received"})
			return
		}

		amountApplyFloat, err := components.ConvertStringToFloat64(InvoicePaymentReceived[i].AmountApply)

		if amountApplyFloat < 0 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid negative amount apply"})
			return
		}

		balanceInvoiceFloat, err := components.ConvertStringToFloat64(invoice.Balance)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse balance of invoice"})
			return
		}

		creditApplyFloat, err := components.ConvertStringToFloat64(invoice.CreditApply)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse credit apply of invoice"})
			return
		}

		balancePaymentReceivedFloat, err := components.ConvertStringToFloat64(PaymentReceived.Balance)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse balance of payment received"})
			return
		}

		usedAmountPaymentReceivedFloat, err := components.ConvertStringToFloat64(PaymentReceived.UsedAmount)
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
		PaymentReceived.Balance = balancePaymentReceivedStr
		PaymentReceived.UsedAmount = usedAmountPaymentReceivedStr

		if err := tx.Save(&PaymentReceived).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save payment received"})
			return
		}

		InvoicePaymentReceived[i].InvoiceAmount = invoice.Amount

		if err := tx.Create(&InvoicePaymentReceived[i]).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invoice payment received"})
			return
		}

		if err := tx.Save(&invoice).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save invoice"})
			return
		}
		tx.Commit()

		// Start a new transaction for the next iteration
		tx = initializers.DB.Begin()
	}

	// Commit the transaction if all operations were successful
	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{"message": "Invoice imputations created successfully"})
}
