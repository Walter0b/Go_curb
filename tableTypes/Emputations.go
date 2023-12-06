package tableTypes

import (
	"Go_curb/Database/components"
	"database/sql"
)

type InvoicePaymentReceived struct {
	IDInvoice         int            `column:"id_invoice;foreignKey:ID"`
	IDPaymentReceived int            `gorm:"column:id_payment_received;foreignKey:ID"`
	GainLossAmount    sql.NullString `gorm:"column:gain_loss_amount"`
	AmountApply       string         `column:"amount_apply"`
	GainLoss          string         `gorm:"column:gain_loss;default:gain"`
	WithholdingTax    sql.NullString `column:"withholding_tax"`
	PaymentAmount     string         `column:"payment_amount"`
	InvoiceAmount     string         `column:"invoice_amount"`
	ID                int            `column:"id"`
	Slug              int64          `column:"slug"`
	HiddenField       string         `column:"hidden_field"`
	AlreadyUsed       int64          `gorm:"column:already_used;default:0"`
	Tag               string         `gorm:"column:tag;default:2"`
}

func (InvoicePaymentReceived) TableName() string {
	return "invoice_payment_received"
}
func (i *InvoicePaymentReceived) AmountApplyFloat() float64 {
	AmountApplyFloat, _ := components.ConvertStringToFloat64(i.AmountApply)
	return AmountApplyFloat
}
