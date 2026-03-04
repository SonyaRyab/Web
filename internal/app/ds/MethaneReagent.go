package ds

type MethaneReagent struct {
	ID         uint   `gorm:"primaryKey"`
	Methane_id uint   `gorm:"not null;uniqueIndex:idx_exp_reagent"`
	Reagent_id uint   `gorm:"not null;uniqueIndex:idx_exp_reagent"`
	Formula    string `gorm:"type:varchar(15);not null"`
	// Дополнительные поля м-м
	Quantity float64 `gorm:"type:decimal(10,2);not null"`
	Methane  Methane `gorm:"foreignKey:MethaneID"`
	Reagent  Reagent `gorm:"foreignKey:ReagentID"`
}
