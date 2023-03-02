package database

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

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

func Open(env *config.Env) (*DB, error) {
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

	sdTable := &schemeDaysTable{DB: db}

	return &DB{
		DB:         db,
		Drugs:      &drugsTable{DB: db},
		Illnesses:  &illnessesTable{DB: db},
		Procedures: &proceduresTable{DB: db},
		Schemes:    &schemesTable{DB: db, sdTable: sdTable},
		SchemeDays: sdTable,
	}, nil
}

func TestDB() *DB {
	env := &config.Env{}
	if err := env.ParseEnv("./../.env"); err != nil {
		log.Fatal(err)
	}

	db, err := Open(env)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func TruncateTables(db *sqlx.DB) {
	db.MustExec("SET FOREIGN_KEY_CHECKS = 0")
	db.MustExec("TRUNCATE TABLE drugs")
	db.MustExec("TRUNCATE TABLE illnesses")
	db.MustExec("TRUNCATE TABLE procedures")
	db.MustExec("TRUNCATE TABLE scheme_days")
	db.MustExec("TRUNCATE TABLE schemes")
	db.MustExec("SET FOREIGN_KEY_CHECKS = 1")
}
