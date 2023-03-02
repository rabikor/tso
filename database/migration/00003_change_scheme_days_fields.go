package migration

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(Up00003, Down00003)
}

func Up00003(tx *sql.Tx) error {
	const (
		q1 = "ALTER TABLE scheme_days RENAME COLUMN day_number to `order`;"
		q2 = "ALTER TABLE scheme_days RENAME COLUMN each_hours to frequency;"
		q3 = "ALTER TABLE scheme_days ADD CONSTRAINT unq_scheme_days UNIQUE (scheme_id, procedure_id, drug_id, `order`);"
	)

	if _, err := tx.Exec(q1); err != nil {
		return err
	}
	if _, err := tx.Exec(q2); err != nil {
		return err
	}
	if _, err := tx.Exec(q3); err != nil {
		return err
	}

	return nil
}

func Down00003(tx *sql.Tx) error {
	const q = `
		DROP TABLE IF EXISTS illnesses;
		DROP TABLE IF EXISTS procedures;
		DROP TABLE IF EXISTS schemes;
		DROP TABLE IF EXISTS scheme_days;
		
		ALTER TABLE schemes
			DROP CONSTRAINT fk_schemes_illnesses;
		ALTER TABLE scheme_days
			DROP CONSTRAINT fk_scheme_days_procedures;
		ALTER TABLE scheme_days
			DROP CONSTRAINT fk_scheme_days_schemes;
		ALTER TABLE scheme_days
			DROP CONSTRAINT fk_scheme_days_drugs;`

	_, err := tx.Exec(q)
	if err != nil {
		return err
	}
	return nil
}
