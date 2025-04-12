package model

type PersonalToken struct {
	ID        string `gorm:"not null;uniqueIndex;primary_key" db:"id, primarykey" json:"id"`
	Token     string `gorm:"not null" db:"token" json:"token"`
	UserID    string `gorm:"not null" db:"user_id" json:"user_id"`
	User      User   `gorm:"foreignKey:UserID" json:"user"`
	CreatedAt int64  `gorm:"not null" db:"created_at" json:"created_at"`
	table     string `gorm:"-"`
}

func (p PersonalToken) TableName() string {
	if p.table != "" {
		return p.table
	}
	return "personal_tokens"
}
