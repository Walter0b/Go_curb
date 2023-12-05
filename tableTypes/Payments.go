package tableTypes

import (
	"Go_curb/Database/components"
	"time"
)

type PaymentReceived struct {
	ID                    int       `gorm:"column:id;primaryKey"`
	Number                string    `gorm:"column:number;not null"`
	Date                  time.Time `gorm:"column:date;not null"`
	Balance               string    `gorm:"column:balance;default:0"`
	Amount                string    `column:"amount" gorm:"type:numeric"`
	CurrencyRate          float64   `gorm:"column:currency_rate;not null"`
	Fop                   string    `gorm:"column:fop;not null"`
	Reference             string    `gorm:"column:reference"`
	DeductedTax           bool      `gorm:"column:deducted_tax;not null"`
	Note                  string    `gorm:"column:note"`
	UsedAmount            string    `gorm:"column:used_amount;not null"`
	Status                string    `gorm:"column:status;not null"`
	BaseAmount            string    `gorm:"column:base_amount;not null"`
	IsReconciled          bool      `gorm:"column:is_reconciled;not null"`
	Slug                  int64     `gorm:"column:slug;not null"`
	CustomerID            int       `gorm:"column:id_customer"`
	IDChartOfAccountsFrom int64     `gorm:"column:id_chart_of_accounts_from;default:94"`
	Type                  string    `gorm:"column:type;default:customer_payment"`
	IDConsultant          int       `gorm:"column:id_consultant;default:6"`
	IDChartOfAccounts     int       `gorm:"column:id_chart_of_accounts;default:33"`
	IDCurrency            int       `gorm:"column:id_currency;default:550"`
	HiddenField           string    `gorm:"column:hidden_field"`
	TransfertType         string    `gorm:"column:transfert_type;default:sales_without_invoices"`
	AlreadyUsed           int       `gorm:"column:already_used;not null"`
	ReceipiantName        string    `gorm:"column:receipiant_name"`
	Tag                   string    `gorm:"column:tag;default:2"`
}

type PaymentCustomer struct {
	PaymentReceived
	Customer Customer `gorm:"primaryKey:CustomerID"`
}

func (PaymentReceived) TableName() string {
	return "payment_received"
}

func (p *PaymentReceived) UsedAmountFloat() float64 {
	usedAmountFloat, _ := components.ConvertStringToFloat64(p.UsedAmount)
	return usedAmountFloat
}
func (i *PaymentReceived) BalanceFloat() float64 {
	balanceFloat, _ := components.ConvertStringToFloat64(i.Balance)
	return balanceFloat
}
func (i *PaymentReceived) AmountFloat() float64 {
	amountFloat, _ := components.ConvertStringToFloat64(i.Amount)
	return amountFloat
}
