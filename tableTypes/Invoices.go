package tableTypes

import (
	"Go_curb/Database/components"
	"time"
)

type Invoice struct {
	ID               uint      `gorm:"column:id;primaryKey"`
	CreationDate     time.Time `gorm:"column:creation_date;not null"`
	InvoiceNumber    string    `gorm:"column:invoice_number;not null"`
	Status           string    `gorm:"column:status;not null"`
	DueDate          time.Time `gorm:"column:due_date"`
	Amount           string    `gorm:"amount" gorm:"column:amount"`
	Balance          string    `gorm:"column:balance"`
	NetAmount        string    `gorm:"column:net_amount;not null"`
	TaxAmount        string    `gorm:"column:tax_amount"`
	BaseAmount       string    `gorm:"column:base_amount;not null"`
	PurchaseOrder    string    `gorm:"column:purchase_order"`
	CustomerNotes    string    `gorm:"column:customer_notes"`
	Terms            int       `gorm:"column:terms"`
	TermsConditions  string    `gorm:"column:terms_conditions"`
	CreditApply      string    `gorm:"column:credit_apply;default:0"`
	CreditUsed       string    `gorm:"column:credit_used;default:0"`
	Email            string    `gorm:"column:email"`
	PrintedName      string    `gorm:"column:printed_name"`
	HiddenField      string    `gorm:"column:hidden_field"`
	HiddenIdentifier string    `gorm:"column:hidden_identifier"`
	AlreadyUsed      int       `gorm:"column:already_used;not null"`
	IsOpeningBalance bool      `gorm:"column:is_opening_balance"`
	Tag              string    `gorm:"column:tag;default:2"`
	CustomerID       int       `gorm:"column:id_customer;"`
}

type InvoiceCustomer struct {
	Invoice
	Customer Customer `gorm:"primaryKey:CustomerID"`
}

type RequestPayload struct {
	CustomerID  int    `column:"CustomerID" binding:"required"`
	DueDate     string `column:"dueDate" binding:"required"` // Format: "2006-01-02"
	Amount      string `column:"amount" binding:"required"`
	Tag         string `gorm:"column:tag;default:2"` // Include the tag field in the request payload
	TravelItems []struct {
		ID         int    `column:"id" binding:"required"`
		TotalPrice string `column:"totalPrice" binding:"required"`
		Status     string `gorm:"column:status;not null"`
	} `column:"travelItems" binding:"required"`
}

func (Invoice) TableName() string {
	return "invoice"
}

func (i *Invoice) AmountFloat() float64 {
	amountFloat, _ := components.ConvertStringToFloat64(i.Amount)
	return amountFloat
}
func (i *Invoice) BalanceFloat() float64 {
	balanceFloat, _ := components.ConvertStringToFloat64(i.Balance)
	return balanceFloat
}
func (i *Invoice) CreditApplyFloat() float64 {
	creditApplytFloat, _ := components.ConvertStringToFloat64(i.CreditApply)
	return creditApplytFloat
}
