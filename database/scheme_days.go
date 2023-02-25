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
	DayNumber   uint `db:"day_number" json:"dayNumber"`
	Times       uint `json:"times"`
	EachHours   uint `db:"each_hours" json:"eachHours"`

	Drug      Drug      `db:"drug" json:"drug"`
	Procedure Procedure `db:"procedure" json:"procedure"`
	Scheme    Scheme    `db:"scheme" json:"-"`
}

func (db *schemeDaysTable) GetByScheme(schemeID, limit, offset int) (sds []SchemeDay, _ error) {
	sql := `
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
		sql,
		schemeID,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}

	return sds, nil
}

var sqlCreateSchemeDay string = `
	INSERT INTO scheme_days 
	    (scheme_id, procedure_id, drug_id, day_number, times, each_hours) 
	VALUES (?, ?, ?, ?, ?, ?)`

func (db *schemeDaysTable) Add(schemeDay *SchemeDay) error {
	if result, err := db.Exec(
		sqlCreateSchemeDay,
		schemeDay.Scheme.ID,
		schemeDay.Procedure.ID,
		schemeDay.Drug.ID,
		schemeDay.DayNumber,
		schemeDay.Times,
		schemeDay.EachHours,
	); err != nil {
		return err
	} else {
		lastId, _ := result.LastInsertId()
		schemeDay.ID = uint(lastId)
	}

	return nil
}
