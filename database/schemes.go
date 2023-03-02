package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
)

type (
	schemesRepository interface {
		ByID(id uint) (s Scheme, err error)
		ByIllness(illnessID, limit, offset int) ([]Scheme, error)
		Add(s Scheme) (uint, error)
	}
	schemesTable struct {
		*sqlx.DB
		sdTable schemeDaysTable
	}
)

type Scheme struct {
	ID        uint `json:"id"`
	IllnessID uint `db:"illness_id" json:"-" bson:"-"`
	Length    uint `json:"length"`

	Illness Illness     `db:"illness" json:"-"`
	Days    []SchemeDay `db:"days" json:"-"`
}

func (db schemesTable) ByIllness(illnessID, limit, offset int) (ss []Scheme, _ error) {
	sql := `
		SELECT s.*, 
		       i.id as "illness.id", 
		       i.title as "illness.title"
		FROM schemes s 
		    INNER JOIN illnesses i on s.illness_id = i.id
		    LEFT JOIN scheme_days sd on s.id = sd.scheme_id
		WHERE s.illness_id = ?
		LIMIT ? OFFSET ?`

	return ss, db.Select(&ss, sql, illnessID, limit, offset)
}

func (db schemesTable) ByID(id uint) (s Scheme, err error) {
	if err = db.Get(&s, "SELECT * FROM schemes WHERE id = ?", id); err != nil {
		return s, err
	}

	if err = db.Select(&s.Days, "SELECT * FROM scheme_days WHERE scheme_id = ?", s.ID); err != nil {
		return s, err
	}

	return s, nil
}

func (db schemesTable) Add(s Scheme) (uint, error) {
	const (
		qScheme = "INSERT INTO schemes (illness_id, length) VALUES (?, ?)"
		qDays   = "INSERT INTO scheme_days (scheme_id, procedure_id, drug_id, `order`, times, frequency) VALUES (?, ?, ?, ?, ?, ?)"
	)

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Error(err)
			}
		}
	}()

	r, err := tx.Exec(qScheme, s.IllnessID, s.Length)
	if err != nil {
		return 0, err
	}

	lastID, _ := r.LastInsertId()
	for _, sd := range s.Days {
		_, err = tx.Exec(qDays, lastID, sd.ProcedureID, sd.DrugID, sd.Order, sd.Times, sd.Frequency)
		if err != nil {
			return 0, err
		}
	}

	return uint(lastID), tx.Commit()
}
