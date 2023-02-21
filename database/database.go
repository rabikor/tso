package database

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"treatment-scheme-organizer/config"
)

type DB struct {
	*sqlx.DB
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

	db, err := sqlx.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DB{
		DB:    db,
		Drugs: &drugsTable{DB: db},
	}, nil
}
