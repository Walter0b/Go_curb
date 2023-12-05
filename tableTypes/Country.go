package tableTypes

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

func (Country) TableName() string {
	return "country"
}
