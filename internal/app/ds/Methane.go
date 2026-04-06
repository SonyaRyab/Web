// заявки
package ds

import (
	"time"
)

type Methane struct {
	ID           uint             `gorm:"primaryKey" json:"id"`
	Name         string           `gorm:"type:varchar(100);not null" json:"name"`
	Status       string           `gorm:"type:varchar(20)"`
	DateCreate   time.Time        `json:"date_create"`
	DateForm     time.Time        `gorm:"default:null" json:"date_update"`
	DateFinish   time.Time        `gorm:"default:null" json:"date_finish"`
	AdminID      uint             `gorm:"not null" json:"-"`
	ModeratorID  *uint            `json:"-"`
	Temperature  float64          `gorm:"type:decimal(6,2)" json:"temperature"`
	MethaneYield float64          `gorm:"type:decimal(5,2)" json:"methane_yield"`
	Admin        User             `gorm:"foreignKey:AdminID" json:"admin"`
	Moderator    User             `gorm:"foreignKey:ModeratorID;references:ID" json:"moderator,omitempty"`
	Reagents     []MethaneReagent `gorm:"foreignKey:MethaneID" json:"reagents,omitempty"`
}

// Сериализатор для списка (без деталей)
type MethaneListSerializer struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	DateCreate   time.Time `json:"date_create"`
	DateForm     time.Time `json:"date_form"`
	ReagentCount int64     `json:"reagent_count"`
}

// Сериализатор для детального просмотра
type FullMethaneSerializer struct {
	Methane
	AdminName     string `json:"admin_name"`
	ModeratorName string `json:"moderator_name,omitempty"`
}
