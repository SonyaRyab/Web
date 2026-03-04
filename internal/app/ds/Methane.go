// заявки
package ds

import (
	"time"
)

type Methane struct {
	ID           uint   `gorm:"primaryKey" gorm:"table:methane"`
	Name         string `gorm:"type:varchar(100);not null"`
	Status       string `gorm:"type:varchar(20)"`
	DateCreate   time.Time
	DateUpdate   time.Time `gorm:"default:null"`
	DateFinish   time.Time `gorm:"default:null"`
	AdminID      uint      `gorm:"not null"`
	ModeratorID  *uint
	Temperature  float64 `gorm:"type:decimal(6,2)"`
	MethaneYield float64 `gorm:"type:decimal(5,2)"`
	Admin        Users   `gorm:"foreignKey:AdminID"`
	Moderator    Users   `gorm:"foreignKey:ModeratorID"`
}
