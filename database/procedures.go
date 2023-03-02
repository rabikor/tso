package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type (
	proceduresRepository interface {
		All(limit, offset int) ([]Procedure, error)
		ByID(id uint) (Procedure, error)
		Add(p Procedure) (uint, error)
	}
	proceduresTable struct {
		*sqlx.DB
	}
)

type Procedure struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

func (db proceduresTable) All(limit, offset int) (ps []Procedure, _ error) {
	return ps, db.Select(&ps, fmt.Sprintf("SELECT * FROM procedures LIMIT %d OFFSET %d", limit, offset))
}

func (db proceduresTable) ByID(id uint) (p Procedure, err error) {
	return p, db.Get(&p, "SELECT * FROM procedures WHERE id = ?", id)
}

func (db proceduresTable) Add(p Procedure) (uint, error) {
	r, err := db.Exec("INSERT INTO procedures (title) VALUES (?)", p.Title)
	if err != nil {
		return 0, err
	}

	lastID, _ := r.LastInsertId()
	return uint(lastID), nil
}
