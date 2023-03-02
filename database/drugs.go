package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type (
	drugsRepository interface {
		All(limit, offset int) (ds []Drug, _ error)
		ByID(id uint) (Drug, error)
		Add(d Drug) (uint, error)
	}
	drugsTable struct {
		*sqlx.DB
	}
)

type Drug struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

func (db drugsTable) All(limit, offset int) (ds []Drug, _ error) {
	return ds, db.Select(&ds, fmt.Sprintf("SELECT * FROM drugs LIMIT %d OFFSET %d", limit, offset))
}

func (db drugsTable) ByID(id uint) (d Drug, err error) {
	return d, db.Get(&d, "SELECT * FROM drugs WHERE id = ?", id)
}

func (db drugsTable) Add(d Drug) (uint, error) {
	result, err := db.Exec("INSERT INTO drugs (title) VALUES (?)", d.Title)
	if err != nil {
		return 0, err
	}

	id, _ := result.LastInsertId()
	return uint(id), nil
}
