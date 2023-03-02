package migration

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(Up00002, Down00002)
}

func Up00002(tx *sql.Tx) error {
	const (
		q1 = `
		CREATE TABLE IF NOT EXISTS illnesses
			(
				id    INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
				title VARCHAR(100) NOT NULL
			) engine = InnoDB;`
		q2 = `CREATE TABLE IF NOT EXISTS procedures
			(
				id    INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
				title VARCHAR(100) NOT NULL
			) engine = InnoDB;`
		q3 = `CREATE TABLE IF NOT EXISTS schemes
			(
				id         INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
				illness_id INT UNSIGNED NOT NULL,
				length     INT UNSIGNED NOT NULL
			) engine = InnoDB;`
		q4 = `CREATE TABLE IF NOT EXISTS scheme_days
			(
				id           INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
				scheme_id    INT UNSIGNED NOT NULL,
				procedure_id INT UNSIGNED NOT NULL,
				drug_id      INT UNSIGNED DEFAULT NULL,
				day_number   INT UNSIGNED NOT NULL,
				times        INT UNSIGNED NOT NULL,
				each_hours   INT UNSIGNED NOT NULL
			) engine = InnoDB;`
		q5 = `ALTER TABLE schemes
				ADD CONSTRAINT fk_schemes_illnesses 
					FOREIGN KEY (illness_id) 
						REFERENCES illnesses (id) 
						ON DELETE CASCADE 
						ON UPDATE NO ACTION;`
		q6 = `ALTER TABLE scheme_days
				ADD CONSTRAINT fk_scheme_days_procedures 
					FOREIGN KEY (procedure_id) 
						REFERENCES procedures (id) 
						ON DELETE NO ACTION 
						ON UPDATE NO ACTION;`
		q7 = `ALTER TABLE scheme_days
				ADD CONSTRAINT fk_scheme_days_schemes 
					FOREIGN KEY (scheme_id) 
						REFERENCES schemes (id) 
						ON DELETE NO ACTION
						ON UPDATE NO ACTION;`
		q8 = `ALTER TABLE scheme_days
				ADD CONSTRAINT fk_scheme_days_drugs 
					FOREIGN KEY (drug_id) 
						REFERENCES drugs (id) 
						ON DELETE NO ACTION 
						ON UPDATE NO ACTION;`
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
	if _, err := tx.Exec(q4); err != nil {
		return err
	}
	if _, err := tx.Exec(q5); err != nil {
		return err
	}
	if _, err := tx.Exec(q6); err != nil {
		return err
	}
	if _, err := tx.Exec(q7); err != nil {
		return err
	}
	if _, err := tx.Exec(q8); err != nil {
		return err
	}

	return nil
}

func Down00002(tx *sql.Tx) error {
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
