package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type (
	drugsRepository interface {
		GetAll(limit, offset int) ([]Drug, error)
		GetById(id uint) (Drug, error)
		Add(drg Drug) error
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
	return ds, db.Select(&ds, fmt.Sprintf("SELECT * FROM drugs LIMIT %d OFFSET %d", limit, offset))
}

func (db *drugsTable) GetById(id uint) (drug Drug, err error) {
	return drug, db.Get(&drug, "SELECT * FROM drugs WHERE id = ?", id)
}

func (db *drugsTable) Add(drg Drug) (int, error) {
	result, err := db.Exec("INSERT INTO drugs (title) VALUES (?)", drg.Title)
	if err != nil {
		return err
	}
	
	id, _ := result.LastInsertId()
	return uint(id), nil
}
