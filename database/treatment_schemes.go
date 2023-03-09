package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type (
	TreatmentSchemesRepository interface {
		CountByTreatment(treatmentID uint) (count uint, _ error)
		ByTreatment(treatmentID uint, limit, offset int) (tss []TreatmentScheme, _ error)
		Add(treatmentID, schemeID, beginFromDay, order uint) (uint, error)
	}
	TreatmentSchemesTable struct {
		*sqlx.DB
	}
)

type TreatmentScheme struct {
	ID           uint      `json:"id"`
	TreatmentID  uint      `db:"treatment_id" json:"-"`
	SchemeID     uint      `db:"scheme_id" json:"-"`
	BeginFromDay uint      `db:"begin_from_day" json:"beginFromDay"`
	Order        uint      `db:"order" json:"order"`
	Treatment    Treatment `db:"treatment" json:"-"`
	Scheme       Scheme    `db:"scheme" json:"-"`
}

func NewTreatmentSchemesRepository(db *sqlx.DB) TreatmentSchemesRepository {
	return TreatmentSchemesTable{db}
}

func (db TreatmentSchemesTable) CountByTreatment(tID uint) (count uint, _ error) {
	row := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM treatment_schemes WHERE treatment_id = %d", tID))
	return count, row.Scan(&count)
}

func (db TreatmentSchemesTable) ByTreatment(tID uint, limit, offset int) (tss []TreatmentScheme, _ error) {
	return tss, db.Select(
		&tss,
		fmt.Sprintf("SELECT * FROM treatment_schemes WHERE treatment_id = %d LIMIT %d OFFSET %d", tID, limit, offset),
	)
}

const qCreateTS = "INSERT INTO treatment_schemes (treatment_id, scheme_id, begin_from_day, `order`) VALUES (?, ?, ?, ?)"

func (db TreatmentSchemesTable) Add(tID, sID, bfd, o uint) (uint, error) {
	result, err := db.Exec(qCreateTS, tID, sID, bfd, o)
	if err != nil {
		return 0, err
	}

	id, _ := result.LastInsertId()
	return uint(id), nil
}
