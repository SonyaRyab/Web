package ds

type Users struct {
	ID       uint   `gorm:"primary_key" json:"id"`
	Username string `gorm:"type:varchar(50);unique;not null" json:"username"`
	Email    string `gorm:"type:varchar(25);unique;not null" json:"email"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
	IsAdmin  bool   `gorm:"type:boolean;default:false" json:"is_admin"`
}
