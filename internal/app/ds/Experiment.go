// заявки
package ds

import (
	"database/sql"
	"time"
)

type Experiment struct {
	ID          uint      `gorm:"primaryKey"`
	Status      string    `gorm:"type:varchar(15);not null"`
	DateCreate  time.Time `gorm:"not null"`
	DateUpdate  time.Time
	DateFinish  sql.NullTime `gorm:"default:null"`
	CreatorID   uint         `gorm:"not null"`
	ModeratorID uint         `gorm:"default:null"`
	Temperature float64      `gorm:"type:decimal(6,2)"`
	//Pressure     float64      `gorm:"type:decimal(8,2)"`
	MethaneYield float64 `gorm:"type:decimal(5,2)"`
	Notes        string  `gorm:"type:text"`
	Creator      Users   `gorm:"foreignKey:CreatorID"`
	Moderator    Users   `gorm:"foreignKey:ModeratorID"`
}
