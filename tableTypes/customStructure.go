package tableTypes

import (
	"time"
)

type PaymentReceived struct {
	ID                    int       `gorm:"column:id;primaryKey"`
	Number                string    `gorm:"column:number;not null"`
	Date                  time.Time `gorm:"column:date;not null"`
	Balance               string    `gorm:"column:balance;not null"`
	Amount                string    `json:"amount" gorm:"type:numeric"`
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
	IDChartOfAccountsFrom int64     `gorm:"column:id_chart_of_accounts_from"`
	Type                  string    `gorm:"column:type;not null"`
	IDConsultant          int       `gorm:"column:id_consultant"`
	IDChartOfAccounts     int       `gorm:"column:id_chart_of_accounts;not null"`
	IDCurrency            int       `gorm:"column:id_currency;not null"`
	HiddenField           string    `gorm:"column:hidden_field"`
	TransfertType         string    `gorm:"column:transfert_type"`
	AlreadyUsed           int       `gorm:"column:already_used;not null"`
	ReceipiantName        string    `gorm:"column:receipiant_name"`
	Tag                   string    `gorm:"column:tag"`
}

type InvoicePaymentReceived struct {
	ID                  int64
	Id_invoice          int64
	Id_payment_received int64
	Gain_loss_amount    string `gorm:"default:null"`
	Amount_apply        string
	Gain_loss           string `gorm:"default:gain"`
	Withholding_tax     string `gorm:"default:null"`
	//Payment_amount      string
	Invoice_amount string `gorm:"default:null"`
	Slug           int64  `gorm:"default:10028"`
	Hidden_field   string `gorm:"default:null"`
	Already_used   int64  `gorm:"default:0"`
	Tag            string `json:"tag"`
}

type Customer struct {
	ID                int `gorm:"primaryKey;"`
	Customer_name     string
	Street            string
	City              string
	State             string
	Zip_code          string
	Notes             string
	Terms             uint
	Account_number    string
	Tax_id            string
	Balance           string
	Is_active         bool
	Is_sub_agency     bool
	Language          string
	Slug              int64
	Id_currency       int64
	Id_country        int64
	Irs_share_key     string
	Currency_rate     float32
	Agency            string
	Avoid_deletion    bool
	Is_editable       bool
	Alias             string
	Already_used      int64
	Ab_key            string
	Tmc_client_number string
	Invoices          []Invoice `gorm:"foreignKey:CustomerID"`
}

type Invoice struct {
	ID               uint       `gorm:"column:id;primaryKey"`
	CreationDate     time.Time  `gorm:"column:creation_date;not null"`
	InvoiceNumber    string     `gorm:"column:invoice_number;not null"`
	Status           string     `gorm:"column:status;not null"`
	DueDate          time.Time  `gorm:"column:due_date"`
	Amount           string     `gorm:"amount" gorm:"column:amount"`
	Balance          string     `gorm:"column:balance"`
	NetAmount        string     `gorm:"column:net_amount;not null"`
	TaxAmount        string     `gorm:"column:tax_amount"`
	BaseAmount       string     `gorm:"column:base_amount;not null"`
	PurchaseOrder    string     `gorm:"column:purchase_order"`
	CustomerNotes    string     `gorm:"column:customer_notes"`
	Terms            int        `gorm:"column:terms"`
	TermsConditions  string     `gorm:"column:terms_conditions"`
	CreditApply      string     `gorm:"column:credit_apply"`
	CreditUsed       string     `gorm:"column:credit_used"`
	Email            string     `gorm:"column:email"`
	PrintedName      string     `gorm:"column:printed_name"`
	HiddenField      string     `gorm:"column:hidden_field"`
	HiddenIdentifier string     `gorm:"column:hidden_identifier"`
	AlreadyUsed      int        `gorm:"column:already_used;not null"`
	IsOpeningBalance bool       `gorm:"column:is_opening_balance"`
	Tag              string     `gorm:"column:tag"`
	CustomerID       int        `gorm:"column:id_customer;"`
	Customers        Customer `gorm:"foreignKey:CustomerID"`
}

type AirBooking struct {
	ID                uint   `gorm:"column:id;primaryKey"`
	TotalPrice        string `gorm:"column:total_price"`
	Itinerary         string `gorm:"column:itinerary"`
	TravelerName      string `gorm:"column:traveler_name;not null"`
	TicketNumber      int64  `gorm:"column:ticket_number;not null"`
	ConjunctionNumber int16  `gorm:"column:conjunction_number"`
	Status            string `gorm:"column:status;not null"`
	ProductType       string `gorm:"column:product_type;not null"`
	TransactionType   string `gorm:"column:transaction_type;not null"`
	IDInvoice         *uint  `gorm:"column:id_invoice"`
}

type Country struct {
	ID           int    `gorm:"column:id;primaryKey"`
	Name         string `gorm:"column:name;not null"`
	NameEN       string `gorm:"column:name_en"`
	Code         string `gorm:"column:code"`
	NameFR       string `gorm:"column:name_fr"`
	NamePO       string `gorm:"column:name_po"`
	CurrencyCode string `gorm:"column:currency_code"`
	Slug         int64  `gorm:"column:slug;not null"`
	AlreadyUsed  int    `gorm:"column:already_used;not null"`
}
type RequestPayload struct {
	CustomerID  int    `json:"CustomerID" binding:"required"`
	DueDate     string `json:"dueDate" binding:"required"` // Format: "2006-01-02"
	Amount      string `json:"amount" binding:"required"`
	Tag         string `json:"tag"` // Include the tag field in the request payload
	TravelItems []struct {
		ID         int    `json:"id" binding:"required"`
		TotalPrice string `json:"totalPrice" binding:"required"`
		Status     string `gorm:"column:status;not null"`
	} `json:"travelItems" binding:"required"`
}

func (InvoicePaymentReceived) TableName() string {
	return "invoice_payment_received"
}

// TableName specifies the table name for the struct
func (Invoice) TableName() string {
	return "invoice"
}

// TableName specifies the table name for the struct
func (Country) TableName() string {
	return "country"
}

// TableName specifies the table name for the struct
func (AirBooking) TableName() string {
	return "air_booking"
}

type Entity struct {
	ID   int64 `gorm:"primaryKey"`
	Name string
}

type Currency Entity

func (Currency) TableName() string {
	return "currency"
}
func (Customer) TableName() string {
	return "customer"
}

func (PaymentReceived) TableName() string {
	return "payment_received"
}
