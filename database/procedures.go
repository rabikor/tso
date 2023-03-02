package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type (
	proceduresRepository interface {
		GetAll(limit, offset int) ([]Procedure, error)
		GetById(id uint) (*Procedure, error)
		Add(procedure *Procedure) error
	}
	proceduresTable struct {
		*sqlx.DB
	}
)

type Procedure struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

func (db *proceduresTable) GetAll(limit, offset int) (ps []Procedure, _ error) {
	err := db.Select(&ps, fmt.Sprintf("SELECT * FROM procedures LIMIT %d OFFSET %d", limit, offset))
	if err != nil {
		return nil, err
	}

	return ps, nil
}

func (db *proceduresTable) GetById(id uint) (procedure *Procedure, err error) {
	procedure = new(Procedure)
	if err = db.Get(procedure, "SELECT * FROM procedures WHERE id = ?", id); err != nil {
		return
	}

	return
}

func (db *proceduresTable) Add(procedure *Procedure) error {
	if result, err := db.Exec("INSERT INTO procedures (title) VALUES (?)", procedure.Title); err != nil {
		return err
	} else {
		lastId, _ := result.LastInsertId()
		procedure.ID = uint(lastId)
	}

	return nil
}
