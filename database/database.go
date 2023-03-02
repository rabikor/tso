package database

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/pressly/goose/v3"

	"treatment-scheme-organizer/config"
)

type DB struct {
	*sqlx.DB
	Drugs      drugsRepository
	Illnesses  illnessesRepository
	Procedures proceduresRepository
	Schemes    schemesRepository
	SchemeDays schemeDaysRepository
}

func Open(env config.Env) (*DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4",
		env.DB.User,
		env.DB.Password,
		env.DB.Host,
		env.DB.Port,
		env.DB.Name,
	)

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	sdt := schemeDaysTable{DB: db}
	return &DB{
		DB:         db,
		Drugs:      drugsTable{DB: db},
		Illnesses:  illnessesTable{DB: db},
		Procedures: proceduresTable{DB: db},
		Schemes:    schemesTable{DB: db, sdTable: sdt},
		SchemeDays: sdt,
	}, nil
}
