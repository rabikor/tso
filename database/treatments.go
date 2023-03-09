package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
)

type (
	TreatmentsRepository interface {
		ByIllness(illnessID uint, limit, offset int) (ds []Treatment, _ error)
		ByID(id uint) (Treatment, error)
		Add(illnessID uint, begunAt, endedAt string, schemes []TreatmentScheme) (uint, error)
	}
	TreatmentsTable struct {
		*sqlx.DB
	}
)

type Treatment struct {
	ID        uint              `json:"id"`
	IllnessID uint              `db:"illness_id" json:"-"`
	BegunAt   string            `db:"begun_at" json:"begunAt"`
	EndedAt   string            `db:"ended_at" json:"endedAt"`
	Illness   Illness           `db:"illness" json:"-"`
	Schemes   []TreatmentScheme `db:"schemes" json:"-"`
}

func (t Treatment) LastSchemeOrder() uint {
	return uint(len(t.Schemes))
}

func NewTreatmentsRepository(db *sqlx.DB) TreatmentsRepository {
	return TreatmentsTable{db}
}

func (db TreatmentsTable) ByIllness(iID uint, limit, offset int) (ts []Treatment, _ error) {
	return ts, db.Select(
		&ts,
		fmt.Sprintf("SELECT * FROM treatments WHERE illness_id = %d LIMIT %d OFFSET %d", iID, limit, offset),
	)
}

func (db TreatmentsTable) ByID(id uint) (t Treatment, err error) {
	return t, db.Get(&t, "SELECT * FROM treatments WHERE id = ?", id)
}

func (db TreatmentsTable) Add(iID uint, bAt, eAt string, schemes []TreatmentScheme) (uint, error) {
	const q = "INSERT INTO treatments (illness_id, begun_at, ended_at) VALUES (?, ?, ?)"

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

	r, err := db.Exec(q, iID, bAt, eAt)
	if err != nil {
		return 0, err
	}

	lastID, _ := r.LastInsertId()

	for loop, s := range schemes {
		_, err = tx.Exec(qCreateTS, lastID, s.SchemeID, s.BeginFromDay, loop+1)
		if err != nil {
			return 0, err
		}
	}
	return uint(lastID), nil
}
