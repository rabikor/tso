package database

import (
	"github.com/jmoiron/sqlx"
)

type (
	IllnessesRepository interface {
		All(limit, offset int) ([]Illness, error)
		ByID(id uint) (Illness, error)
		Add(i Illness) (uint, error)
	}
	IllnessesTable struct {
		*sqlx.DB
	}
)

type Illness struct {
	ID    uint   `db:"id" json:"id"`
	Title string `db:"title" json:"title"`
}

func NewIllnessesRepository(db *sqlx.DB) IllnessesRepository {
	return IllnessesTable{db}
}

func (db IllnessesTable) All(limit, offset int) (is []Illness, _ error) {
	return is, db.Select(&is, "SELECT * FROM illnesses LIMIT ? OFFSET ?", limit, offset)
}

func (db IllnessesTable) ByID(id uint) (i Illness, err error) {
	return i, db.Get(&i, "SELECT * FROM illnesses WHERE id = ?", id)
}

func (db IllnessesTable) Add(i Illness) (uint, error) {
	r, err := db.Exec("INSERT INTO illnesses (title) VALUES (?)", i.Title)
	if err != nil {
		return 0, err
	}

	lastID, _ := r.LastInsertId()
	return uint(lastID), nil
}
