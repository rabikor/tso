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

func (db *schemesTable) Add(s *Scheme) error {
	var q string
	q = "INSERT INTO schemes (illness_id, length) VALUES (?, ?)"

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if s.IllnessID == 0 && s.Illness.ID > 0 {
		s.IllnessID = s.Illness.ID
	}

	if res, err := tx.Exec(q, s.IllnessID, s.Length); err != nil {
		tx.Rollback()
		return err
	} else {
		lastId, _ := res.LastInsertId()
		s.ID = uint(lastId)
	}

	q = "INSERT INTO scheme_days (scheme_id, procedure_id, drug_id, `order`, times, frequency) VALUES (?, ?, ?, ?, ?, ?)"

	for _, sd := range s.Days {
		if res, err := tx.Exec(q, s.ID, sd.ProcedureID, sd.DrugID, sd.Order, sd.Times, sd.Frequency); err != nil {
			tx.Rollback()
			return err
		} else {
			lastId, _ := res.LastInsertId()
			sd.ID = uint(lastId)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
