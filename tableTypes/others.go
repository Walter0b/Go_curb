package tableTypes

type Entity struct {
	ID   int64 `gorm:"primaryKey"`
	Name string
}

type Currency Entity

func (Currency) TableName() string {
	return "currency"
}
