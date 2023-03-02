package migration

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(Up00001, Down00001)
}

func Up00001(tx *sql.Tx) error {
	const q = `CREATE TABLE IF NOT EXISTS drugs
		(
			id    INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(100) NOT NULL
		) engine = InnoDB;`

	_, err := tx.Exec(q)
	if err != nil {
		return err
	}
	return nil
}

func Down00001(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE IF EXISTS drugs;")
	if err != nil {
		return err
	}
	return nil
}
