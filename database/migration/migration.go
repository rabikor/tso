package migration

import (
	"embed"
	"errors"
	"io/fs"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
	"github.com/pressly/goose/v3"
)

func Migrate(db *sqlx.DB) error {
	var embedMigrations embed.FS

	const (
		dialect = "mysql"
		dir     = "."
	)

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(dialect); err != nil {
		return err
	}

	if _, err := fs.Stat(embedMigrations, "."); errors.Is(err, fs.ErrNotExist) {
		log.Error(err)
	}

	if err := goose.Up(db.DB, dir); err != nil {
		return err
	}

	return nil
}
