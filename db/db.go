package db

import (
	"log"
	"log/slog"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func getDBLogLevel() logger.LogLevel {
	switch os.Getenv("DB_LOG_LEVEL") {
	case "info":
		return logger.Info
	case "error":
		return logger.Error
	default:
		return logger.Silent
	}
}

func Init() error {
	os.MkdirAll("storage", os.ModePerm)

	db, err := gorm.Open(sqlite.Open("./storage/db.sqlite?_foreign_keys=on"), &gorm.Config{
		Logger: logger.Default.LogMode(getDBLogLevel()),
	})
	if err != nil {
		return err
	}

	DB = db
	return nil
}

func CloseDB() {
	db, err := DB.DB()
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
	slog.Info("Closed database connection")
}
