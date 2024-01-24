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

func CheckUsernameExists(username string) (bool, error) {
	var exists bool
	if err := db.DB.Model(&User{}).
		Select("count(*) > 0").Where("username = ?", username).
		Find(&exists).Error; err != nil {
		return false, err
	}
	return exists, nil
}

func FindUserByUsername(username string) (*User, error) {
	var user User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
