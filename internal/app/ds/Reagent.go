package ds

type Reagent struct {
	ID          uint    `gorm:"primaryKey"`
	Name        string  `gorm:"type:varchar(25);not null"`
	Img         string  `gorm:"type:varchar(100)"`
	Video       string  `gorm:"type:varchar(100)"`
	Description string  `gorm:"type:varchar(200)"`
	IsDeleted   bool    `gorm:"column:is_deleted;default:false"`
	Formula     string  `gorm:"type:varchar(15);not null"`
	MolarMass   float64 `gorm:"type:decimal(10,2)"`
}
