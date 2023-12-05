package tableTypes

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

func (AirBooking) TableName() string {
	return "air_booking"
}
