package ds

type Reagent struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Name        string  `gorm:"type:varchar(25);not null" json:"name"`
	Img         string  `gorm:"type:varchar(100)" json:"img"`
	Video       string  `gorm:"type:varchar(100)" json:"video"`
	Description string  `gorm:"type:varchar(200)" json:"description"`
	IsDeleted   bool    `gorm:"column:is_deleted;default:false" json:"-"`
	Formula     string  `gorm:"type:varchar(15);not null" json:"formula"`
	MolarMass   float64 `gorm:"type:decimal(10,2)" json:"molar_mass"`
}
