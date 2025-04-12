package model

type User struct {
	ID        string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Username  string `gorm:"not null;uniqueIndex" db:"username" json:"username"`
	Email     string `gorm:"not null;uniqueIndex" db:"email" json:"email"`
	Password  string `gorm:"not null" db:"password" json:"-"`
	CreatedAt int64  `gorm:"not null" db:"created_at" json:"created_at"`
	table     string `gorm:"-"`
}

func (p User) TableName() string {
	if p.table != "" {
		return p.table
	}
	return "users"
}
