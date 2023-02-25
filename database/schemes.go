package database

import (
	"github.com/jmoiron/sqlx"
)

type (
	schemesRepository interface {
		GetById(id uint) (scheme *Scheme, err error)
		GetByIllness(illnessID, limit, offset int) ([]Scheme, error)
		Add(scheme *Scheme) error
	}
	schemesTable struct {
		*sqlx.DB
		sdTable *schemeDaysTable
	}
)

type Scheme struct {
	ID        uint `json:"id"`
	IllnessID uint `db:"illness_id" json:"-" bson:"-"`
	Length    uint `json:"length"`

	Illness Illness `db:"illness" json:"-"`
}

func (db *schemesTable) GetByIllness(illnessID, limit, offset int) (schemes []Scheme, _ error) {
	sql := `
		SELECT s.*, 
		       i.id as "illness.id", 
		       i.title as "illness.title"
		FROM schemes s 
		    INNER JOIN illnesses i on s.illness_id = i.id
		    LEFT JOIN scheme_days sd on s.id = sd.scheme_id
		WHERE s.illness_id = ?
		LIMIT ? OFFSET ?`

	err := db.Select(&schemes, sql, illnessID, limit, offset)
	if err != nil {
		return nil, err
	}

	return schemes, nil
}

func (db *schemesTable) GetById(id uint) (scheme *Scheme, err error) {
	scheme = new(Scheme)
	if err = db.Get(scheme, "SELECT * FROM schemes WHERE id = ?", id); err != nil {
		return
	}

	return
}

func (db *schemesTable) Add(scheme *Scheme) error {
	var sql = "INSERT INTO schemes (illness_id, length) VALUES (?, ?)"

	if res, err := db.Exec(sql, scheme.Illness.ID, scheme.Length); err != nil {
		return err
	} else {
		lastId, _ := res.LastInsertId()
		scheme.ID = uint(lastId)
	}

	return nil
}
