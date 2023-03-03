package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type (
	DrugsRepository interface {
		All(limit, offset int) (ds []Drug, _ error)
		ByID(id uint) (Drug, error)
		Add(title string) (uint, error)
	}
	DrugsTable struct {
		*sqlx.DB
	}
)

type Drug struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

func NewDrugsRepository(db *sqlx.DB) DrugsRepository {
	return DrugsTable{db}
}

func (db DrugsTable) All(limit, offset int) (ds []Drug, _ error) {
	return ds, db.Select(&ds, fmt.Sprintf("SELECT * FROM drugs LIMIT %d OFFSET %d", limit, offset))
}

func (db DrugsTable) ByID(id uint) (d Drug, err error) {
	return d, db.Get(&d, "SELECT * FROM drugs WHERE id = ?", id)
}

func (db DrugsTable) Add(t string) (uint, error) {
	result, err := db.Exec("INSERT INTO drugs (title) VALUES (?)", t)
	if err != nil {
		return 0, err
	}

	id, _ := result.LastInsertId()
	return uint(id), nil
}
