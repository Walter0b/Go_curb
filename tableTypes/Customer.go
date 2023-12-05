package tableTypes

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
}

type CustomerEmbed struct {
	Customer
	Invoices []Invoice         `gorm:"foreignKey:CustomerID"`
	Payments []PaymentReceived `gorm:"foreignKey:CustomerID"`
}

func (Customer) TableName() string {
	return "customer"
}
