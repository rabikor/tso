package migration

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(Up00004, Down00004)
}

func Up00004(tx *sql.Tx) error {
	const (
		q1 = `CREATE  TABLE treatments ( 
					id INT UNSIGNED NOT NULL AUTO_INCREMENT  PRIMARY KEY,
					illness_id INT UNSIGNED NOT NULL,
					begun_at TIMESTAMP NOT NULL,
					ended_at TIMESTAMP DEFAULT (NULL)    
			) engine=InnoDB;`
		q2 = `ALTER TABLE treatments ADD CONSTRAINT fk_treatments_illnesses FOREIGN KEY ( illness_id ) REFERENCES illnesses ( id ) ON DELETE NO ACTION ON UPDATE NO ACTION;`
	)

	if _, err := tx.Exec(q1); err != nil {
		return err
	}

	if _, err := tx.Exec(q2); err != nil {
		return err
	}

	return nil
}

func Down00004(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE IF EXISTS treatments; ALTER TABLE treatments DROP CONSTRAINT fk_treatments_illnesses;`)
	if err != nil {
		return err
	}
	return nil
}
