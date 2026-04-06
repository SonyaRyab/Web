package ds

type User struct {
	ID          uint   `gorm:"primary_key" json:"id"`
	Username    string `gorm:"type:varchar(25)" json:"username"`
	Login       string `gorm:"type:varchar(25);unique;not null" json:"login"`
	Email       string `gorm:"type:varchar(25);unique;" json:"email"`
	Password    string `gorm:"type:varchar(100);not null" json:"-"`
	IsModerator bool   `gorm:"type:boolean;default:false" json:"is_moderator"`
}

type UserMethanes struct {
	User     User
	Methanes []Methane `json:"methanes"`
}
