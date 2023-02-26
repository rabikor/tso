package database

import (
	"github.com/jmoiron/sqlx"
)

type (
	schemeDaysRepository interface {
		GetByScheme(schemeID, limit, offset int) ([]SchemeDay, error)
		Add(sd *SchemeDay) error
	}
	schemeDaysTable struct {
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

func (db *schemeDaysTable) GetByScheme(schemeID, limit, offset int) (sds []SchemeDay, _ error) {
	q := `
		SELECT sd.*, 
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

	err := db.Select(
		&sds,
		q,
		schemeID,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}

	return sds, nil
}

func (db *schemeDaysTable) Add(schemeDay *SchemeDay) error {
	var q = "INSERT INTO scheme_days (scheme_id, procedure_id, drug_id, `order`, times, frequency) VALUES (?, ?, ?, ?, ?, ?)"

	if result, err := db.Exec(
		q,
		schemeDay.Scheme.ID,
		schemeDay.Procedure.ID,
		schemeDay.Drug.ID,
		schemeDay.Order,
		schemeDay.Times,
		schemeDay.Frequency,
	); err != nil {
		return err
	} else {
		lastId, _ := result.LastInsertId()
		schemeDay.ID = uint(lastId)
	}

	return nil
}
