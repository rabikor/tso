package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"treatment-scheme-organizer/config"
)

type DB struct {
	*gorm.DB
	Drugs drugsRepository
}

func Open() (*DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4",
		config.Env.DB.User,
		config.Env.DB.Password,
		config.Env.DB.Host,
		config.Env.DB.Port,
		config.Env.DB.Name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	dbb, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := dbb.Ping(); err != nil {
		return nil, err
	}

	return &DB{
		DB:    db,
		Drugs: &drugsTable{DB: db},
	}, nil
}
