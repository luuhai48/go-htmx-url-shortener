package models

import "luuhai48/short/db"

// Types =======================================================================
type AccountState string

const (
	AccountStateNormal  AccountState = "normal"
	AccountStateBlocked AccountState = "blocked"
)

// Models ======================================================================
type User struct {
	BaseModel

	Username string `gorm:"not null;unique;index;size:100"`
	Password string `gorm:"not null"`

	AccountStatus AccountState `gorm:"not null;default:'normal'"`
}

// Functions ===================================================================
func CreateUser(u *User) error {
	return db.DB.Create(&u).Error
}
