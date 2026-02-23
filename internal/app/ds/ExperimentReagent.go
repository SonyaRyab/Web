package ds

//type MessageChat struct {
//	ID uint `gorm:"primaryKey"`
//	// здесь создаем Unique key, указывая общий uniqueIndex
//	MessageID uint `gorm:"not null;uniqueIndex:idx_message_chat"`
//	ChatID    uint `gorm:"not null;uniqueIndex:idx_message_chat"`
//
//	Sound bool `gorm:"default:true"`
//
//	Message Message `gorm:"foreignKey:MessageID"`
//	Chat    Chat    `gorm:"foreignKey:ChatID"`
//}

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
