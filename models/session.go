package models

import (
	"time"

	"luuhai48/short/db"
)

// Models ======================================================================

type Session struct {
	BaseModel

	UserID   string
	User     User `gorm:"constraint:OnDelete:CASCADE"`
	Username string

	Valid      bool      `gorm:"default:true"`
	ValidUntil time.Time `gorm:"not null"`
}

// Functions ===================================================================

func CreateSession(s *Session) error {
	return db.DB.Create(&s).Error
}

func FindSessionByID(ID string) (*Session, error) {
	var session Session
	if err := db.DB.Where("id = ?", ID).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func DeleteSessionById(ID string) error {
	return db.DB.Where("id = ?", ID).Delete(&Session{}).Error
}

func DeleteOldSessions() error {
	return db.DB.Delete(&Session{}, "valid_until < ?", time.Now()).Error
}
