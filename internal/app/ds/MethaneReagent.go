package ds

type MethaneReagent struct {
	ID        uint `gorm:"primaryKey" json:"-"`
	MethaneID uint `gorm:"not null;uniqueIndex:idx_exp_reagent" json:"-"`
	Reagent_id uint `gorm:"column:reagent_id;not null;uniqueIndex:idx_exp_reagent" json:"-"`
	Volume float64 `gorm:"type:decimal(10,2);not null" json:"volume"`
	// MethaneYield float64 `gorm:"type:decimal(6,2);not null" json:"methane_yield"`

	Reagent Reagent `gorm:"foreignKey:Reagent_id;->" json:"reagent"`
}
