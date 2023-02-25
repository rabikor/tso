package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type (
	drugsRepository interface {
		GetAll(limit, offset int) ([]Drug, error)
		GetById(id uint) (*Drug, error)
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

func (db *drugsTable) GetById(id uint) (drug *Drug, err error) {
	drug = new(Drug)
	if err = db.Get(drug, "SELECT * FROM drugs WHERE id = ?", id); err != nil {
		return
	}

	return
}

func (db *drugsTable) Add(drg *Drug) error {
	if result, err := db.Exec("INSERT INTO drugs (title) VALUES (?)", drg.Title); err != nil {
		return err
	} else {
		lastId, _ := result.LastInsertId()
		drg.ID = uint(lastId)
	}

	return nil
}
