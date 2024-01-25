package models

import "luuhai48/short/db"

// Models ======================================================================

type Short struct {
	BaseModel

	UserID string
	User   User `gorm:"constraint:OnDelete:CASCADE"`

	Url string `gorm:"not null"`
}

// Functions ===================================================================

func CreateShort(s *Short) error {
	return db.DB.Create(&s).Error
}

func FindShortByID(ID string) (*Short, error) {
	var short Short
	if err := db.DB.Where("id = ?", ID).First(&short).Error; err != nil {
		return nil, err
	}
	return &short, nil
}

func ListShortOfUser(uid string) ([]Short, error) {
	shorts := []Short{}
	err := db.DB.Model(&Short{}).Where("user_id = ?", uid).Order("created_at desc").
		Find(&shorts).Error
	return shorts, err
}

func DeleteShortByID(ID string) error {
	return db.DB.Where("id = ?", ID).Delete(&Short{}).Error
}
