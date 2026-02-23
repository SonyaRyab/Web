package ds

type Reagent struct {
	ID          int     `gorm:"primaryKey"`
	IsDelete    bool    `gorm:"type:boolean not null;default:false"`
	Img         string  `gorm:"type:varchar(100)"`
	Name        string  `gorm:"type:varchar(25);not null"`
	Formula     string  `gorm:"type:varchar(15);not null"`
	Description string  `gorm:"type:varchar(200)"`
	MolarMass   float64 `gorm:"type:decimal(10,2)"`
	Coefficient float64 `gorm:"type:decimal(10,4)"`
	Purity      float64 `gorm:"type:decimal(5,2)"`
	Price       float64 `gorm:"type:decimal(10,2)"`
}
