package ds

type MethaneReagent struct {
	ID        uint `gorm:"primaryKey" json:"-"`
	MethaneID uint `gorm:"not null;uniqueIndex:idx_exp_reagent" json:"-"`
	ReagentID uint `gorm:"not null;uniqueIndex:idx_exp_reagent" json:"-"`
	// Дополнительные поля м-м
	Quantity float64 `gorm:"type:decimal(10,2);not null" json:"quantity"`
	//Methane  Methane `gorm:"foreignKey:MethaneID" json:"-"`
	//Reagent  Reagent `gorm:"foreignKey:ReagentID" json:"reagent"`
}