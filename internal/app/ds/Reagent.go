package ds

type Reagent struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Name        string  `gorm:"type:varchar(25);not null" json:"name"`
	Formula     string  `gorm:"type:varchar(100)" json:"formula"`
	Temperature float64 `gorm:"type:decimal(6,2)" json:"temperature"`
	Img         string  `gorm:"type:varchar(100)" json:"img"`
	Video       string  `gorm:"type:varchar(100)" json:"video"`
	Description string  `gorm:"type:varchar(200)" json:"description"`
	MolarMass   float64 `gorm:"column:molarmass;type:decimal(10,2)" json:"molarmass"`
	IsDeleted   bool    `gorm:"column:is_deleted;default:false" json:"-"`
}

type ReagentsPage struct {
	Items      []Reagent `json:"items"`
	Total      int64     `json:"total"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	TotalPages int       `json:"totalPages"`
}
