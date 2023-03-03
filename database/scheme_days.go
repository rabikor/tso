package database

import (
	"github.com/jmoiron/sqlx"
)

type (
	SchemeDaysRepository interface {
		ByScheme(schemeID, limit, offset int) ([]SchemeDay, error)
		Add(sd SchemeDay) (uint, error)
	}
	SchemeDaysTable struct {
		*sqlx.DB
	}
)

type SchemeDay struct {
	ID          uint `json:"id"`
	SchemeID    uint `db:"scheme_id" json:"-" bson:"-"`
	ProcedureID uint `db:"procedure_id" json:"-" bson:"-"`
	DrugID      uint `db:"drug_id" json:"-" bson:"-"`
	Order       uint `db:"order" json:"dayNumber"`
	Times       uint `db:"times" json:"times"`
	Frequency   uint `db:"frequency" json:"eachHours"`

	Drug      Drug      `db:"drug" json:"drug"`
	Procedure Procedure `db:"procedure" json:"procedure"`
	Scheme    Scheme    `db:"scheme" json:"-"`
}

func NewSchemeDaysRepository(db *sqlx.DB) SchemeDaysRepository {
	return SchemeDaysTable{db}
}

func (db SchemeDaysTable) ByScheme(schemeID, limit, offset int) (sds []SchemeDay, _ error) {
	const q = `SELECT sd.*, 
		       d.id as "drug.id", 
		       d.title as "drug.title",
		       p.id as "procedure.id", 
		       p.title as "procedure.title"
		FROM scheme_days sd
		INNER JOIN drugs d on sd.drug_id = d.id
		INNER JOIN procedures p on sd.procedure_id = p.id
		INNER JOIN schemes s on sd.scheme_id = s.id
		WHERE sd.scheme_id = ? 
		LIMIT ? OFFSET ?`

	return sds, db.Select(&sds, q, schemeID, limit, offset)
}

const qCreateSD = "INSERT INTO scheme_days (scheme_id, procedure_id, drug_id, `order`, times, frequency) VALUES (?, ?, ?, ?, ?, ?)"

// Add creates a new scheme day in the table.
// Do not use integrated relational fields to pass the arguments.
// For example, use SchemeID instead od Scheme.ID to pass the actual scheme ID.
func (db SchemeDaysTable) Add(sd SchemeDay) (uint, error) {
	r, err := db.Exec(
		qCreateSD,
		sd.SchemeID,
		sd.ProcedureID,
		sd.DrugID,
		sd.Order,
		sd.Times,
		sd.Frequency,
	)
	if err != nil {
		return 0, err
	}

	lastID, _ := r.LastInsertId()
	return uint(lastID), nil
}
