package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
)

type (
	SchemesRepository interface {
		ByID(id uint) (s Scheme, err error)
		ByIllness(illnessID, limit, offset int) ([]Scheme, error)
		Add(s Scheme) (uint, error)
	}
	SchemesTable struct {
		*sqlx.DB
	}
)

type Scheme struct {
	ID        uint `json:"id"`
	IllnessID uint `db:"illness_id" json:"-" bson:"-"`
	Length    uint `json:"length"`

	Illness Illness     `db:"illness" json:"-"`
	Days    []SchemeDay `db:"days" json:"-"`
}

func NewSchemesRepository(db *sqlx.DB) SchemesRepository {
	return SchemesTable{db}
}

func (db SchemesTable) ByIllness(illnessID, limit, offset int) (ss []Scheme, _ error) {
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

func (db SchemesTable) ByID(id uint) (s Scheme, err error) {
	if err = db.Get(&s, "SELECT * FROM schemes WHERE id = ?", id); err != nil {
		return s, err
	}

	if err = db.Select(&s.Days, "SELECT * FROM scheme_days WHERE scheme_id = ?", s.ID); err != nil {
		return s, err
	}

	return s, nil
}

func (db SchemesTable) Add(s Scheme) (uint, error) {
	const q = "INSERT INTO schemes (illness_id, length) VALUES (?, ?)"

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

	r, err := tx.Exec(qCreateSD, s.IllnessID, s.Length)
	if err != nil {
		return 0, err
	}

	lastID, _ := r.LastInsertId()
	for _, sd := range s.Days {
		_, err = tx.Exec(q, lastID, sd.ProcedureID, sd.DrugID, sd.Order, sd.Times, sd.Frequency)
		if err != nil {
			return 0, err
		}
	}

	return uint(lastID), tx.Commit()
}
