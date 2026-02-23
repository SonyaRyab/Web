package ds

type ExperimentReagent struct {
	ID           uint `gorm:"primaryKey"`
	ExperimentID uint `gorm:"not null;uniqueIndex:idx_exp_reagent"`
	ReagentID    uint `gorm:"not null;uniqueIndex:idx_exp_reagent"`
	// Дополнительные поля м-м
	Quantity     float64    `gorm:"type:decimal(10,2);not null"`
	ActualAmount float64    `gorm:"type:decimal(10,2)"`
	OrderNum     int        `gorm:"default:1"`
	Experiment   Experiment `gorm:"foreignKey:ExperimentID"`
	Reagent      Reagent    `gorm:"foreignKey:ReagentID"`
}
