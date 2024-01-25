package main

import (
	"log/slog"

	"luuhai48/short/db"
	"luuhai48/short/models"

	"gorm.io/gorm"
)

func getModels() []interface{} {
	return append(
		make([]interface{}, 0),

		&models.User{},
		&models.Session{},
		&models.Short{},
	)
}

func migrateDB() error {
	models := getModels()

	if len(models) > 0 {
		slog.Info("Migrating database model(s)...")

		if err := db.DB.AutoMigrate(models...); err != nil {
			return err
		}

		// Remove unused column(s)
		for _, model := range models {
			stmt := &gorm.Statement{DB: db.DB}
			stmt.Parse(model)
			fields := stmt.Schema.Fields
			columns, _ := db.DB.Migrator().ColumnTypes(model)

			for i := range columns {
				found := false
				for j := range fields {
					if columns[i].Name() == fields[j].DBName {
						found = true
						break
					}
				}
				if !found {
					if err := db.DB.Migrator().DropColumn(model, columns[i].Name()); err != nil {
						return err
					}
				}
			}
		}

		slog.Info("Done database migration!")
	}
	return nil
}
