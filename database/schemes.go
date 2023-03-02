package database

import (
	"github.com/jmoiron/sqlx"
)

type (
	schemesRepository interface {
		GetById(id uint) (scheme *Scheme, err error)
		GetByIllness(illnessID, limit, offset int) ([]Scheme, error)
		Add(illnessID, length uint, days []SchemeDay) error
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

	Illness Illness     `db:"illness" json:"-"`
	Days    []SchemeDay `db:"days" json:"-"`
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

func (db *schemesTable) GetById(id uint) (s *Scheme, err error) {
	s = new(Scheme)
	if err = db.Get(s, "SELECT * FROM schemes WHERE id = ?", id); err != nil {
		return
	}

	if err = db.Select(&s.Days, "SELECT * FROM scheme_days WHERE scheme_id = ?", s.ID); err != nil {
		return
	}

	return
}

func (db *schemesTable) Add(illnessID, length uint, days []SchemeDay) error {
	const (
		qScheme = "INSERT INTO schemes (illness_id, length) VALUES (?, ?)"
		qDays   = "INSERT INTO scheme_days (scheme_id, procedure_id, drug_id, `order`, times, frequency) VALUES (?, ?, ?, ?, ?, ?)"
	)
	
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				// TODO: log the rollback error
			}
		}
	}()

	var result *sqlx.Result
	result, err = tx.Exec(qScheme, illnessID, length)
	if err != nil {
		return err
	}
	
	lastID, _ := result.LastInsertId()
	for _, sd := range days {
		_, err = tx.Exec(qDays, lastID, sd.ProcedureID, sd.DrugID, sd.Order, sd.Times, sd.Frequency)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
