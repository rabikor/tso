package rdb

import "treatment-scheme-organizer/pkg/models"

func Migrate() error {
	return DB.AutoMigrate(
		&models.Drug{},
	)
}
