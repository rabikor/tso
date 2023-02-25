package database

import (
	"github.com/jmoiron/sqlx"
)

type (
	illnessesRepository interface {
		GetAll(limit, offset int) ([]*Illness, error)
		GetById(id uint) (*Illness, error)
		Add(illness *Illness) error
	}
	illnessesTable struct {
		*sqlx.DB
	}
)

type Illness struct {
	ID    uint   `db:"id" json:"id"`
	Title string `db:"title" json:"title"`
}

func (db *illnessesTable) GetAll(limit, offset int) (illnesses []*Illness, _ error) {
	err := db.Select(&illnesses, "SELECT * FROM illnesses LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}

	return illnesses, nil
}

func (db *illnessesTable) GetById(id uint) (illness *Illness, err error) {
	illness = new(Illness)
	if err = db.Get(illness, "SELECT * FROM illnesses WHERE id = ?", id); err != nil {
		return
	}

	return
}

func (db *illnessesTable) Add(illness *Illness) error {
	if result, err := db.Exec("INSERT INTO illnesses (title) VALUES (?)", illness.Title); err != nil {
		return err
	} else {
		lastId, _ := result.LastInsertId()
		illness.ID = uint(lastId)
	}

	return nil
}
