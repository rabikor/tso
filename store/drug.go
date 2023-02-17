package store

import (
	"treatment-scheme-organizer/module"

	"gorm.io/gorm"
)

type DrugStore struct {
	db *gorm.DB
}

func NewDrugStore(DB *gorm.DB) *DrugStore {
	return &DrugStore{db: DB}
}

func (dr *DrugStore) GetAll(limit, page uint) ([]module.Drug, error) {
	var drugs []module.Drug

	result := dr.db.Limit(int(limit)).Offset(int((page - 1) * limit)).Find(&drugs)

	if result.Error != nil {
		panic(result.Error)
	}

	return drugs, nil
}

func (dr *DrugStore) Add(drug *module.Drug) error {
	result := dr.db.Create(&drug)

	if result.Error != nil {
		panic(result.Error)
	}

	return nil
}
