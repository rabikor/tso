package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type (
	drugsRepository interface {
		GetAll(limit, offset int) ([]Drug, error)
		Add(drg *Drug) error
	}
	drugsTable struct {
		*sqlx.DB
	}
)

type Drug struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

func (db *drugsTable) GetAll(limit, offset int) (ds []Drug, _ error) {
	err := db.Select(&ds, fmt.Sprintf("SELECT * FROM drugs LIMIT %d OFFSET %d", limit, offset))
	if err != nil {
		return nil, err
	}

	return ds, nil
}

func (db *drugsTable) Add(drg *Drug) error {
	if _, err := db.Exec("INSERT INTO drugs (title) VALUES (?)", drg.Title); err != nil {
		return err
	}

	return nil
}
