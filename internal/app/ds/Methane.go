// заявки
package ds

import (
	"time"
)

type Methane struct {
	ID 			 uint 			  `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Name         string           `gorm:"type:varchar(100);not null;column:name" json:"name"`
	Status       string           `gorm:"type:varchar(20);column:status" json:"status"`
	DateCreate   time.Time        `gorm:"column:date_create" json:"date_create"`
	DateForm     *time.Time       `gorm:"default:null;column:date_form" json:"date_form"`
	DateFinish   *time.Time       `gorm:"default:null;column:date_finish" json:"date_finish"`

	AdminID      uint             `gorm:"not null;column:admin_id" json:"-"`
	ModeratorID  *uint            `gorm:"column:moderator_id" json:"-"`

	Temperature  *float64         `gorm:"column:temperature" json:"temperature,omitempty"`
	MethaneYield *float64         `gorm:"column:methane_yield" json:"methane_yield,omitempty"`

	Admin        User             `gorm:"foreignKey:AdminID;references:ID" json:"researcher"`
	Moderator    *User            `gorm:"foreignKey:ModeratorID;references:ID" json:"professor,omitempty"`
	Reagents     []MethaneReagent `gorm:"foreignKey:MethaneID" json:"reagents,omitempty"`
}

type MethaneListSerializer struct {
	ID           uint       `json:"id"`
	Name         string     `json:"name"`
	Status       string     `json:"status"`
	DateCreate   time.Time  `json:"date_create"`
	DateForm     *time.Time `json:"date_form"`
	ReagentCount int64      `json:"reagent_count"`
}

type FullMethaneSerializer struct {
	Methane
	ResearcherName string `json:"researcher_name"`
	ProfessorName  string `json:"professor_name,omitempty"`
}